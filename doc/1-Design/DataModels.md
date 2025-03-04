# Data Models

## 📑 Navigation

- 🏠 **[Project Home](../../README.md)** - Return to project root
- 📘 **[Documentation Home](../README.md)** - Main documentation
- 📋 **[Instructions](../0-Instructions/INSTRUCTIONS_Backend.md)** - Project requirements
- 🔄 **[Methodology](../1-Methodology/README.md)** - Development approach
- 🏗️ **[Architecture](./Architecture.md)** - System design
- 📊 **[Analysis](../2-Analysis/Analysis.md)** - Requirements analysis
- ⚙️ **[Technical Specs](../3-Specifications/TechnicalSpecifications.md)** - Technical details
- 🔌 **[API Specs](../3-Specifications/APISpecifications.md)** - API documentation
- 🔒 **[Security](../3-Specifications/SecuritySpecifications.md)** - Security measures
- 📝 **[Implementation Plan](../4-Todo/README.md)** - Development tasks

---

## Overview

This document describes the data models used in the Live Ops Events System. The system uses a SQLite database to store data about live events and API keys.

## Entity Relationship Diagram

```
┌───────────────┐       ┌───────────────┐
│   LiveEvent   │       │    APIKey     │
├───────────────┤       ├───────────────┤
│ id            │       │ id            │
│ title         │       │ key           │
│ description   │       │ role          │
│ start_time    │       │ created_at    │
│ end_time      │       │ last_used_at  │
│ rewards       │       │ is_active     │
│ created_at    │       └───────────────┘
│ updated_at    │
└───────────────┘
```

## Data Models

### LiveEvent

The `LiveEvent` model represents a live event in the game.

#### Go Struct

```go
type LiveEvent struct {
    ID          string    `json:"id" db:"id"`
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    StartTime   time.Time `json:"start_time" db:"start_time"`
    EndTime     time.Time `json:"end_time" db:"end_time"`
    Rewards     string    `json:"rewards" db:"rewards"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

#### SQLite Schema

```sql
CREATE TABLE live_events (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    rewards TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

#### Field Descriptions

- `id`: A UUID that uniquely identifies the event.
- `title`: The title of the event.
- `description`: A detailed description of the event.
- `start_time`: The time when the event starts.
- `end_time`: The time when the event ends.
- `rewards`: A JSON string describing the rewards for the event.
- `created_at`: The time when the event was created.
- `updated_at`: The time when the event was last updated.

### APIKey

The `APIKey` model represents an API key used for authentication.

#### Go Struct

```go
type APIKey struct {
    ID         string    `json:"id" db:"id"`
    Key        string    `json:"key" db:"key"`
    Role       string    `json:"role" db:"role"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    LastUsedAt time.Time `json:"last_used_at" db:"last_used_at"`
    IsActive   bool      `json:"is_active" db:"is_active"`
}
```

#### SQLite Schema

```sql
CREATE TABLE api_keys (
    id TEXT PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);
```

#### Field Descriptions

- `id`: A UUID that uniquely identifies the API key.
- `key`: The actual API key string.
- `role`: The role associated with the API key (e.g., "http_user" or "grpc_admin").
- `created_at`: The time when the API key was created.
- `last_used_at`: The time when the API key was last used.
- `is_active`: A boolean indicating whether the API key is active.

## Protobuf Definitions

### LiveEvent

```protobuf
syntax = "proto3";

package events;

import "google/protobuf/timestamp.proto";

option go_package = "github.com/yourusername/liveops/pkg/proto/events";

message LiveEvent {
    string id = 1;
    string title = 2;
    string description = 3;
    google.protobuf.Timestamp start_time = 4;
    google.protobuf.Timestamp end_time = 5;
    string rewards = 6;
    google.protobuf.Timestamp created_at = 7;
    google.protobuf.Timestamp updated_at = 8;
}

message CreateEventRequest {
    string title = 1;
    string description = 2;
    google.protobuf.Timestamp start_time = 3;
    google.protobuf.Timestamp end_time = 4;
    string rewards = 5;
}

message CreateEventResponse {
    LiveEvent event = 1;
}

message UpdateEventRequest {
    string id = 1;
    string title = 2;
    string description = 3;
    google.protobuf.Timestamp start_time = 4;
    google.protobuf.Timestamp end_time = 5;
    string rewards = 6;
}

message UpdateEventResponse {
    LiveEvent event = 1;
}

message DeleteEventRequest {
    string id = 1;
}

message DeleteEventResponse {
    bool success = 1;
}

message ListEventsRequest {
    // Empty request
}

message ListEventsResponse {
    repeated LiveEvent events = 1;
}

service EventService {
    rpc CreateEvent(CreateEventRequest) returns (CreateEventResponse);
    rpc UpdateEvent(UpdateEventRequest) returns (UpdateEventResponse);
    rpc DeleteEvent(DeleteEventRequest) returns (DeleteEventResponse);
    rpc ListEvents(ListEventsRequest) returns (ListEventsResponse);
}
```

## JSON Representations

### LiveEvent

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Summer Festival",
  "description": "Join the summer festival and earn exclusive rewards!",
  "start_time": "2023-06-01T00:00:00Z",
  "end_time": "2023-06-30T23:59:59Z",
  "rewards": "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}",
  "created_at": "2023-05-15T12:00:00Z",
  "updated_at": "2023-05-15T12:00:00Z"
}
```

### APIKey

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "key": "api_key_123456789",
  "role": "http_user",
  "created_at": "2023-05-15T12:00:00Z",
  "last_used_at": "2023-05-15T12:00:00Z",
  "is_active": true
}
```

## Data Validation Rules

### LiveEvent

- `title`: Required, maximum length 100 characters.
- `description`: Required, maximum length 1000 characters.
- `start_time`: Required, must be a valid timestamp.
- `end_time`: Required, must be a valid timestamp and must be after `start_time`.
- `rewards`: Required, must be a valid JSON string.

### APIKey

- `key`: Required, must be unique, minimum length 16 characters.
- `role`: Required, must be either "http_user" or "grpc_admin".
- `is_active`: Required, must be a boolean.

## Indexes

### LiveEvent

- Primary key on `id`.
- Index on `start_time` and `end_time` for efficient querying of active events.

### APIKey

- Primary key on `id`.
- Unique index on `key` for efficient lookup during authentication.
- Index on `role` for efficient filtering by role.

## Conclusion

These data models provide a solid foundation for the Live Ops Events System, with clear definitions of the entities and their relationships. The models are designed to be efficient, secure, and easy to work with, while meeting all the requirements specified in the project brief. 