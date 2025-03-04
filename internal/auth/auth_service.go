package auth

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/tombombadilom/liveops/internal/db"
	"github.com/tombombadilom/liveops/internal/models"
)

// AuthService handles authentication and authorization
type AuthService struct {
	userRepo   *db.UserRepository
	apiKeyRepo *db.APIKeyRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *db.UserRepository, apiKeyRepo *db.APIKeyRepository) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		apiKeyRepo: apiKeyRepo,
	}
}

// AuthenticateAPIKey validates an API key and returns the associated user
func (s *AuthService) AuthenticateAPIKey(apiKey string) (*models.User, error) {
	// Validate API key format
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, models.ErrInvalidAPIKey
	}

	// Get API key from database
	key, err := s.apiKeyRepo.GetAPIKeyByKey(apiKey)
	if err != nil {
		return nil, models.ErrInvalidAPIKey
	}

	// Check if key is expired
	if key.IsExpired() {
		return nil, models.ErrInvalidAPIKey
	}

	// Update last used timestamp
	key.UpdateLastUsed()
	if err := s.apiKeyRepo.UpdateAPIKeyLastUsed(key.ID, key.LastUsed); err != nil {
		// Log error but continue (non-critical)
	}

	// Get user associated with the API key
	user, err := s.userRepo.GetUserByID(key.UserID)
	if err != nil {
		return nil, models.ErrUnauthorized
	}

	return user, nil
}

// CheckPermission checks if a user has permission for an action
func (s *AuthService) CheckPermission(user *models.User, action string) error {
	if user == nil {
		return models.ErrUnauthorized
	}

	if !models.HasPermission(user.Role, action) {
		return models.ErrForbidden
	}

	return nil
}

// CreateAPIKey creates a new API key for a user
func (s *AuthService) CreateAPIKey(userID string, validDays int) (*models.APIKey, error) {
	// Parse UUID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, models.ErrInvalidID
	}

	// Check if user exists
	user, err := s.userRepo.GetUserByID(uid)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate API key
	apiKey, err := models.GenerateAPIKey(user.ID, validDays)
	if err != nil {
		return nil, err
	}

	// Save to database
	if err := s.apiKeyRepo.CreateAPIKey(apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

// RevokeAPIKey revokes an API key
func (s *AuthService) RevokeAPIKey(apiKeyID string) error {
	// Parse UUID
	id, err := uuid.Parse(apiKeyID)
	if err != nil {
		return models.ErrInvalidID
	}

	// Delete from database
	return s.apiKeyRepo.DeleteAPIKey(id)
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(username string, role models.Role) (*models.User, error) {
	// Validate username
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	// Check if username already exists
	_, err := s.userRepo.GetUserByUsername(username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	// Create user
	user := models.NewUser(username, role)

	// Save to database
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *AuthService) GetUser(userID string) (*models.User, error) {
	// Parse UUID
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, models.ErrInvalidID
	}

	// Get from database
	return s.userRepo.GetUserByID(id)
}

// ListUsers retrieves all users
func (s *AuthService) ListUsers() ([]*models.User, error) {
	return s.userRepo.ListUsers()
}

// ListAPIKeys retrieves all API keys for a user
func (s *AuthService) ListAPIKeys(userID string) ([]*models.APIKey, error) {
	// Parse UUID
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, models.ErrInvalidID
	}

	// Get from database
	return s.apiKeyRepo.ListAPIKeysByUserID(id)
}
