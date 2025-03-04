# Security Specifications

## 📑 Navigation

- 🏠 **[Project Home](../../README.md)** - Return to project root
- 📘 **[Documentation Home](../README.md)** - Main documentation
- 📋 **[Instructions](../0-Instructions/INSTRUCTIONS_Backend.md)** - Project requirements
- 🔄 **[Methodology](../1-Methodology/README.md)** - Development approach
- 🏗️ **[Architecture](../1-Design/Architecture.md)** - System design
- 💾 **[Data Models](../1-Design/DataModels.md)** - Data structures
- 📊 **[Analysis](../2-Analysis/Analysis.md)** - Requirements analysis
- ⚙️ **[Technical Specs](./TechnicalSpecifications.md)** - Technical details
- 🔌 **[API Specs](./APISpecifications.md)** - API documentation
- 📝 **[Implementation Plan](../4-Todo/README.md)** - Development tasks

---

## Overview

This document outlines the security measures implemented in the Live Ops Events System. Security is a critical aspect of the system, as it handles sensitive game data and provides administrative capabilities that could impact player experience.

## Authentication and Authorization

### API Key Authentication

The system uses API keys for authentication, with a robust implementation:

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

#### Key Generation

- API keys are generated using cryptographically secure random number generators
- Keys are at least 32 bytes (256 bits) in length and base64url-encoded
- Each key is assigned a unique ID for tracking and management

```go
func GenerateAPIKey() (string, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(key), nil
}
```

#### Key Storage

- API keys are stored securely in the database
- The system does not store plaintext keys in logs or configuration files
- Keys are indexed for efficient lookup during authentication

#### Key Rotation and Revocation

- The system supports key rotation through the creation of new keys and deactivation of old ones
- Keys can be immediately revoked by setting `IsActive` to `false`
- The `LastUsedAt` field tracks key usage for auditing and cleanup

### Role-Based Access Control (RBAC)

The system implements strict role-based access control:

- `http_user`: Can only access HTTP endpoints for retrieving event data
- `grpc_admin`: Can only access gRPC endpoints for managing events

```go
func (a *Authenticator) CheckRole(role string, requiredRole string) bool {
    return role == requiredRole
}
```

#### Authorization Flow

1. Extract API key from the request header
2. Validate the API key against the database
3. Check if the key is active
4. Verify the key's role against the required role for the endpoint
5. Allow or deny access based on the role check

## Data Protection

### Input Validation

All input data is validated to prevent injection attacks and ensure data integrity:

```go
func ValidateEvent(event *models.LiveEvent) error {
    if event.Title == "" {
        return errors.New("title is required")
    }
    if len(event.Title) > 100 {
        return errors.New("title must be less than 100 characters")
    }
    if event.Description == "" {
        return errors.New("description is required")
    }
    if len(event.Description) > 1000 {
        return errors.New("description must be less than 1000 characters")
    }
    if event.StartTime.IsZero() {
        return errors.New("start time is required")
    }
    if event.EndTime.IsZero() {
        return errors.New("end time is required")
    }
    if event.EndTime.Before(event.StartTime) {
        return errors.New("end time must be after start time")
    }
    if event.Rewards == "" {
        return errors.New("rewards is required")
    }
    
    // Validate rewards JSON
    var rewards map[string]interface{}
    if err := json.Unmarshal([]byte(event.Rewards), &rewards); err != nil {
        return errors.New("rewards must be a valid JSON string")
    }
    
    return nil
}
```

### Database Security

The system implements several measures to secure the SQLite database:

- **Prepared Statements**: All database queries use prepared statements to prevent SQL injection
- **Input Sanitization**: All user input is sanitized before being used in database operations
- **Minimal Privileges**: The application uses the minimum required database privileges
- **Database File Security**: The SQLite database file has restricted permissions
- **Transaction Management**: Critical operations use transactions to maintain data integrity

```go
func (r *EventRepository) CreateEvent(ctx context.Context, event *models.LiveEvent) error {
    query := `
        INSERT INTO live_events (id, title, description, start_time, end_time, rewards, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    _, err := r.db.ExecContext(
        ctx,
        query,
        event.ID,
        event.Title,
        event.Description,
        event.StartTime,
        event.EndTime,
        event.Rewards,
        event.CreatedAt,
        event.UpdatedAt,
    )
    
    return err
}
```

## Communication Security

### Transport Layer Security

- All API communications (both HTTP and gRPC) should be secured with TLS in production
- The system supports configurable TLS settings, including cipher suites and certificate paths
- Certificate validation is enforced for all TLS connections

```go
func setupTLS() (*tls.Config, error) {
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        return nil, err
    }
    
    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        MinVersion:   tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
    }, nil
}
```

### API Security Headers

For HTTP endpoints, the system implements security headers:

```go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Next()
    }
}
```

## Rate Limiting and DoS Protection

### Rate Limiting

The system implements rate limiting to prevent abuse and ensure fair resource allocation:

- HTTP API: 100 requests per minute per API key
- gRPC API: 50 requests per minute per API key

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

### DoS Protection

Additional measures to protect against Denial of Service attacks:

- **Request Timeout**: All requests have a configurable timeout
- **Request Size Limiting**: Maximum request size is enforced
- **Connection Limiting**: Maximum number of concurrent connections is enforced
- **Resource Monitoring**: System resources are monitored to detect abnormal usage patterns

## Error Handling and Logging

### Secure Error Handling

The system implements secure error handling to prevent information leakage:

- Detailed error messages are logged internally but not exposed to clients
- Generic error messages are returned to clients
- Error codes are used to provide context without revealing implementation details

```go
func handleError(err error, c *gin.Context, logger *zap.Logger) {
    // Log the detailed error internally
    logger.Error("Error occurred", zap.Error(err))
    
    // Return a generic error to the client
    switch {
    case errors.Is(err, models.ErrNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
    case errors.Is(err, models.ErrInvalidInput):
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
    case errors.Is(err, models.ErrUnauthorized):
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    case errors.Is(err, models.ErrForbidden):
        c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
    }
}
```

### Security Logging and Monitoring

The system implements comprehensive logging for security events:

- All authentication attempts (successful and failed) are logged
- All administrative actions are logged
- All security-related errors are logged
- Logs include contextual information but exclude sensitive data

```go
func (a *Authenticator) Authenticate(c *gin.Context) {
    // Get the API key from the Authorization header
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        a.logger.Warn("Authentication attempt with missing API key",
            zap.String("ip", c.ClientIP()),
            zap.String("user_agent", c.Request.UserAgent()))
        
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized: Missing API key",
        })
        return
    }
    
    // Extract the API key
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        a.logger.Warn("Authentication attempt with invalid Authorization header format",
            zap.String("ip", c.ClientIP()),
            zap.String("user_agent", c.Request.UserAgent()))
        
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized: Invalid Authorization header format",
        })
        return
    }
    
    apiKey := parts[1]
    
    // Validate the API key
    key, err := a.service.ValidateAPIKey(c, apiKey)
    if err != nil {
        a.logger.Warn("Authentication attempt with invalid API key",
            zap.String("ip", c.ClientIP()),
            zap.String("user_agent", c.Request.UserAgent()))
        
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized: Invalid API key",
        })
        return
    }
    
    a.logger.Info("Successful authentication",
        zap.String("key_id", key.ID),
        zap.String("role", key.Role),
        zap.String("ip", c.ClientIP()))
    
    // Set the API key in the context
    c.Set("api_key", key)
    
    c.Next()
}
```

## Dependency Management and Vulnerability Scanning

### Dependency Management

The system implements secure dependency management practices:

- All dependencies are explicitly versioned
- Dependencies are regularly updated to include security patches
- A minimal set of dependencies is used to reduce the attack surface

### Vulnerability Scanning

The development process includes regular vulnerability scanning:

- Static code analysis to detect potential security issues
- Dependency scanning to detect known vulnerabilities
- Container scanning for Docker images
- Regular security audits of the codebase

## Deployment Security

### Secure Docker Configuration

The Docker deployment is configured with security in mind:

- Non-root user for running the application
- Minimal base image to reduce attack surface
- Read-only file system where possible
- No unnecessary capabilities
- Resource limits to prevent resource exhaustion

```dockerfile
# Final stage
FROM alpine:3.18

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Install SQLite
RUN apk add --no-cache sqlite-libs

# Copy the binary from the build stage
COPY --from=build /app/server .

# Copy the migrations
COPY --from=build /app/internal/db/migrations ./internal/db/migrations

# Set permissions
RUN chown -R appuser:appgroup /app && \
    chmod -R 755 /app

# Use the non-root user
USER appuser

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./server"]
```

### Environment Configuration

Secure environment configuration practices:

- Sensitive configuration is passed through environment variables, not command-line arguments
- Default values are secure
- Configuration validation to prevent insecure settings
- Secrets management for production deployments

## Security Testing

### Automated Security Testing

The CI/CD pipeline includes automated security testing:

- Unit tests for security-critical components
- Integration tests for authentication and authorization
- Penetration testing scripts
- Fuzzing for API endpoints

### Manual Security Testing

Regular manual security testing:

- Code reviews with a security focus
- Penetration testing
- Security audits
- Threat modeling

## Incident Response

### Security Incident Response Plan

The system includes a security incident response plan:

1. **Detection**: Monitoring systems to detect potential security incidents
2. **Containment**: Procedures to contain the impact of an incident
3. **Eradication**: Removing the cause of the incident
4. **Recovery**: Restoring systems to normal operation
5. **Lessons Learned**: Analyzing the incident to prevent recurrence

### Backup and Recovery

Data protection measures:

- Regular database backups
- Backup encryption
- Backup testing
- Disaster recovery procedures

## Compliance and Standards

The system is designed to comply with relevant security standards and best practices:

- OWASP Top 10
- NIST Cybersecurity Framework
- CWE/SANS Top 25 Most Dangerous Software Errors
- Industry-specific standards for gaming platforms

## Conclusion

Security is a fundamental aspect of the Live Ops Events System. This document outlines the comprehensive security measures implemented throughout the system, from authentication and authorization to secure deployment and incident response. By following these security specifications, the system provides a robust and secure platform for managing live events in games.

The security measures are designed to protect against common threats while maintaining performance and usability. Regular security reviews and updates will ensure that the system remains secure as new threats emerge and best practices evolve. 