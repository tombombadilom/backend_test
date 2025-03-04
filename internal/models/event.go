package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// LiveEvent represents a live event in the system
type LiveEvent struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Rewards     string    `json:"rewards"` // JSON string
}

// NewLiveEvent creates a new LiveEvent with a generated UUID
func NewLiveEvent(title, description string, startTime, endTime time.Time, rewards string) (*LiveEvent, error) {
	// Validate rewards is valid JSON
	if rewards != "" {
		var js json.RawMessage
		if err := json.Unmarshal([]byte(rewards), &js); err != nil {
			return nil, err
		}
	}

	return &LiveEvent{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		Rewards:     rewards,
	}, nil
}

// IsActive returns true if the event is currently active
func (e *LiveEvent) IsActive() bool {
	now := time.Now()
	return now.After(e.StartTime) && now.Before(e.EndTime)
}

// Validate checks if the event data is valid
func (e *LiveEvent) Validate() error {
	if e.Title == "" {
		return ErrEmptyTitle
	}
	if e.StartTime.After(e.EndTime) {
		return ErrInvalidTimeRange
	}
	if e.Rewards != "" {
		var js json.RawMessage
		if err := json.Unmarshal([]byte(e.Rewards), &js); err != nil {
			return ErrInvalidRewardsJSON
		}
	}
	return nil
} 