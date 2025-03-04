# Technical Specifications

## 📑 Navigation

- 🏠 **[Project Home](../../README.md)** - Return to project root
- 📘 **[Documentation Home](../README.md)** - Main documentation
- 📋 **[Instructions](../0-Instructions/INSTRUCTIONS_Backend.md)** - Project requirements
- 🔄 **[Methodology](../1-Methodology/README.md)** - Development approach
- 🏗️ **[Architecture](../1-Design/Architecture.md)** - System design
- 💾 **[Data Models](../1-Design/DataModels.md)** - Data structures
- 📊 **[Analysis](../2-Analysis/Analysis.md)** - Requirements analysis
- 🔌 **[API Specs](./APISpecifications.md)** - API documentation
- 🔒 **[Security](./SecuritySpecifications.md)** - Security measures
- 📝 **[Implementation Plan](../4-Todo/README.md)** - Development tasks

---

## Overview

This document provides the technical specifications for the Live Ops Events System. It covers the technology stack, implementation details, deployment, and testing strategies.

## Technology Stack

### Programming Language and Framework

- **Language**: Go 1.21 or later
- **HTTP Framework**: [Gin](https://github.com/gin-gonic/gin) for HTTP API
- **gRPC Framework**: [gRPC-Go](https://github.com/grpc/grpc-go) for gRPC API
- **Database**: [SQLite](https://www.sqlite.org/index.html) with [go-sqlite3](https://github.com/mattn/go-sqlite3) driver

### Libraries and Tools

- **Configuration**: [Viper](https://github.com/spf13/viper) for configuration management
- **Logging**: [Zap](https://github.com/uber-go/zap) for structured logging
- **Testing**: [Testify](https://github.com/stretchr/testify) for testing
- **Mocking**: [gomock](https://github.com/golang/mock) for mocking interfaces
- **Database Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations
- **API Documentation**: [Swagger](https://github.com/swaggo/swag) for HTTP API documentation
- **Containerization**: [Docker](https://www.docker.com/) for containerization
- **CI/CD**: [GitHub Actions](https://github.com/features/actions) for continuous integration and deployment

## Implementation Details

### Project Structure

The project follows a standard Go project layout:

```
/
├── cmd/                  # Main applications
│   └── server/           # The server application
├── internal/             # Private application and library code
│   ├── api/              # API handlers (HTTP and gRPC)
│   │   ├── http/         # HTTP handlers
│   │   └── grpc/         # gRPC handlers
│   ├── auth/             # Authentication and authorization
│   ├── config/           # Configuration
│   ├── db/               # Database access
│   │   ├── migrations/   # Database migrations
│   │   └── sqlite/       # SQLite implementation
│   ├── models/           # Domain models
│   └── service/          # Business logic
├── pkg/                  # Public library code
│   └── proto/            # Protobuf definitions
├── scripts/              # Scripts for development and CI/CD
├── test/                 # Additional test applications and test data
└── vendor/               # Application dependencies
```

### HTTP and gRPC Server

The application runs both HTTP and gRPC servers on the same port using a multiplexer. The multiplexer determines whether an incoming request is HTTP or gRPC based on the content type and routes it to the appropriate server.

```go
func main() {
    // Create a new gRPC server
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(
            grpcMiddleware.ChainUnaryServer(
                authInterceptor.Unary(),
                grpcMiddleware.UnaryServerRecover(),
            ),
        ),
    )

    // Register the gRPC service
    pb.RegisterEventServiceServer(grpcServer, grpcHandler)

    // Create a new HTTP server
    httpServer := gin.Default()
    httpServer.Use(authMiddleware.Authenticate())

    // Register the HTTP routes
    httpServer.GET("/events", httpHandler.GetEvents)
    httpServer.GET("/events/:id", httpHandler.GetEvent)

    // Create a multiplexer
    mux := cmux.New(listener)

    // Match gRPC requests
    grpcListener := mux.MatchWithWriters(
        cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
    )

    // Match HTTP requests
    httpListener := mux.Match(cmux.Any())

    // Start the servers
    go grpcServer.Serve(grpcListener)
    go httpServer.RunListener(httpListener)

    // Start the multiplexer
    mux.Serve()
}
```

### Authentication and Authorization

The application uses API keys for authentication and role-based access control. API keys are stored in the database with a role that determines which endpoints the key can access.

```go
func (a *Authenticator) Authenticate(c *gin.Context) {
    // Get the API key from the Authorization header
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized: Missing API key",
        })
        return
    }

    // Extract the API key
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized: Invalid Authorization header format",
        })
        return
    }

    apiKey := parts[1]

    // Validate the API key
    key, err := a.service.ValidateAPIKey(c, apiKey)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized: Invalid API key",
        })
        return
    }

    // Check the role
    if key.Role != "http_user" {
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "error": "Forbidden: API key does not have the required role",
        })
        return
    }

    // Set the API key in the context
    c.Set("api_key", key)

    c.Next()
}
```

### Database Access

The application uses SQLite for persistent storage. The database is accessed through a repository layer that provides an abstraction over the database.

```go
func (r *EventRepository) GetEvents(ctx context.Context) ([]*models.LiveEvent, error) {
    query := `
        SELECT id, title, description, start_time, end_time, rewards, created_at, updated_at
        FROM live_events
        WHERE start_time <= datetime('now')
        AND end_time >= datetime('now')
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var events []*models.LiveEvent
    for rows.Next() {
        var event models.LiveEvent
        err := rows.Scan(
            &event.ID,
            &event.Title,
            &event.Description,
            &event.StartTime,
            &event.EndTime,
            &event.Rewards,
            &event.CreatedAt,
            &event.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        events = append(events, &event)
    }

    return events, nil
}
```

### Error Handling

The application uses a consistent error handling approach across both HTTP and gRPC APIs.

```go
func (h *HTTPHandler) GetEvent(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Bad Request: Missing event ID",
        })
        return
    }

    event, err := h.service.GetEvent(c, id)
    if err != nil {
        if errors.Is(err, models.ErrEventNotFound) {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Not Found: Event not found",
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Internal Server Error",
        })
        return
    }

    c.JSON(http.StatusOK, event)
}
```

### Logging

The application uses structured logging to provide detailed information about the application's behavior.

```go
func (s *Service) GetEvent(ctx context.Context, id string) (*models.LiveEvent, error) {
    logger := s.logger.With(zap.String("event_id", id))
    logger.Info("Getting event")

    event, err := s.repo.GetEvent(ctx, id)
    if err != nil {
        logger.Error("Failed to get event", zap.Error(err))
        return nil, err
    }

    logger.Info("Got event", zap.String("event_title", event.Title))
    return event, nil
}
```

## Deployment

### Docker

The application is containerized using Docker. The Dockerfile is designed to create a small, secure image.

```dockerfile
# Build stage
FROM golang:1.21-alpine AS build

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o server ./cmd/server

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install SQLite
RUN apk add --no-cache sqlite-libs

# Copy the binary from the build stage
COPY --from=build /app/server .

# Copy the migrations
COPY --from=build /app/internal/db/migrations ./internal/db/migrations

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./server"]
```

### Environment Variables

The application is configured using environment variables:

- `PORT`: The port to listen on (default: 8080)
- `DB_PATH`: The path to the SQLite database file (default: ./data/liveops.db)
- `LOG_LEVEL`: The log level (default: info)
- `API_KEY_HTTP`: The API key for HTTP access (for development only)
- `API_KEY_GRPC`: The API key for gRPC access (for development only)

### Health Checks

The application provides health check endpoints:

- `GET /health`: Returns 200 OK if the application is healthy
- `GET /health/ready`: Returns 200 OK if the application is ready to serve requests

## Testing

### Unit Tests

Unit tests are written for all components of the application. The tests use the Testify library for assertions and gomock for mocking interfaces.

```go
func TestGetEvent(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockEventRepository(ctrl)
    mockLogger := mocks.NewMockLogger(ctrl)

    service := NewService(mockRepo, mockLogger)

    event := &models.LiveEvent{
        ID:          "550e8400-e29b-41d4-a716-446655440000",
        Title:       "Summer Festival",
        Description: "Join the summer festival and earn exclusive rewards!",
        StartTime:   time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
        EndTime:     time.Date(2023, 6, 30, 23, 59, 59, 0, time.UTC),
        Rewards:     "{\"items\":[{\"id\":\"item1\",\"name\":\"Sun Hat\",\"quantity\":1},{\"id\":\"item2\",\"name\":\"Beach Ball\",\"quantity\":1}]}",
        CreatedAt:   time.Date(2023, 5, 15, 12, 0, 0, 0, time.UTC),
        UpdatedAt:   time.Date(2023, 5, 15, 12, 0, 0, 0, time.UTC),
    }

    mockRepo.EXPECT().GetEvent(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000").Return(event, nil)
    mockLogger.EXPECT().With(gomock.Any()).Return(mockLogger)
    mockLogger.EXPECT().Info(gomock.Any())
    mockLogger.EXPECT().Info(gomock.Any(), gomock.Any())

    result, err := service.GetEvent(context.Background(), "550e8400-e29b-41d4-a716-446655440000")
    assert.NoError(t, err)
    assert.Equal(t, event, result)
}
```

### Integration Tests

Integration tests are written to test the interaction between components. The tests use a real SQLite database.

```go
func TestIntegrationGetEvent(t *testing.T) {
    db, err := sql.Open("sqlite3", ":memory:")
    require.NoError(t, err)
    defer db.Close()

    // Create the schema
    _, err = db.Exec(`
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
    `)
    require.NoError(t, err)

    // Insert test data
    _, err = db.Exec(`
        INSERT INTO live_events (id, title, description, start_time, end_time, rewards, created_at, updated_at)
        VALUES (
            '550e8400-e29b-41d4-a716-446655440000',
            'Summer Festival',
            'Join the summer festival and earn exclusive rewards!',
            '2023-06-01 00:00:00',
            '2023-06-30 23:59:59',
            '{"items":[{"id":"item1","name":"Sun Hat","quantity":1},{"id":"item2","name":"Beach Ball","quantity":1}]}',
            '2023-05-15 12:00:00',
            '2023-05-15 12:00:00'
        );
    `)
    require.NoError(t, err)

    repo := sqlite.NewEventRepository(db)
    logger := zap.NewNop().Sugar()
    service := service.NewService(repo, logger)

    event, err := service.GetEvent(context.Background(), "550e8400-e29b-41d4-a716-446655440000")
    assert.NoError(t, err)
    assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", event.ID)
    assert.Equal(t, "Summer Festival", event.Title)
}
```

### End-to-End Tests

End-to-end tests are written to test the complete flow from API to database and back. The tests use a real HTTP and gRPC client.

```go
func TestE2EGetEvent(t *testing.T) {
    // Start the server
    go main()

    // Wait for the server to start
    time.Sleep(1 * time.Second)

    // Create an HTTP client
    client := &http.Client{}

    // Create a request
    req, err := http.NewRequest("GET", "http://localhost:8080/events/550e8400-e29b-41d4-a716-446655440000", nil)
    require.NoError(t, err)

    // Set the API key
    req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY_HTTP"))

    // Send the request
    resp, err := client.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    // Check the response
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // Parse the response
    var event models.LiveEvent
    err = json.NewDecoder(resp.Body).Decode(&event)
    require.NoError(t, err)

    // Check the event
    assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", event.ID)
    assert.Equal(t, "Summer Festival", event.Title)
}
```

## Performance Considerations

### Database Optimization

The application uses indexes to optimize database queries:

```sql
CREATE INDEX idx_live_events_start_time_end_time ON live_events(start_time, end_time);
```

### Caching

The application uses in-memory caching to reduce database load:

```go
func (r *EventRepository) GetEvents(ctx context.Context) ([]*models.LiveEvent, error) {
    // Check the cache
    if r.cache != nil {
        if events, ok := r.cache.Get("events"); ok {
            return events.([]*models.LiveEvent), nil
        }
    }

    // Query the database
    query := `
        SELECT id, title, description, start_time, end_time, rewards, created_at, updated_at
        FROM live_events
        WHERE start_time <= datetime('now')
        AND end_time >= datetime('now')
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var events []*models.LiveEvent
    for rows.Next() {
        var event models.LiveEvent
        err := rows.Scan(
            &event.ID,
            &event.Title,
            &event.Description,
            &event.StartTime,
            &event.EndTime,
            &event.Rewards,
            &event.CreatedAt,
            &event.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        events = append(events, &event)
    }

    // Update the cache
    if r.cache != nil {
        r.cache.Set("events", events, 5*time.Minute)
    }

    return events, nil
}
```

### Connection Pooling

The application uses connection pooling to reduce the overhead of creating new database connections:

```go
func NewDB(path string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }

    // Set connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    return db, nil
}
```

## Security Considerations

### Input Validation

The application validates all input to prevent injection attacks:

```go
func (s *Service) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
    // Validate the request
    if req.Title == "" {
        return nil, status.Error(codes.InvalidArgument, "Title is required")
    }
    if req.Description == "" {
        return nil, status.Error(codes.InvalidArgument, "Description is required")
    }
    if req.StartTime == nil {
        return nil, status.Error(codes.InvalidArgument, "Start time is required")
    }
    if req.EndTime == nil {
        return nil, status.Error(codes.InvalidArgument, "End time is required")
    }
    if req.Rewards == "" {
        return nil, status.Error(codes.InvalidArgument, "Rewards is required")
    }

    // Validate the rewards JSON
    var rewards map[string]interface{}
    if err := json.Unmarshal([]byte(req.Rewards), &rewards); err != nil {
        return nil, status.Error(codes.InvalidArgument, "Rewards must be a valid JSON string")
    }

    // Create the event
    event := &models.LiveEvent{
        ID:          uuid.New().String(),
        Title:       req.Title,
        Description: req.Description,
        StartTime:   req.StartTime.AsTime(),
        EndTime:     req.EndTime.AsTime(),
        Rewards:     req.Rewards,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    // Save the event
    if err := s.repo.CreateEvent(ctx, event); err != nil {
        return nil, status.Error(codes.Internal, "Failed to create event")
    }

    // Return the response
    return &pb.CreateEventResponse{
        Event: &pb.LiveEvent{
            Id:          event.ID,
            Title:       event.Title,
            Description: event.Description,
            StartTime:   req.StartTime,
            EndTime:     req.EndTime,
            Rewards:     event.Rewards,
            CreatedAt:   timestamppb.New(event.CreatedAt),
            UpdatedAt:   timestamppb.New(event.UpdatedAt),
        },
    }, nil
}
```

### API Key Security

The application securely stores API keys in the database:

```go
func (s *Service) CreateAPIKey(ctx context.Context, role string) (*models.APIKey, error) {
    // Generate a random API key
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return nil, err
    }
    apiKey := base64.URLEncoding.EncodeToString(key)

    // Create the API key
    key := &models.APIKey{
        ID:         uuid.New().String(),
        Key:        apiKey,
        Role:       role,
        CreatedAt:  time.Now(),
        LastUsedAt: time.Time{},
        IsActive:   true,
    }

    // Save the API key
    if err := s.repo.CreateAPIKey(ctx, key); err != nil {
        return nil, err
    }

    return key, nil
}
```

### Rate Limiting

The application implements rate limiting to prevent abuse:

```go
func (m *Middleware) RateLimit(c *gin.Context) {
    // Get the API key from the context
    apiKey, exists := c.Get("api_key")
    if !exists {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized: Missing API key",
        })
        return
    }

    // Get the rate limiter for this API key
    limiter := m.getLimiter(apiKey.(*models.APIKey).Key)

    // Check if the request is allowed
    if !limiter.Allow() {
        c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
            "error": "Too Many Requests: Rate limit exceeded",
        })
        return
    }

    c.Next()
}
```

## Conclusion

This technical specification provides a comprehensive guide to the implementation of the Live Ops Events System. It covers the technology stack, implementation details, deployment, and testing strategies. The system is designed to be secure, efficient, and easy to maintain, while meeting all the requirements specified in the project brief. 