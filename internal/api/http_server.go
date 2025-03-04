package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tombombadilom/liveops/internal/auth"
	"github.com/tombombadilom/liveops/internal/models"
	"github.com/tombombadilom/liveops/internal/service"
)

// HTTPServer handles HTTP API requests
type HTTPServer struct {
	router       *gin.Engine
	eventService *service.EventService
	authService  *auth.AuthService
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(eventService *service.EventService, authService *auth.AuthService) *HTTPServer {
	// Create router
	router := gin.New()

	// Use middleware
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware())

	server := &HTTPServer{
		router:       router,
		eventService: eventService,
		authService:  authService,
	}

	// Register routes
	server.registerRoutes()

	return server
}

// Handler returns the HTTP handler
func (s *HTTPServer) Handler() http.Handler {
	return s.router
}

// registerRoutes sets up all API routes
func (s *HTTPServer) registerRoutes() {
	// Public routes
	s.router.GET("/health", s.healthCheck)

	// API routes (require authentication)
	api := s.router.Group("/api")
	api.Use(s.authMiddleware())
	{
		// Events
		events := api.Group("/events")
		{
			events.GET("", s.listEvents)
			events.GET("/active", s.listActiveEvents)
			events.GET("/:id", s.getEvent)
			events.POST("", s.createEvent)
			events.PUT("/:id", s.updateEvent)
			events.DELETE("/:id", s.deleteEvent)
		}

		// Admin routes (require admin role)
		admin := api.Group("/admin")
		admin.Use(s.adminMiddleware())
		{
			// Users
			admin.GET("/users", s.listUsers)
			admin.POST("/users", s.createUser)
			admin.GET("/users/:id", s.getUser)

			// API Keys
			admin.GET("/users/:id/keys", s.listAPIKeys)
			admin.POST("/users/:id/keys", s.createAPIKey)
			admin.DELETE("/keys/:id", s.revokeAPIKey)
		}
	}
}

// loggerMiddleware logs HTTP requests
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		status := c.Writer.Status()

		log.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Msg("HTTP request")
	}
}

// authMiddleware authenticates API requests
func (s *HTTPServer) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
			})
			return
		}

		// Authenticate API key
		user, err := s.authService.AuthenticateAPIKey(apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			return
		}

		// Store user in context
		c.Set("user", user)
		c.Next()
	}
}

// adminMiddleware ensures the user has admin role
func (s *HTTPServer) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			return
		}

		// Check if user has admin role
		if err := s.authService.CheckPermission(user.(*models.User), "admin"); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			return
		}

		c.Next()
	}
}

// healthCheck handles health check requests
func (s *HTTPServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// listEvents handles GET /api/events
func (s *HTTPServer) listEvents(c *gin.Context) {
	events, err := s.eventService.ListEvents(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// listActiveEvents handles GET /api/events/active
func (s *HTTPServer) listActiveEvents(c *gin.Context) {
	events, err := s.eventService.ListEvents(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// getEvent handles GET /api/events/:id
func (s *HTTPServer) getEvent(c *gin.Context) {
	id := c.Param("id")

	event, err := s.eventService.GetEvent(id)
	if err != nil {
		if err == models.ErrEventNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, event)
}

// createEvent handles POST /api/events
func (s *HTTPServer) createEvent(c *gin.Context) {
	// Get user from context
	user := c.MustGet("user").(*models.User)

	// Check permission
	if err := s.authService.CheckPermission(user, "create"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Parse request
	var req struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time" binding:"required"`
		EndTime     time.Time `json:"end_time" binding:"required"`
		Rewards     string    `json:"rewards"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create event
	event, err := s.eventService.CreateEvent(req.Title, req.Description, req.StartTime, req.EndTime, req.Rewards)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// updateEvent handles PUT /api/events/:id
func (s *HTTPServer) updateEvent(c *gin.Context) {
	// Get user from context
	user := c.MustGet("user").(*models.User)

	// Check permission
	if err := s.authService.CheckPermission(user, "update"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Get event ID
	id := c.Param("id")

	// Parse request
	var req struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time" binding:"required"`
		EndTime     time.Time `json:"end_time" binding:"required"`
		Rewards     string    `json:"rewards"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update event
	event, err := s.eventService.UpdateEvent(id, req.Title, req.Description, req.StartTime, req.EndTime, req.Rewards)
	if err != nil {
		if err == models.ErrEventNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, event)
}

// deleteEvent handles DELETE /api/events/:id
func (s *HTTPServer) deleteEvent(c *gin.Context) {
	// Get user from context
	user := c.MustGet("user").(*models.User)

	// Check permission
	if err := s.authService.CheckPermission(user, "delete"); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Get event ID
	id := c.Param("id")

	// Delete event
	err := s.eventService.DeleteEvent(id)
	if err != nil {
		if err == models.ErrEventNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// listUsers handles GET /api/admin/users
func (s *HTTPServer) listUsers(c *gin.Context) {
	users, err := s.authService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// createUser handles POST /api/admin/users
func (s *HTTPServer) createUser(c *gin.Context) {
	// Parse request
	var req struct {
		Username string      `json:"username" binding:"required"`
		Role     models.Role `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user, err := s.authService.CreateUser(req.Username, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// getUser handles GET /api/admin/users/:id
func (s *HTTPServer) getUser(c *gin.Context) {
	id := c.Param("id")

	user, err := s.authService.GetUser(id)
	if err != nil {
		if err == models.ErrInvalidID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// listAPIKeys handles GET /api/admin/users/:id/keys
func (s *HTTPServer) listAPIKeys(c *gin.Context) {
	id := c.Param("id")

	keys, err := s.authService.ListAPIKeys(id)
	if err != nil {
		if err == models.ErrInvalidID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, keys)
}

// createAPIKey handles POST /api/admin/users/:id/keys
func (s *HTTPServer) createAPIKey(c *gin.Context) {
	id := c.Param("id")

	// Parse request
	var req struct {
		ValidDays int `json:"valid_days" binding:"required,min=1,max=365"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create API key
	key, err := s.authService.CreateAPIKey(id, req.ValidDays)
	if err != nil {
		if err == models.ErrInvalidID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, key)
}

// revokeAPIKey handles DELETE /api/admin/keys/:id
func (s *HTTPServer) revokeAPIKey(c *gin.Context) {
	id := c.Param("id")

	err := s.authService.RevokeAPIKey(id)
	if err != nil {
		if err == models.ErrInvalidID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
