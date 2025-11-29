# Close Poll Scheduler Integration Tests

This directory contains integration tests for the `ClosePollScheduler` component.

## Overview

The integration tests verify that the `ClosePollScheduler` correctly:
- Closes polls that are due (past their `scheduled_closes_at` time)
- Skips polls that are not yet due
- Responds to trigger updates when new polls are created
- Handles edge cases like no polls in the database
- Works correctly with multiple polls scheduled at different times

## Prerequisites

To run these integration tests, you need:

1. **PostgreSQL Database**: A running PostgreSQL instance for testing
2. **Go 1.25+**: As specified in go.mod

## Running the Tests

### Option 1: Using Default Database Configuration

The tests will use the following default values if environment variables are not set:

```bash
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=weeate_user
TEST_DB_PASSWORD=weeate_password
TEST_DB_NAME=weeate_test_db
```

Run the tests with:

```bash
cd backend
go test -v ./internal/features/polls/infrastructure/jobs/...
```

### Option 2: Using Custom Database Configuration

Set environment variables to customize the test database connection:

```bash
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=your_user
export TEST_DB_PASSWORD=your_password
export TEST_DB_NAME=weeate_test_db

cd backend
go test -v ./internal/features/polls/infrastructure/jobs/...
```

### Option 3: Run Specific Tests

To run a specific test:

```bash
cd backend
go test -v -run TestClosePollScheduler_Integration_CloseDuePolls ./internal/features/polls/infrastructure/jobs/...
```

## Test Cases

### TestClosePollScheduler_Integration_CloseDuePolls
Tests that the scheduler correctly closes polls that are past their scheduled closing time and leaves future polls open.

### TestClosePollScheduler_Integration_StartAndTrigger
Tests the scheduler's start functionality with a mock clock to verify it closes polls at the correct time.

### TestClosePollScheduler_Integration_TriggerUpdate
Tests that the scheduler can be triggered to recalculate its schedule when new polls are created.

### TestClosePollScheduler_Integration_MultiplePolls
Tests the scheduler with multiple polls scheduled at different times to verify correct behavior with batch operations.

### TestClosePollScheduler_Integration_NoPolls
Tests that the scheduler handles the case where no polls exist in the database without crashing.

## Test Database Setup

The tests will automatically:
1. Connect to the configured test database
2. Run migrations to create necessary tables
3. Clean up test data before each test
4. Clean up test data after each test completes

**Note**: If the test cannot connect to the database, it will be skipped with a message indicating the connection failure.

## Mock Clock

The tests use a `MockClock` implementation (in `backend/internal/common/infrastructure/clock/mock_clock.go`) to control time progression during tests. This allows for deterministic testing of time-based behavior without waiting for real time to pass.

## Troubleshooting

### Tests are skipped
- Ensure PostgreSQL is running and accessible
- Verify database credentials and permissions
- Check that the test database exists (it will be created if it doesn't exist and you have permissions)

### Tests fail with connection errors
- Verify firewall settings allow connection to PostgreSQL
- Check PostgreSQL logs for connection issues
- Ensure the test database user has appropriate permissions

### Tests fail with migration errors
- Ensure the test database user has CREATE TABLE permissions
- Check for conflicting schema from previous test runs
- Try dropping and recreating the test database

## CI/CD Integration

For CI/CD pipelines, you can use Docker to run a PostgreSQL instance:

```bash
docker run -d \
  --name postgres-test \
  -e POSTGRES_USER=weeate_user \
  -e POSTGRES_PASSWORD=weeate_password \
  -e POSTGRES_DB=weeate_test_db \
  -p 5432:5432 \
  postgres:16

# Wait for postgres to be ready
sleep 5

# Run tests
go test -v ./internal/features/polls/infrastructure/jobs/...

# Cleanup
docker stop postgres-test
docker rm postgres-test
```
