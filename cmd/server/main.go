package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tombombadilom/liveops/internal/api"
	"github.com/tombombadilom/liveops/internal/auth"
	"github.com/tombombadilom/liveops/internal/config"
	"github.com/tombombadilom/liveops/internal/db"
	"github.com/tombombadilom/liveops/internal/service"
)

func main() {
	// Load configuration
	cfg := config.New()
	cfg.ParseFlags()

	// Configure logging
	configureLogging(cfg.LogLevel)
	log.Info().Msg("Starting Live Ops Events System")

	// Create a context that is canceled on interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	setupSignalHandler(cancel)

	// Initialize database
	database, err := db.New(cfg.DBPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", cfg.DBPath).Msg("Failed to initialize database")
	}
	defer database.Close()
	log.Info().Str("path", cfg.DBPath).Msg("Database initialized")

	// Create repositories
	eventRepo := db.NewEventRepository(database)
	userRepo := db.NewUserRepository(database)
	apiKeyRepo := db.NewAPIKeyRepository(database)

	// Create services
	eventService := service.NewEventService(eventRepo)
	authService := auth.NewAuthService(userRepo, apiKeyRepo)

	// Create and start server
	server := api.NewServer(cfg.Port, eventService, authService)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	log.Info().Int("port", cfg.Port).Msg("Server started")

	// Wait for context cancellation (shutdown signal)
	<-ctx.Done()
	log.Info().Msg("Server shutting down")
	server.Stop()
}

// configureLogging sets up the logger with the specified log level
func configureLogging(level string) {
	// Set up pretty console logging
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	// Set log level
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// setupSignalHandler sets up a handler for interrupt signals
func setupSignalHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Info().Msg("Received interrupt signal")
		cancel()
	}()
}
