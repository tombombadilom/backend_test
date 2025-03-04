# API Specifications

## 📑 Navigation

- 🏠 **[Project Home](../../README.md)** - Return to project root
- 📘 **[Documentation Home](../README.md)** - Main documentation
- 📋 **[Instructions](../0-Instructions/INSTRUCTIONS_Backend.md)** - Project requirements
- 🔄 **[Methodology](../1-Methodology/README.md)** - Development approach
- 🏗️ **[Architecture](../1-Design/Architecture.md)** - System design
- 💾 **[Data Models](../1-Design/DataModels.md)** - Data structures
- 📊 **[Analysis](../2-Analysis/Analysis.md)** - Requirements analysis
- ⚙️ **[Technical Specs](./TechnicalSpecifications.md)** - Technical details
- 🔒 **[Security](./SecuritySpecifications.md)** - Security measures
- 📝 **[Implementation Plan](../4-Todo/README.md)** - Development tasks

---

## Overview

This document specifies the HTTP and gRPC APIs for the Live Ops Events System. The system exposes two APIs:

1. **HTTP API**: For public access to retrieve live events.
2. **gRPC API**: For internal administrative operations to manage live events.

Both APIs are served on the same port and use API keys for authentication and role-based access control.

## Authentication

All API requests must include an API key in the `Authorization` header. The format of the header is:

```
Authorization: Bearer <api_key>
```

API keys have roles:
- `http_user`: Can access only HTTP endpoints.
- `grpc_admin`: Can access only gRPC endpoints.

## HTTP API

### Endpoints

#### GET /events

Retrieves all active events.

**Authentication**: Requires API Key with `http_user` role.

**Request**:
```http
GET /events HTTP/1.1
Host: example.com
Authorization: Bearer <api_key>
```

**Response**:
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "events": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Summer Festival",
      "description": "Join the summer festival and earn exclusive rewards!",
      "start_time": "2023-06-01T00:00:00Z",
      "end_time": "2023-06-30T23:59:59Z",
      "rewards": "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "title": "Winter Wonderland",
      "description": "Experience the magic of winter and collect special items!",
      "start_time": "2023-12-01T00:00:00Z",
      "end_time": "2023-12-31T23:59:59Z",
      "rewards": "{\"items\":[{\"id\":\"item3\",\"name\":\"Snow Globe\",\"quantity\":1},{\"id\":\"item4\",\"name\":\"Winter Hat\",\"quantity\":1}]}"
    }
  ]
}
```

**Error Responses**:

- **401 Unauthorized**: Invalid or missing API key.
- **403 Forbidden**: API key does not have the required role.
- **500 Internal Server Error**: Server error.

#### GET /events/{id}

Retrieves details of a specific event.

**Authentication**: Requires API Key with `http_user` role.

**Request**:
```http
GET /events/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
Host: example.com
Authorization: Bearer <api_key>
```

**Response**:
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Summer Festival",
  "description": "Join the summer festival and earn exclusive rewards!",
  "start_time": "2023-06-01T00:00:00Z",
  "end_time": "2023-06-30T23:59:59Z",
  "rewards": "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}"
}
```

**Error Responses**:

- **401 Unauthorized**: Invalid or missing API key.
- **403 Forbidden**: API key does not have the required role.
- **404 Not Found**: Event not found.
- **500 Internal Server Error**: Server error.

## gRPC API

### Service Definition

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

### Methods

#### CreateEvent

Creates a new live event.

**Authentication**: Requires API Key with `grpc_admin` role.

**Request**:
```protobuf
{
  "title": "Summer Festival",
  "description": "Join the summer festival and earn exclusive rewards!",
  "start_time": "2023-06-01T00:00:00Z",
  "end_time": "2023-06-30T23:59:59Z",
  "rewards": "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}"
}
```

**Response**:
```protobuf
{
  "event": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Summer Festival",
    "description": "Join the summer festival and earn exclusive rewards!",
    "start_time": "2023-06-01T00:00:00Z",
    "end_time": "2023-06-30T23:59:59Z",
    "rewards": "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}",
    "created_at": "2023-05-15T12:00:00Z",
    "updated_at": "2023-05-15T12:00:00Z"
  }
}
```

**Error Responses**:

- **UNAUTHENTICATED**: Invalid or missing API key.
- **PERMISSION_DENIED**: API key does not have the required role.
- **INVALID_ARGUMENT**: Invalid request parameters.
- **INTERNAL**: Server error.

#### UpdateEvent

Updates an existing event.

**Authentication**: Requires API Key with `grpc_admin` role.

**Request**:
```protobuf
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Updated Summer Festival",
  "description": "Join the updated summer festival and earn exclusive rewards!",
  "start_time": "2023-06-01T00:00:00Z",
  "end_time": "2023-06-30T23:59:59Z",
  "rewards": "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}"
}
```

**Response**:
```protobuf
{
  "event": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Updated Summer Festival",
    "description": "Join the updated summer festival and earn exclusive rewards!",
    "start_time": "2023-06-01T00:00:00Z",
    "end_time": "2023-06-30T23:59:59Z",
    "rewards": "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}",
    "created_at": "2023-05-15T12:00:00Z",
    "updated_at": "2023-05-15T12:30:00Z"
  }
}
```

**Error Responses**:

- **UNAUTHENTICATED**: Invalid or missing API key.
- **PERMISSION_DENIED**: API key does not have the required role.
- **INVALID_ARGUMENT**: Invalid request parameters.
- **NOT_FOUND**: Event not found.
- **INTERNAL**: Server error.

#### DeleteEvent

Deletes an event.

**Authentication**: Requires API Key with `grpc_admin` role.

**Request**:
```protobuf
{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response**:
```protobuf
{
  "success": true
}
```

**Error Responses**:

- **UNAUTHENTICATED**: Invalid or missing API key.
- **PERMISSION_DENIED**: API key does not have the required role.
- **INVALID_ARGUMENT**: Invalid request parameters.
- **NOT_FOUND**: Event not found.
- **INTERNAL**: Server error.

#### ListEvents

Retrieves all events (including past events).

**Authentication**: Requires API Key with `grpc_admin` role.

**Request**:
```protobuf
{}
```

**Response**:
```protobuf
{
  "events": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Summer Festival",
      "description": "Join the summer festival and earn exclusive rewards!",
      "start_time": "2023-06-01T00:00:00Z",
      "end_time": "2023-06-30T23:59:59Z",
      "rewards": "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}",
      "created_at": "2023-05-15T12:00:00Z",
      "updated_at": "2023-05-15T12:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "title": "Winter Wonderland",
      "description": "Experience the magic of winter and collect special items!",
      "start_time": "2023-12-01T00:00:00Z",
      "end_time": "2023-12-31T23:59:59Z",
      "rewards": "{\"items\":[{\"id\":\"item3\",\"name\":\"Snow Globe\",\"quantity\":1},{\"id\":\"item4\",\"name\":\"Winter Hat\",\"quantity\":1}]}",
      "created_at": "2023-05-15T12:00:00Z",
      "updated_at": "2023-05-15T12:00:00Z"
    }
  ]
}
```

**Error Responses**:

- **UNAUTHENTICATED**: Invalid or missing API key.
- **PERMISSION_DENIED**: API key does not have the required role.
- **INTERNAL**: Server error.

## Error Handling

### HTTP API

HTTP errors are returned with appropriate status codes and a JSON body:

```json
{
  "error": {
    "code": 401,
    "message": "Unauthorized: Invalid API key"
  }
}
```

### gRPC API

gRPC errors are returned with appropriate status codes and error messages:

```
{
  "code": 16,
  "message": "Unauthenticated: Invalid API key"
}
```

## Rate Limiting

To prevent abuse, the API implements rate limiting:

- HTTP API: 100 requests per minute per API key.
- gRPC API: 50 requests per minute per API key.

Exceeding these limits will result in a `429 Too Many Requests` (HTTP) or `RESOURCE_EXHAUSTED` (gRPC) error.

## Versioning

The API is versioned to ensure backward compatibility:

- HTTP API: Version is included in the URL path (e.g., `/v1/events`).
- gRPC API: Version is included in the package name (e.g., `events.v1`).

## Conclusion

This API specification provides a comprehensive guide to the HTTP and gRPC APIs exposed by the Live Ops Events System. It covers authentication, endpoints, request/response formats, error handling, rate limiting, and versioning. The APIs are designed to be secure, efficient, and easy to use, while meeting all the requirements specified in the project brief. 