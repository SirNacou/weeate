package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/data"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/services"
	"github.com/gofrs/uuid/v5"
	"github.com/jonboulle/clockwork"
	"gorm.io/gorm"
)

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Use environment variables for test database configuration
	// Default to localhost if not set
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

	// Auto-migrate the schema
	if err := data.MigratePostgresDB(db); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Clean up tables before test
	cleanupTestDB(t, db)

	return db
}

// cleanupTestDB removes all test data
func cleanupTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Clean up in reverse dependency order
	db.Exec("DELETE FROM votes")
	db.Exec("DELETE FROM poll_options")
	db.Exec("DELETE FROM polls")
	db.Exec("DELETE FROM orders")
}

// setupTestDependencies creates the necessary dependencies for testing
func setupTestDependencies(t *testing.T, db *gorm.DB) (*services.ClosePollCommandHandler, error) {
	t.Helper()

	// Get underlying sql.DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Set to Error to reduce test noise
	}))

	// Create Bus
	testBus, err := bus.NewBus(sqlDB, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create bus: %w", err)
	}

	// Create Centrifugo client with test config
	testConfig := configs.Config{
		CENTRIFUGO_GRPC_HOST: "localhost",
		CENTRIFUGO_GRPC_PORT: 10000,
	}
	testCentrifugo, err := centrifugo.NewCentrifugoClient(testConfig.CENTRIFUGO_GRPC_HOST, testConfig.CENTRIFUGO_GRPC_PORT)
	if err != nil {
		// For integration tests, we can skip if centrifugo is not available
		t.Logf("Warning: Could not create Centrifugo client: %v", err)
		// Create a nil centrifugo for tests that don't need it
		testCentrifugo = nil
	}

	// Create ClosePollCommandHandler
	closePollHandler := services.NewClosePollCommandHandler(db, testBus, testCentrifugo)

	return closePollHandler, nil
}

func createTestPoll(t *testing.T, db *gorm.DB, scheduledClosesAt time.Time, buyerID string) *domain.Poll {

	orderDate := time.Now().UTC().Truncate(24 * time.Hour)

	// Create poll options using the constructor
	pollOptions := []domain.PollOption{
		*domain.NewPollOption(uuid.Nil, uuid.Must(uuid.NewV7()), 1099), // Price in cents
		*domain.NewPollOption(uuid.Nil, uuid.Must(uuid.NewV7()), 1599), // Price in cents
	}

	poll, err := domain.NewPoll(orderDate, scheduledClosesAt, domain.OrderConsensus, pollOptions, buyerID)
	if err != nil {
		t.Fatalf("Failed to create test poll: %v", err)
	}

	result := db.Create(poll)
	if result.Error != nil {
		t.Fatalf("Failed to save test poll: %v", result.Error)
	}

	return poll
}

func TestClosePollScheduler_Integration_CloseDuePolls(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Test was skipped
	}
	defer cleanupTestDB(t, db)

	// Setup dependencies
	closePollHandler, err := setupTestDependencies(t, db)
	if err != nil {
		t.Skipf("Skipping test due to dependency setup failure: %v", err)
		return
	}

	// Create a poll that should be closed (scheduled in the past)
	pastTime := time.Now().UTC().Add(-1 * time.Hour)
	poll1 := createTestPoll(t, db, pastTime, "buyer1")

	// Create a poll that should NOT be closed (scheduled in the future)
	futureTime := time.Now().UTC().Add(2 * time.Hour)
	poll2 := createTestPoll(t, db, futureTime, "buyer2")

	// Create scheduler
	scheduler := NewClosePollScheduler(closePollHandler, db, clockwork.NewRealClock())

	// Call closeDuePolls directly
	ctx := context.Background()
	scheduler.closeDuePolls(ctx)

	// Verify poll1 is closed
	var updatedPoll1 domain.Poll
	err = db.Where("id = ?", poll1.ID).First(&updatedPoll1).Error
	if err != nil {
		t.Fatalf("Failed to fetch poll1: %v", err)
	}

	if updatedPoll1.ClosedAt == nil {
		t.Error("Expected poll1 to be closed, but ClosedAt is nil")
	}

	// Verify poll2 is NOT closed
	var updatedPoll2 domain.Poll
	err = db.Where("id = ?", poll2.ID).First(&updatedPoll2).Error
	if err != nil {
		t.Fatalf("Failed to fetch poll2: %v", err)
	}

	if updatedPoll2.ClosedAt != nil {
		t.Error("Expected poll2 to NOT be closed, but ClosedAt is not nil")
	}
}

func TestClosePollScheduler_Integration_StartAndTrigger(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Test was skipped
	}
	defer cleanupTestDB(t, db)

	// Setup dependencies
	closePollHandler, err := setupTestDependencies(t, db)
	if err != nil {
		t.Skipf("Skipping test due to dependency setup failure: %v", err)
		return
	}

	// Create scheduler with fake clock
	fakeClock := clockwork.NewFakeClock()
	scheduler := NewClosePollScheduler(closePollHandler, db, fakeClock)

	// Create a poll that will close in 5 seconds
	scheduledTime := fakeClock.Now().Add(5 * time.Second)
	poll := createTestPoll(t, db, scheduledTime, "buyer1")

	// Start scheduler in a goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go scheduler.Start(ctx)

	// Wait a bit for scheduler to start
	time.Sleep(100 * time.Millisecond)

	// Advance the fake clock by 6 seconds to trigger the poll closure
	fakeClock.Advance(6 * time.Second)

	// Wait for the poll to be processed
	time.Sleep(500 * time.Millisecond)

	// Verify the poll is closed
	var updatedPoll domain.Poll
	err = db.Where("id = ?", poll.ID).First(&updatedPoll).Error
	if err != nil {
		t.Fatalf("Failed to fetch poll: %v", err)
	}

	if updatedPoll.ClosedAt == nil {
		t.Error("Expected poll to be closed after time advancement, but ClosedAt is nil")
	}
}

func TestClosePollScheduler_Integration_TriggerUpdate(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Test was skipped
	}
	defer cleanupTestDB(t, db)

	// Setup dependencies
	closePollHandler, err := setupTestDependencies(t, db)
	if err != nil {
		t.Skipf("Skipping test due to dependency setup failure: %v", err)
		return
	}

	// Create scheduler
	scheduler := NewClosePollScheduler(closePollHandler, db, clockwork.NewRealClock())

	// Start scheduler in a goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go scheduler.Start(ctx)

	// Wait for scheduler to start
	time.Sleep(100 * time.Millisecond)

	// Trigger an update - this should wake up the scheduler
	scheduler.TriggerUpdate()

	// Wait a bit to ensure the trigger was processed
	time.Sleep(200 * time.Millisecond)

	// Test passes if no panic occurs and the scheduler handles the trigger gracefully
	// Additional verification: try triggering multiple times (non-blocking)
	for i := 0; i < 5; i++ {
		scheduler.TriggerUpdate()
	}

	time.Sleep(100 * time.Millisecond)
}

func TestClosePollScheduler_Integration_MultiplePolls(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Test was skipped
	}
	defer cleanupTestDB(t, db)

	// Setup dependencies
	closePollHandler, err := setupTestDependencies(t, db)
	if err != nil {
		t.Skipf("Skipping test due to dependency setup failure: %v", err)
		return
	}

	now := time.Now().UTC()

	// Create multiple polls with different scheduled times
	poll1 := createTestPoll(t, db, now.Add(-2*time.Hour), "buyer1")    // Past - should close
	poll2 := createTestPoll(t, db, now.Add(-1*time.Hour), "buyer2")    // Past - should close
	poll3 := createTestPoll(t, db, now.Add(1*time.Hour), "buyer3")     // Future - should NOT close
	poll4 := createTestPoll(t, db, now.Add(-30*time.Minute), "buyer4") // Past - should close

	// Create scheduler
	scheduler := NewClosePollScheduler(closePollHandler, db, clockwork.NewRealClock())

	// Close due polls
	ctx := context.Background()
	scheduler.closeDuePolls(ctx)

	// Verify polls 1, 2, and 4 are closed
	pollsToCheck := []struct {
		poll           *domain.Poll
		shouldBeClosed bool
		name           string
	}{
		{poll1, true, "poll1"},
		{poll2, true, "poll2"},
		{poll3, false, "poll3"},
		{poll4, true, "poll4"},
	}

	for _, tc := range pollsToCheck {
		var updatedPoll domain.Poll
		err := db.Where("id = ?", tc.poll.ID).First(&updatedPoll).Error
		if err != nil {
			t.Fatalf("Failed to fetch %s: %v", tc.name, err)
		}

		if tc.shouldBeClosed && updatedPoll.ClosedAt == nil {
			t.Errorf("Expected %s to be closed, but ClosedAt is nil", tc.name)
		}

		if !tc.shouldBeClosed && updatedPoll.ClosedAt != nil {
			t.Errorf("Expected %s to NOT be closed, but ClosedAt is %v", tc.name, updatedPoll.ClosedAt)
		}
	}
}

func TestClosePollScheduler_Integration_NoPolls(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Test was skipped
	}
	defer cleanupTestDB(t, db)

	// Setup dependencies
	closePollHandler, err := setupTestDependencies(t, db)
	if err != nil {
		t.Skipf("Skipping test due to dependency setup failure: %v", err)
		return
	}

	// Create scheduler
	scheduler := NewClosePollScheduler(closePollHandler, db, clockwork.NewRealClock())

	// Start scheduler with no polls in the database
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// This should not panic and should sleep for 1 hour (default when no polls exist)
	go scheduler.Start(ctx)

	// Wait for context to timeout
	<-ctx.Done()

	// Test passes if no panic occurs
}
