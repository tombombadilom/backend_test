package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tombombadilom/liveops/internal/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser adds a new user to the database
func (r *UserRepository) CreateUser(user *models.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (id, username, role, created_at)
		VALUES (?, ?, ?, ?)
	`, user.ID.String(), user.Username, string(user.Role), user.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	var idStr, roleStr string
	var createdAt string

	err := r.db.QueryRow(`
		SELECT id, username, role, created_at
		FROM users
		WHERE id = ?
	`, id.String()).Scan(&idStr, &user.Username, &roleStr, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse UUID
	user.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in database: %w", err)
	}

	// Parse role
	user.Role = models.Role(roleStr)

	// Parse timestamp
	user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at time in database: %w", err)
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	var idStr, roleStr string
	var createdAt string

	err := r.db.QueryRow(`
		SELECT id, username, role, created_at
		FROM users
		WHERE username = ?
	`, username).Scan(&idStr, &user.Username, &roleStr, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse UUID
	user.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in database: %w", err)
	}

	// Parse role
	user.Role = models.Role(roleStr)

	// Parse timestamp
	user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at time in database: %w", err)
	}

	return &user, nil
}

// ListUsers retrieves all users
func (r *UserRepository) ListUsers() ([]*models.User, error) {
	rows, err := r.db.Query(`
		SELECT id, username, role, created_at
		FROM users
		ORDER BY username
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		var idStr, roleStr string
		var createdAt string

		if err := rows.Scan(&idStr, &user.Username, &roleStr, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		// Parse UUID
		user.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID in database: %w", err)
		}

		// Parse role
		user.Role = models.Role(roleStr)

		// Parse timestamp
		user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("invalid created_at time in database: %w", err)
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// DeleteUser removes a user by ID
func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id.String())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// APIKeyRepository handles database operations for API keys
type APIKeyRepository struct {
	db *DB
}

// NewAPIKeyRepository creates a new API key repository
func NewAPIKeyRepository(db *DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// CreateAPIKey adds a new API key to the database
func (r *APIKeyRepository) CreateAPIKey(apiKey *models.APIKey) error {
	_, err := r.db.Exec(`
		INSERT INTO api_keys (id, user_id, key, created_at, expires_at, last_used)
		VALUES (?, ?, ?, ?, ?, ?)
	`, apiKey.ID.String(), apiKey.UserID.String(), apiKey.Key, apiKey.CreatedAt, apiKey.ExpiresAt, apiKey.LastUsed)

	if err != nil {
		return fmt.Errorf("failed to create API key: %w", err)
	}

	return nil
}

// GetAPIKeyByKey retrieves an API key by its key string
func (r *APIKeyRepository) GetAPIKeyByKey(key string) (*models.APIKey, error) {
	var apiKey models.APIKey
	var idStr, userIDStr string
	var createdAt, expiresAt, lastUsed string

	err := r.db.QueryRow(`
		SELECT id, user_id, key, created_at, expires_at, last_used
		FROM api_keys
		WHERE key = ?
	`, key).Scan(&idStr, &userIDStr, &apiKey.Key, &createdAt, &expiresAt, &lastUsed)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrInvalidAPIKey
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	// Parse UUIDs
	apiKey.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid API key ID in database: %w", err)
	}

	apiKey.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in database: %w", err)
	}

	// Parse timestamps
	apiKey.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at time in database: %w", err)
	}

	apiKey.ExpiresAt, err = time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid expires_at time in database: %w", err)
	}

	apiKey.LastUsed, err = time.Parse(time.RFC3339, lastUsed)
	if err != nil {
		return nil, fmt.Errorf("invalid last_used time in database: %w", err)
	}

	return &apiKey, nil
}

// UpdateAPIKeyLastUsed updates the last_used timestamp for an API key
func (r *APIKeyRepository) UpdateAPIKeyLastUsed(id uuid.UUID, lastUsed time.Time) error {
	_, err := r.db.Exec(`
		UPDATE api_keys
		SET last_used = ?
		WHERE id = ?
	`, lastUsed, id.String())

	if err != nil {
		return fmt.Errorf("failed to update API key last_used: %w", err)
	}

	return nil
}

// DeleteAPIKey removes an API key by ID
func (r *APIKeyRepository) DeleteAPIKey(id uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM api_keys WHERE id = ?", id.String())
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

// ListAPIKeysByUserID retrieves all API keys for a user
func (r *APIKeyRepository) ListAPIKeysByUserID(userID uuid.UUID) ([]*models.APIKey, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, key, created_at, expires_at, last_used
		FROM api_keys
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query API keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []*models.APIKey

	for rows.Next() {
		var apiKey models.APIKey
		var idStr, userIDStr string
		var createdAt, expiresAt, lastUsed string

		if err := rows.Scan(&idStr, &userIDStr, &apiKey.Key, &createdAt, &expiresAt, &lastUsed); err != nil {
			return nil, fmt.Errorf("failed to scan API key row: %w", err)
		}

		// Parse UUIDs
		apiKey.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid API key ID in database: %w", err)
		}

		apiKey.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID in database: %w", err)
		}

		// Parse timestamps
		apiKey.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("invalid created_at time in database: %w", err)
		}

		apiKey.ExpiresAt, err = time.Parse(time.RFC3339, expiresAt)
		if err != nil {
			return nil, fmt.Errorf("invalid expires_at time in database: %w", err)
		}

		apiKey.LastUsed, err = time.Parse(time.RFC3339, lastUsed)
		if err != nil {
			return nil, fmt.Errorf("invalid last_used time in database: %w", err)
		}

		apiKeys = append(apiKeys, &apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating API key rows: %w", err)
	}

	return apiKeys, nil
}
