package models

import "errors"

// Common errors for models
var (
	ErrEmptyTitle         = errors.New("title cannot be empty")
	ErrInvalidTimeRange   = errors.New("start time must be before end time")
	ErrInvalidRewardsJSON = errors.New("rewards must be valid JSON")
	ErrEventNotFound      = errors.New("event not found")
	ErrInvalidID          = errors.New("invalid ID format")
	ErrInvalidAPIKey      = errors.New("invalid API key")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("forbidden action")
)
