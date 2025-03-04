package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	port     = flag.Int("port", 8080, "The server port")
	dbPath   = flag.String("db_path", "./data/liveops.db", "Path to SQLite database file")
	logLevel = flag.String("log_level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	// Parse command line flags
	flag.Parse()

	// Configure logging
	configureLogging(*logLevel)
	log.Info().Msg("Starting Live Ops Events System")

	// Create a context that is canceled on interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	setupSignalHandler(cancel)

	// TODO: Initialize database
	// TODO: Set up HTTP and gRPC servers
	// TODO: Start the server

	log.Info().Msgf("Server listening on port %d", *port)

	// Keep the server running until context is canceled
	<-ctx.Done()
	log.Info().Msg("Server shutting down")
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