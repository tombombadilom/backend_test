package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tombombadilom/liveops/internal/models"
)

// EventRepository handles database operations for events
type EventRepository struct {
	db *DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *DB) *EventRepository {
	return &EventRepository{db: db}
}

// Create adds a new event to the database
func (r *EventRepository) Create(event *models.LiveEvent) error {
	_, err := r.db.Exec(`
		INSERT INTO events (id, title, description, start_time, end_time, rewards, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`, event.ID.String(), event.Title, event.Description, event.StartTime, event.EndTime, event.Rewards)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

// GetByID retrieves an event by its ID
func (r *EventRepository) GetByID(id uuid.UUID) (*models.LiveEvent, error) {
	var event models.LiveEvent
	var idStr string
	var startTime, endTime string

	err := r.db.QueryRow(`
		SELECT id, title, description, start_time, end_time, rewards
		FROM events
		WHERE id = ?
	`, id.String()).Scan(&idStr, &event.Title, &event.Description, &startTime, &endTime, &event.Rewards)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrEventNotFound
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Parse UUID
	event.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID in database: %w", err)
	}

	// Parse timestamps
	event.StartTime, err = time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time in database: %w", err)
	}

	event.EndTime, err = time.Parse(time.RFC3339, endTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time in database: %w", err)
	}

	return &event, nil
}

// Update updates an existing event
func (r *EventRepository) Update(event *models.LiveEvent) error {
	result, err := r.db.Exec(`
		UPDATE events
		SET title = ?, description = ?, start_time = ?, end_time = ?, rewards = ?, updated_at = datetime('now')
		WHERE id = ?
	`, event.Title, event.Description, event.StartTime, event.EndTime, event.Rewards, event.ID.String())

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrEventNotFound
	}

	return nil
}

// Delete removes an event by its ID
func (r *EventRepository) Delete(id uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM events WHERE id = ?", id.String())
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrEventNotFound
	}

	return nil
}

// List retrieves all events, optionally filtered by active status
func (r *EventRepository) List(activeOnly bool) ([]*models.LiveEvent, error) {
	var query string
	var args []interface{}

	if activeOnly {
		query = `
			SELECT id, title, description, start_time, end_time, rewards
			FROM events
			WHERE datetime('now') BETWEEN start_time AND end_time
			ORDER BY start_time
		`
	} else {
		query = `
			SELECT id, title, description, start_time, end_time, rewards
			FROM events
			ORDER BY start_time
		`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []*models.LiveEvent

	for rows.Next() {
		var event models.LiveEvent
		var idStr string
		var startTime, endTime string

		if err := rows.Scan(&idStr, &event.Title, &event.Description, &startTime, &endTime, &event.Rewards); err != nil {
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}

		// Parse UUID
		event.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid event ID in database: %w", err)
		}

		// Parse timestamps
		event.StartTime, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start time in database: %w", err)
		}

		event.EndTime, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			return nil, fmt.Errorf("invalid end time in database: %w", err)
		}

		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event rows: %w", err)
	}

	return events, nil
}
