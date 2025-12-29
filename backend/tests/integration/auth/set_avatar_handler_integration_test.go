//go:build integration

package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/SirNacou/weeate/backend/internal/features/auth/domain"
	"github.com/SirNacou/weeate/backend/internal/features/auth/services"
	"github.com/SirNacou/weeate/backend/tests/testhelpers"
	"github.com/gofrs/uuid/v5"
	_ "github.com/lib/pq" // postgres driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDB *sql.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Suppress library logs during tests
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	// Setup test database with user_profiles table
	db, cleanup, err := testhelpers.SetupSupabaseTestDB(ctx)
	if err != nil {
		panic(err)
	}
	testDB = db

	exitCode := m.Run()

	cleanup()
	os.Exit(exitCode)
}

// mockSupabaseService implements the same interface as SupabaseService but uses direct SQL
type mockSupabaseService struct {
	db *sql.DB
}

func (m *mockSupabaseService) GetUserProfileByID(userID string) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	err := m.db.QueryRow(
		`SELECT id, avatar_url, display_name, created_at FROM user_profiles WHERE id = $1`,
		userID,
	).Scan(&profile.ID, &profile.AvatarURL, &profile.DisplayName, &profile.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user profile not found")
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (m *mockSupabaseService) GetUserProfilesByIDs(userIDs ...string) ([]domain.UserProfile, error) {
	// Convert userIDs to JSON array for ANY query
	userIDsJSON, _ := json.Marshal(userIDs)

	rows, err := m.db.Query(
		`SELECT id, avatar_url, display_name, created_at FROM user_profiles WHERE id = ANY($1)`,
		userIDsJSON,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []domain.UserProfile
	for rows.Next() {
		var profile domain.UserProfile
		if err := rows.Scan(&profile.ID, &profile.AvatarURL, &profile.DisplayName, &profile.CreatedAt); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, rows.Err()
}

func (m *mockSupabaseService) UpdateUserProfile(userID string, avatarURL string) error {
	result, err := m.db.Exec(
		`UPDATE user_profiles SET avatar_url = $1 WHERE id = $2`,
		avatarURL, userID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user profile not found")
	}
	return nil
}

func TestSetAvatarHandler_Integration(t *testing.T) {
	t.Run("successfully updates user avatar", func(t *testing.T) {
		ctx := context.Background()

		// Arrange - Create a test user
		userID := uuid.Must(uuid.NewV7())
		testUser := &domain.User{
			ID: userID.String(),
		}

		// Create a user profile in test database
		initialAvatarURL := "https://example.com/avatar1.jpg"
		err := createTestUserProfile(t, ctx, userID, initialAvatarURL)
		require.NoError(t, err, "Failed to create test user profile")
		defer cleanupTestUserProfile(t, ctx, userID)

		// Add user to context
		userCtx := domain.WithUser(ctx, testUser)

		// Create handler with mock Supabase service
		mockSupabase := &mockSupabaseService{db: testDB}
		handler := services.NewSetAvatarHandler(mockSupabase)

		// Act - Update avatar
		newAvatarURL := "https://example.com/avatar2.jpg"
		req := &services.SetAvatarRequest{
			AvatarURL: newAvatarURL,
		}
		err = handler.Handle(userCtx, req)

		// Assert
		require.NoError(t, err, "Failed to update avatar")

		// Verify the update in database
		profile, err := mockSupabase.GetUserProfileByID(userID.String())
		require.NoError(t, err, "Failed to fetch updated profile")
		assert.Equal(t, newAvatarURL, profile.AvatarURL, "Avatar URL should be updated")
	})

	t.Run("returns error when user not in context", func(t *testing.T) {
		ctx := context.Background()

		mockSupabase := &mockSupabaseService{db: testDB}
		handler := services.NewSetAvatarHandler(mockSupabase)

		// Act
		req := &services.SetAvatarRequest{
			AvatarURL: "https://example.com/avatar.jpg",
		}
		err := handler.Handle(ctx, req)

		// Assert
		require.Error(t, err)
		assert.Equal(t, domain.ErrUserNotFoundInContext, err)
	})

	t.Run("returns error when user profile does not exist", func(t *testing.T) {
		ctx := context.Background()

		// Arrange - Create a user that doesn't have a profile
		userID := uuid.Must(uuid.NewV7())
		testUser := &domain.User{
			ID: userID.String(),
		}
		userCtx := domain.WithUser(ctx, testUser)

		// Create handler
		mockSupabase := &mockSupabaseService{db: testDB}
		handler := services.NewSetAvatarHandler(mockSupabase)

		// Act
		req := &services.SetAvatarRequest{
			AvatarURL: "https://example.com/avatar.jpg",
		}
		err := handler.Handle(userCtx, req)

		// Assert
		require.Error(t, err, "Should return error for non-existent user profile")
	})

	t.Run("handles empty avatar URL", func(t *testing.T) {
		ctx := context.Background()

		// Arrange
		userID := uuid.Must(uuid.NewV7())
		testUser := &domain.User{
			ID: userID.String(),
		}

		// Create a user profile in test database
		initialAvatarURL := "https://example.com/avatar1.jpg"
		err := createTestUserProfile(t, ctx, userID, initialAvatarURL)
		require.NoError(t, err, "Failed to create test user profile")
		defer cleanupTestUserProfile(t, ctx, userID)

		userCtx := domain.WithUser(ctx, testUser)

		// Create handler
		mockSupabase := &mockSupabaseService{db: testDB}
		handler := services.NewSetAvatarHandler(mockSupabase)

		// Act - Update with empty URL (clearing avatar)
		req := &services.SetAvatarRequest{
			AvatarURL: "",
		}
		err = handler.Handle(userCtx, req)

		// Assert
		require.NoError(t, err, "Should handle empty avatar URL")

		// Verify the update in database
		profile, err := mockSupabase.GetUserProfileByID(userID.String())
		require.NoError(t, err, "Failed to fetch updated profile")
		assert.Equal(t, "", profile.AvatarURL, "Avatar URL should be empty")
	})
}

// Helper functions

func createTestUserProfile(t *testing.T, ctx context.Context, userID uuid.UUID, avatarURL string) error {
	t.Helper()

	_, err := testDB.ExecContext(ctx,
		`INSERT INTO user_profiles (id, avatar_url, display_name, created_at) 
		 VALUES ($1, $2, $3, $4)`,
		userID, avatarURL, "Test User", time.Now())
	return err
}

func cleanupTestUserProfile(t *testing.T, ctx context.Context, userID uuid.UUID) {
	t.Helper()

	_, err := testDB.ExecContext(ctx, `DELETE FROM user_profiles WHERE id = $1`, userID)
	if err != nil {
		t.Logf("Warning: failed to cleanup test user profile: %v", err)
	}
}
