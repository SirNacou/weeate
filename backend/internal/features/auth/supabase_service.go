package auth

import (
	"encoding/json"
	"errors"

	"github.com/supabase-community/supabase-go"
)

type SupabaseService struct {
	supabaseClient *supabase.Client
}

func NewSupabaseService(client *supabase.Client) *SupabaseService {
	return &SupabaseService{
		supabaseClient: client,
	}
}

func (s *SupabaseService) GetUserProfileByID(userID string) (*UserProfile, error) {
	var profiles []UserProfile
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

func (s *SupabaseService) GetUserProfilesByIDs(userIDs []string) ([]UserProfile, error) {
	var profiles []UserProfile
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
