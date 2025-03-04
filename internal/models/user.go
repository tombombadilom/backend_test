package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

// Role represents user roles in the system
type Role string

const (
	// RoleAdmin has full access to all operations
	RoleAdmin Role = "admin"
	// RoleEditor can create and modify events
	RoleEditor Role = "editor"
	// RoleViewer can only view events
	RoleViewer Role = "viewer"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// APIKey represents an API key for authentication
type APIKey struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	LastUsed  time.Time `json:"last_used"`
}

// NewUser creates a new user with the specified role
func NewUser(username string, role Role) *User {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		Role:      role,
		CreatedAt: time.Now(),
	}
}

// GenerateAPIKey creates a new API key for a user
func GenerateAPIKey(userID uuid.UUID, validDays int) (*APIKey, error) {
	// Generate a random key
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	key := base64.URLEncoding.EncodeToString(b)

	now := time.Now()
	expiresAt := now.AddDate(0, 0, validDays)

	return &APIKey{
		ID:        uuid.New(),
		UserID:    userID,
		Key:       key,
		CreatedAt: now,
		ExpiresAt: expiresAt,
		LastUsed:  now,
	}, nil
}

// IsExpired checks if the API key has expired
func (k *APIKey) IsExpired() bool {
	return time.Now().After(k.ExpiresAt)
}

// UpdateLastUsed updates the last used timestamp
func (k *APIKey) UpdateLastUsed() {
	k.LastUsed = time.Now()
}

// HasPermission checks if a role has permission for the specified action
func HasPermission(role Role, action string) bool {
	switch action {
	case "read":
		// All roles can read
		return true
	case "create", "update":
		// Only admin and editor can create/update
		return role == RoleAdmin || role == RoleEditor
	case "delete":
		// Only admin can delete
		return role == RoleAdmin
	default:
		return false
	}
}
