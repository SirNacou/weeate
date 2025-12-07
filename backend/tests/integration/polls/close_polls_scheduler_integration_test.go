//go:build integration

package polls

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	food_domain "github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/infrastructure/jobs"
	"github.com/SirNacou/weeate/backend/internal/features/polls/services"
	"github.com/SirNacou/weeate/backend/tests/testhelpers"
	"github.com/gofrs/uuid/v5"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	db               *gorm.DB
	centrifugoClient *centrifugo.CentrifugoClient
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Suppress library logs during tests
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors
	})))

	testdb, cleanup, err := testhelpers.SetupGORM(ctx)

	if err != nil {
		panic(err)
	}

	client, centrifugoCleanup, err := testhelpers.SetupCentrifugoClient(ctx)
	if err != nil {
		panic(err)
	}

	db = testdb
	centrifugoClient = client

	exitCode := m.Run()

	cleanup()
	centrifugoCleanup()
	os.Exit(exitCode)
}

// cleanupTestData removes all test data from the database to prevent test pollution
func cleanupTestData(t *testing.T) {
	t.Helper()
	// Delete in order to respect foreign key constraints
	if err := db.Exec("DELETE FROM votes").Error; err != nil {
		t.Fatalf("failed to clean up votes: %v", err)
	}
	if err := db.Exec("DELETE FROM poll_options").Error; err != nil {
		t.Fatalf("failed to clean up poll_options: %v", err)
	}
	if err := db.Exec("DELETE FROM polls").Error; err != nil {
		t.Fatalf("failed to clean up polls: %v", err)
	}
	if err := db.Exec("DELETE FROM foods").Error; err != nil {
		t.Fatalf("failed to clean up foods: %v", err)
	}
}

func TestClosePollsSchedulerIntegration(t *testing.T) {
	// Test implementation goes here
	t.Run("scheduler closes polls when polls reach their end time", func(t *testing.T) {
		cleanupTestData(t)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Arrange - Create fake clock
		fakeClock := clockwork.NewFakeClock()

		expiredPoll, err := domain.NewPoll(fakeClock.Now().Add(24*time.Hour),
			fakeClock.Now().Add(-1*time.Hour),
			domain.OrderMultiple,
			[]domain.PollOption{
				*domain.NewPollOption(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()), 20000)},
			"buyer-1")
		if err != nil {
			t.Fatalf("failed to create expired poll: %v", err)
		}
		err = gorm.G[domain.Poll](db).Create(ctx, expiredPoll)
		if err != nil {
			t.Fatalf("failed to insert expired poll: %v", err)
		}
		sql, err := db.DB()
		if err != nil {
			t.Fatalf("failed to get raw DB: %v", err)
		}

		// Use a discarding logger for cleaner test output
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		b, err := bus.NewBus(sql, logger)
		if err != nil {
			t.Fatalf("failed to create bus: %v", err)
		}

		// Start bus with graceful shutdown
		busCtx, busCancel := context.WithCancel(ctx)
		defer busCancel()

		go func() {
			if err := b.Start(busCtx); err != nil && err != context.Canceled {
				t.Logf("bus stopped: %v", err)
			}
		}()

		// Wait for bus to start
		time.Sleep(100 * time.Millisecond)

		closeCmd := services.NewClosePollCommandHandler(db, b, centrifugoClient)

		job := jobs.NewClosePollScheduler(closeCmd, db, fakeClock)

		// Act - Start scheduler in background
		schedulerCtx, schedulerCancel := context.WithCancel(ctx)
		defer schedulerCancel()

		go job.Start(schedulerCtx)

		// Advance fake clock to trigger poll processing
		fakeClock.Advance(2 * time.Hour)
		time.Sleep(500 * time.Millisecond)

		// Assert
		closedPoll, err := gorm.G[domain.Poll](db).Where("id = ?", expiredPoll.ID).First(ctx)
		if err != nil {
			t.Fatalf("failed to fetch poll after scheduler run: %v", err)
		}

		assert.NotNil(t, closedPoll.ClosedAt, "Poll should be closed by the scheduler")
	})

	t.Run("scheduler reacts to new poll creation", func(t *testing.T) {
		cleanupTestData(t)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Arrange - Create fake clock
		fakeClock := clockwork.NewFakeClock()

		sql, err := db.DB()
		if err != nil {
			t.Fatalf("failed to get raw DB: %v", err)
		}

		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		b, err := bus.NewBus(sql, logger)
		if err != nil {
			t.Fatalf("failed to create bus: %v", err)
		}

		// Start bus with graceful shutdown
		busCtx, busCancel := context.WithCancel(ctx)
		defer busCancel()

		go func() {
			if err := b.Start(busCtx); err != nil && err != context.Canceled {
				t.Logf("bus stopped: %v", err)
			}
		}()

		time.Sleep(100 * time.Millisecond)

		closeCmd := services.NewClosePollCommandHandler(db, b, centrifugoClient)
		scheduler := jobs.NewClosePollScheduler(closeCmd, db, fakeClock)

		// Start scheduler in background
		schedulerCtx, schedulerCancel := context.WithCancel(ctx)
		defer schedulerCancel()

		go scheduler.Start(schedulerCtx)

		// Wait for scheduler to initialize
		time.Sleep(100 * time.Millisecond)

		// Act - Create a poll that expires in 1 second (fake clock time)
		newPoll, err := domain.NewPoll(
			fakeClock.Now().Add(24*time.Hour),
			fakeClock.Now().Add(1*time.Second),
			domain.OrderMultiple,
			[]domain.PollOption{
				*domain.NewPollOption(uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()), 15000)},
			"buyer-2")
		if err != nil {
			t.Fatalf("failed to create new poll: %v", err)
		}

		err = gorm.G[domain.Poll](db).Create(ctx, newPoll)
		if err != nil {
			t.Fatalf("failed to insert new poll: %v", err)
		}

		// Trigger scheduler to recalculate
		scheduler.TriggerUpdate()

		// Advance fake clock to trigger poll closure
		fakeClock.Advance(2 * time.Second)
		time.Sleep(500 * time.Millisecond)

		// Assert
		closedPoll, err := gorm.G[domain.Poll](db).Where("id = ?", newPoll.ID).First(ctx)
		if err != nil {
			t.Fatalf("failed to fetch poll after scheduler run: %v", err)
		}

		assert.NotNil(t, closedPoll.ClosedAt, "Poll should be closed by the scheduler after trigger")

		// Graceful shutdown
		t.Log("Shutting down scheduler and bus...")
		schedulerCancel()
		busCancel()
		time.Sleep(200 * time.Millisecond)
		t.Log("Shutdown complete")
	})

	t.Run("scheduler triggered by poll creation closes poll automatically", func(t *testing.T) {
		cleanupTestData(t)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Arrange - Create fake clock
		fakeClock := clockwork.NewFakeClock()

		sql, err := db.DB()
		if err != nil {
			t.Fatalf("failed to get raw DB: %v", err)
		}

		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		b, err := bus.NewBus(sql, logger)
		if err != nil {
			t.Fatalf("failed to create bus: %v", err)
		}

		busCtx, busCancel := context.WithCancel(ctx)
		defer busCancel()

		go func() {
			if err := b.Start(busCtx); err != nil && err != context.Canceled {
				t.Logf("bus stopped: %v", err)
			}
		}()

		time.Sleep(100 * time.Millisecond)

		// Setup scheduler and command handlers
		closeCmd := services.NewClosePollCommandHandler(db, b, centrifugoClient)
		scheduler := jobs.NewClosePollScheduler(closeCmd, db, fakeClock)
		createCmd := services.NewCreatePollCommandHandler(db, centrifugoClient, scheduler)

		// Start scheduler
		schedulerCtx, schedulerCancel := context.WithCancel(ctx)
		defer schedulerCancel()

		go scheduler.Start(schedulerCtx)
		time.Sleep(100 * time.Millisecond)

		userID := uuid.Must(uuid.NewV7())
		food, err := food_domain.NewFood(
			"Test Food",
			"",
			nil,
			"Delicious test food",
			10.0,
			userID,
		)
		if err != nil {
			t.Fatalf("can not create food: %v", err)
		}

		err = gorm.G[food_domain.Food](db).Create(ctx, food)
		if err != nil {
			t.Fatalf("can not save food: %v", err)
		}

		// Act - Use CreatePollCommandHandler (simulates real application flow)
		// This will automatically call scheduler.TriggerUpdate()
		strategy, err := domain.OrderMultiple.Value()
		req := services.CreatePollCommand{
			OrderDate:        fakeClock.Now().Add(24 * time.Hour),
			ScheduledCloseAt: fakeClock.Now().Add(1 * time.Hour),
			Strategy:         strategy.(string),
			FoodIDs:          []string{food.ID.String()},
		}

		userCtx := auth.WithUser(ctx, &auth.User{
			ID: userID.String(),
		})
		res, err := createCmd.Handle(userCtx, req)
		if err != nil {
			t.Fatalf("failed to create poll: %v", err)
		}

		t.Logf("Poll created with ID %s, scheduled to close at %s", res.PollID, req.ScheduledCloseAt)

		// Advance fake clock to trigger poll closure
		fakeClock.Advance(1 * time.Hour)
		time.Sleep(500 * time.Millisecond)

		// Assert
		closedPoll, err := gorm.G[domain.Poll](db).Where("id = ?", res.PollID).First(ctx)
		if err != nil {
			t.Fatalf("failed to fetch poll: %v", err)
		}

		assert.NotNil(t, closedPoll.ClosedAt, "Poll should be automatically closed by scheduler after creation")
		t.Logf("Poll successfully closed at %s", closedPoll.ClosedAt)
	})
}

func TestSchedulerTriggerWhenPollCreated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Arrange - Create fake clock
	fakeClock := clockwork.NewFakeClock()

	sql, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get raw DB: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	b, err := bus.NewBus(sql, logger)
	if err != nil {
		t.Fatalf("failed to create bus: %v", err)
	}

	busCtx, busCancel := context.WithCancel(ctx)
	defer busCancel()

	go func() {
		if err := b.Start(busCtx); err != nil && err != context.Canceled {
			t.Logf("bus stopped: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Setup scheduler and command handlers with fake clock
	closeCmd := services.NewClosePollCommandHandler(db, b, centrifugoClient)
	scheduler := jobs.NewClosePollScheduler(closeCmd, db, fakeClock)
	createCmd := services.NewCreatePollCommandHandler(db, centrifugoClient, scheduler)

	// Start scheduler FIRST
	schedulerCtx, schedulerCancel := context.WithCancel(ctx)
	defer schedulerCancel()

	go scheduler.Start(schedulerCtx)
	time.Sleep(100 * time.Millisecond)

	userID := uuid.Must(uuid.NewV7())
	food, err := food_domain.NewFood(
		"Test Food",
		"",
		nil,
		"Delicious test food",
		10.0,
		userID,
	)
	if err != nil {
		t.Fatalf("can not create food: %v", err)
	}

	err = gorm.G[food_domain.Food](db).Create(ctx, food)
	if err != nil {
		t.Fatalf("can not save food: %v", err)
	}

	userCtx := auth.WithUser(ctx, &auth.User{
		ID: userID.String(),
	})

	strategy, err := domain.OrderMultiple.Value()
	if err != nil {
		t.Fatalf("failed to get strategy value: %v", err)
	}

	// Create a poll that expires in 2 seconds (fake clock time)
	req := services.CreatePollCommand{
		OrderDate:        fakeClock.Now().Add(24 * time.Hour),
		ScheduledCloseAt: fakeClock.Now().Add(1 * time.Hour),
		Strategy:         strategy.(string),
		FoodIDs:          []string{food.ID.String()},
	}
	res, err := createCmd.Handle(userCtx, req)
	if err != nil {
		t.Fatalf("failed to create poll: %v", err)
	}

	t.Logf("Poll created with ID %s, scheduled to close at %s", res.PollID, req.ScheduledCloseAt)

	// Advance fake clock to trigger poll closure
	fakeClock.Advance(1 * time.Hour)
	time.Sleep(500 * time.Millisecond)

	// Assert
	pollID, err := uuid.FromString(res.PollID)
	if err != nil {
		t.Fatalf("invalid poll ID: %v", err)
	}

	closedPoll, err := gorm.G[domain.Poll](db).Where("id = ?", pollID).First(ctx)
	if err != nil {
		t.Fatalf("failed to fetch poll: %v", err)
	}

	assert.NotNil(t, closedPoll.ClosedAt, "Poll should be automatically closed by scheduler after creation")
	t.Logf("Poll successfully closed at %s", closedPoll.ClosedAt)

	// Graceful shutdown
	schedulerCancel()
	busCancel()
	time.Sleep(200 * time.Millisecond)
}
