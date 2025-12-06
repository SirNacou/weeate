package jobs

import (
"context"
"fmt"
"log/slog"
"os"
"testing"
"time"

"github.com/SirNacou/weeate/backend/internal/common/events"
"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
"github.com/SirNacou/weeate/backend/internal/common/infrastructure/data"
"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
"github.com/SirNacou/weeate/backend/internal/features/polls/services"
"github.com/ThreeDotsLabs/watermill/components/cqrs"
"github.com/gofrs/uuid/v5"
"gorm.io/gorm"
)

// setupEndToEndTestDB creates a test database connection for end-to-end tests
func setupEndToEndTestDB(t *testing.T) *gorm.DB {
t.Helper()

dbHost := os.Getenv("TEST_DB_HOST")
if dbHost == "" {
dbHost = "localhost"
}

dbPort := os.Getenv("TEST_DB_PORT")
if dbPort == "" {
dbPort = "5432"
}

dbUser := os.Getenv("TEST_DB_USER")
if dbUser == "" {
dbUser = "weeate_user"
}

dbPassword := os.Getenv("TEST_DB_PASSWORD")
if dbPassword == "" {
dbPassword = "weeate_password"
}

dbName := os.Getenv("TEST_DB_NAME")
if dbName == "" {
dbName = "weeate_test_db"
}

dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
dbHost, dbUser, dbPassword, dbName, dbPort)

ctx := context.Background()
db, err := data.ConnectToPostgres(ctx, dsn)
if err != nil {
t.Skipf("Skipping integration test: could not connect to test database: %v", err)
return nil
}

if err := data.MigratePostgresDB(db); err != nil {
t.Fatalf("Failed to migrate test database: %v", err)
}

cleanupTestDB(t, db)

return db
}

// TestClosePollScheduler_EndToEnd_EventDelivery tests the complete flow from closing a poll to receiving the event
func TestClosePollScheduler_EndToEnd_EventDelivery(t *testing.T) {
db := setupEndToEndTestDB(t)
if db == nil {
return
}
defer cleanupTestDB(t, db)

sqlDB, err := db.DB()
if err != nil {
t.Fatalf("Failed to get sql.DB: %v", err)
}

// Create logger
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
Level: slog.LevelDebug,
}))

// Create Bus
testBus, err := bus.NewBus(sqlDB, logger)
if err != nil {
t.Fatalf("Failed to create bus: %v", err)
}

// Track received events
eventReceived := make(chan *events.PollClosedEvent, 1)

// Register event handler
testBus.EventProcessor.AddHandlers(
cqrs.NewEventHandler("test.onPollClosed", func(ctx context.Context, event *events.PollClosedEvent) error {
t.Logf("Received PollClosedEvent: PollID=%v, BuyerID=%v", event.PollID, event.BuyerID)
eventReceived <- event
return nil
}),
)

// Start the bus in a goroutine
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

go func() {
if err := testBus.Start(ctx); err != nil {
t.Logf("Bus stopped: %v", err)
}
}()

// Wait for bus to start
time.Sleep(500 * time.Millisecond)

// Create test poll
orderDate := time.Now().UTC().Truncate(24 * time.Hour)
scheduledClosesAt := time.Now().UTC().Add(-1 * time.Hour) // Already past due
pollOptions := []domain.PollOption{
*domain.NewPollOption(uuid.Nil, uuid.Must(uuid.NewV7()), 1099),
}
poll, err := domain.NewPoll(orderDate, scheduledClosesAt, domain.OrderConsensus, pollOptions, "test-buyer-123")
if err != nil {
t.Fatalf("Failed to create poll: %v", err)
}
if err := db.Create(poll).Error; err != nil {
t.Fatalf("Failed to save poll: %v", err)
}

t.Logf("Created poll with ID: %v", poll.ID)

// Create close poll handler
closePollHandler := services.NewClosePollCommandHandler(db, testBus, nil)

// Close the poll
t.Log("Closing poll...")
err = closePollHandler.Handle(ctx, services.ClosePollCommand{PollID: poll.ID})
if err != nil {
t.Fatalf("Failed to close poll: %v", err)
}
t.Log("Poll closed successfully")

// Wait for event
select {
case event := <-eventReceived:
t.Logf("SUCCESS: Event received! PollID=%v", event.PollID)
if event.PollID != poll.ID {
t.Errorf("Expected PollID %v, got %v", poll.ID, event.PollID)
}
case <-time.After(10 * time.Second):
t.Error("FAILURE: Timeout waiting for PollClosedEvent - this indicates the event was not delivered!")
}

// Verify poll is closed
var updatedPoll domain.Poll
if err := db.Where("id = ?", poll.ID).First(&updatedPoll).Error; err != nil {
t.Fatalf("Failed to fetch poll: %v", err)
}
if updatedPoll.ClosedAt == nil {
t.Error("Poll should be closed but ClosedAt is nil")
}
}

// TestClosePollScheduler_EndToEnd_SchedulerToEvent tests the full flow from scheduler to event
func TestClosePollScheduler_EndToEnd_SchedulerToEvent(t *testing.T) {
db := setupEndToEndTestDB(t)
if db == nil {
return
}
defer cleanupTestDB(t, db)

sqlDB, err := db.DB()
if err != nil {
t.Fatalf("Failed to get sql.DB: %v", err)
}

// Create logger
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
Level: slog.LevelDebug,
}))

// Create Bus
testBus, err := bus.NewBus(sqlDB, logger)
if err != nil {
t.Fatalf("Failed to create bus: %v", err)
}

// Track received events
eventReceived := make(chan *events.PollClosedEvent, 1)

// Register event handler
testBus.EventProcessor.AddHandlers(
cqrs.NewEventHandler("test.onPollClosed", func(ctx context.Context, event *events.PollClosedEvent) error {
t.Logf("Received PollClosedEvent from scheduler: PollID=%v", event.PollID)
eventReceived <- event
return nil
}),
)

// Start the bus
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

go func() {
if err := testBus.Start(ctx); err != nil {
t.Logf("Bus stopped: %v", err)
}
}()

// Wait for bus to start
time.Sleep(500 * time.Millisecond)

// Create close poll handler
closePollHandler := services.NewClosePollCommandHandler(db, testBus, nil)

// Create scheduler
scheduler := NewClosePollScheduler(closePollHandler, db)

// Create test poll that is past due
orderDate := time.Now().UTC().Truncate(24 * time.Hour)
scheduledClosesAt := time.Now().UTC().Add(-1 * time.Minute) // Already past
pollOptions := []domain.PollOption{
*domain.NewPollOption(uuid.Nil, uuid.Must(uuid.NewV7()), 1099),
}
poll, err := domain.NewPoll(orderDate, scheduledClosesAt, domain.OrderConsensus, pollOptions, "test-buyer-scheduler")
if err != nil {
t.Fatalf("Failed to create poll: %v", err)
}
if err := db.Create(poll).Error; err != nil {
t.Fatalf("Failed to save poll: %v", err)
}

t.Logf("Created poll with ID: %v, scheduled to close at: %v", poll.ID, scheduledClosesAt)

// Start scheduler in goroutine
go scheduler.Start(ctx)

// Wait for scheduler to pick it up and close it
select {
case event := <-eventReceived:
t.Logf("SUCCESS: Scheduler flow completed! PollID=%v", event.PollID)
if event.PollID != poll.ID {
t.Errorf("Expected PollID %v, got %v", poll.ID, event.PollID)
}
case <-time.After(15 * time.Second):
t.Error("FAILURE: Timeout waiting for scheduler to close poll and deliver event!")
}
}
