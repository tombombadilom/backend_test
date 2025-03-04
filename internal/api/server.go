package api

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/soheilhy/cmux"
	"github.com/tombombadilom/liveops/internal/auth"
	"github.com/tombombadilom/liveops/internal/service"
)

// Server represents the API server that handles both HTTP and gRPC
type Server struct {
	httpServer *HTTPServer
	grpcServer *GRPCServer
	listener   net.Listener
	port       int
}

// NewServer creates a new API server
func NewServer(port int, eventService *service.EventService, authService *auth.AuthService) *Server {
	return &Server{
		httpServer: NewHTTPServer(eventService, authService),
		grpcServer: NewGRPCServer(eventService, authService),
		port:       port,
	}
}

// Start starts the server
func (s *Server) Start() error {
	// Create listener
	addr := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener

	// Create connection multiplexer
	mux := cmux.New(listener)

	// Match gRPC connections
	grpcListener := mux.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
	)

	// Match HTTP connections
	httpListener := mux.Match(cmux.Any())

	// Start servers
	go s.serveGRPC(grpcListener)
	go s.serveHTTP(httpListener)

	// Start multiplexer
	log.Info().Int("port", s.port).Msg("Server started, listening on port")
	return mux.Serve()
}

// Stop stops the server
func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
}

// serveGRPC starts the gRPC server
func (s *Server) serveGRPC(listener net.Listener) {
	log.Info().Msg("Starting gRPC server")
	server := s.grpcServer.Server()
	if err := server.Serve(listener); err != nil {
		log.Error().Err(err).Msg("gRPC server error")
	}
}

// serveHTTP starts the HTTP server
func (s *Server) serveHTTP(listener net.Listener) {
	log.Info().Msg("Starting HTTP server")
	server := &http.Server{
		Handler: s.httpServer.Handler(),
	}
	if err := server.Serve(listener); err != nil {
		log.Error().Err(err).Msg("HTTP server error")
	}
}
