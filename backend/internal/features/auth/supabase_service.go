package auth

import (
	"encoding/json"
	"errors"

	"github.com/SirNacou/weeate/backend/internal/features/auth/domain"
	"github.com/supabase-community/supabase-go"
)

type SupabaseService struct {
	supabaseClient *supabase.Client
}

func NewSupabaseService(url, apiKey string) (*SupabaseService, error) {
	supabaseClient, err := supabase.NewClient(url, apiKey, &supabase.ClientOptions{})
	if err != nil {
		return nil, err
	}
	return &SupabaseService{
		supabaseClient: supabaseClient,
	}, nil
}

func (s *SupabaseService) GetUserProfileByID(userID string) (*domain.UserProfile, error) {
	var profiles []domain.UserProfile
	res, _, err := s.supabaseClient.From("user_profiles").
		Select("*", "", false).
		Eq("id", userID).
		Execute()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(res, &profiles); err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, errors.New("user profile not found")
	}
	return &profiles[0], nil
}

func (s *SupabaseService) GetUserProfilesByIDs(userIDs ...string) ([]domain.UserProfile, error) {
	var profiles []domain.UserProfile
	res, _, err := s.supabaseClient.From("user_profiles").
		Select("*", "", false).
		In("id", userIDs).
		Execute()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(res, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}
