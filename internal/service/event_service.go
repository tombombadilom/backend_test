package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tombombadilom/liveops/internal/db"
	"github.com/tombombadilom/liveops/internal/models"
)

// EventService handles business logic for events
type EventService struct {
	eventRepo *db.EventRepository
}

// NewEventService creates a new event service
func NewEventService(eventRepo *db.EventRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
	}
}

// CreateEvent creates a new event
func (s *EventService) CreateEvent(title, description string, startTime, endTime time.Time, rewards string) (*models.LiveEvent, error) {
	// Create new event
	event, err := models.NewLiveEvent(title, description, startTime, endTime, rewards)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Validate event
	if err := event.Validate(); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.eventRepo.Create(event); err != nil {
		return nil, fmt.Errorf("failed to save event: %w", err)
	}

	return event, nil
}

// GetEvent retrieves an event by ID
func (s *EventService) GetEvent(id string) (*models.LiveEvent, error) {
	// Parse UUID
	eventID, err := uuid.Parse(id)
	if err != nil {
		return nil, models.ErrInvalidID
	}

	// Get from database
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// UpdateEvent updates an existing event
func (s *EventService) UpdateEvent(id, title, description string, startTime, endTime time.Time, rewards string) (*models.LiveEvent, error) {
	// Parse UUID
	eventID, err := uuid.Parse(id)
	if err != nil {
		return nil, models.ErrInvalidID
	}

	// Get existing event
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, err
	}

	// Update fields
	event.Title = title
	event.Description = description
	event.StartTime = startTime
	event.EndTime = endTime
	event.Rewards = rewards

	// Validate event
	if err := event.Validate(); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.eventRepo.Update(event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return event, nil
}

// DeleteEvent removes an event by ID
func (s *EventService) DeleteEvent(id string) error {
	// Parse UUID
	eventID, err := uuid.Parse(id)
	if err != nil {
		return models.ErrInvalidID
	}

	// Delete from database
	if err := s.eventRepo.Delete(eventID); err != nil {
		return err
	}

	return nil
}

// ListEvents retrieves all events, optionally filtered by active status
func (s *EventService) ListEvents(activeOnly bool) ([]*models.LiveEvent, error) {
	// Get from database
	events, err := s.eventRepo.List(activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	return events, nil
}
