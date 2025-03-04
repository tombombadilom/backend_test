package api

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tombombadilom/liveops/internal/auth"
	"github.com/tombombadilom/liveops/internal/models"
	"github.com/tombombadilom/liveops/internal/service"
	pb "github.com/tombombadilom/liveops/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCServer handles gRPC API requests
type GRPCServer struct {
	pb.UnimplementedEventServiceServer
	eventService *service.EventService
	authService  *auth.AuthService
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(eventService *service.EventService, authService *auth.AuthService) *GRPCServer {
	return &GRPCServer{
		eventService: eventService,
		authService:  authService,
	}
}

// Server returns a configured gRPC server
func (s *GRPCServer) Server() *grpc.Server {
	// Create gRPC server with interceptors
	server := grpc.NewServer(
		grpc.UnaryInterceptor(s.loggingInterceptor),
	)

	// Register services
	pb.RegisterEventServiceServer(server, s)

	return server
}

// loggingInterceptor logs gRPC requests
func (s *GRPCServer) loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	method := info.FullMethod

	// Process request
	resp, err := handler(ctx, req)

	// Log request
	latency := time.Since(start)
	status := "OK"
	if err != nil {
		status = err.Error()
	}

	log.Info().
		Str("method", method).
		Dur("latency", latency).
		Str("status", status).
		Msg("gRPC request")

	return resp, err
}

// authenticate checks the API key in the request
func (s *GRPCServer) authenticate(ctx context.Context) (*models.User, error) {
	// Get API key from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	apiKeys := md.Get("x-api-key")
	if len(apiKeys) == 0 {
		return nil, status.Error(codes.Unauthenticated, "API key required")
	}

	// Authenticate API key
	user, err := s.authService.AuthenticateAPIKey(apiKeys[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid API key")
	}

	return user, nil
}

// checkPermission checks if the user has permission for an action
func (s *GRPCServer) checkPermission(user *models.User, action string) error {
	if err := s.authService.CheckPermission(user, action); err != nil {
		return status.Error(codes.PermissionDenied, "permission denied")
	}
	return nil
}

// ListEvents implements the gRPC ListEvents method
func (s *GRPCServer) ListEvents(ctx context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	// Authenticate request
	_, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// Get events from service
	events, err := s.eventService.ListEvents(req.ActiveOnly)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to protobuf response
	pbEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		pbEvents[i] = &pb.Event{
			Id:          event.ID.String(),
			Title:       event.Title,
			Description: event.Description,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     timestamppb.New(event.EndTime),
			Rewards:     event.Rewards,
		}
	}

	return &pb.ListEventsResponse{
		Events: pbEvents,
	}, nil
}

// GetEvent implements the gRPC GetEvent method
func (s *GRPCServer) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.Event, error) {
	// Authenticate request
	_, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// Get event from service
	event, err := s.eventService.GetEvent(req.Id)
	if err != nil {
		if err == models.ErrEventNotFound {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to protobuf response
	return &pb.Event{
		Id:          event.ID.String(),
		Title:       event.Title,
		Description: event.Description,
		StartTime:   timestamppb.New(event.StartTime),
		EndTime:     timestamppb.New(event.EndTime),
		Rewards:     event.Rewards,
	}, nil
}

// CreateEvent implements the gRPC CreateEvent method
func (s *GRPCServer) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.Event, error) {
	// Authenticate request
	user, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// Check permission
	if err := s.checkPermission(user, "create"); err != nil {
		return nil, err
	}

	// Create event
	event, err := s.eventService.CreateEvent(
		req.Title,
		req.Description,
		req.StartTime.AsTime(),
		req.EndTime.AsTime(),
		req.Rewards,
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert to protobuf response
	return &pb.Event{
		Id:          event.ID.String(),
		Title:       event.Title,
		Description: event.Description,
		StartTime:   timestamppb.New(event.StartTime),
		EndTime:     timestamppb.New(event.EndTime),
		Rewards:     event.Rewards,
	}, nil
}

// UpdateEvent implements the gRPC UpdateEvent method
func (s *GRPCServer) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.Event, error) {
	// Authenticate request
	user, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// Check permission
	if err := s.checkPermission(user, "update"); err != nil {
		return nil, err
	}

	// Update event
	event, err := s.eventService.UpdateEvent(
		req.Id,
		req.Title,
		req.Description,
		req.StartTime.AsTime(),
		req.EndTime.AsTime(),
		req.Rewards,
	)
	if err != nil {
		if err == models.ErrEventNotFound {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert to protobuf response
	return &pb.Event{
		Id:          event.ID.String(),
		Title:       event.Title,
		Description: event.Description,
		StartTime:   timestamppb.New(event.StartTime),
		EndTime:     timestamppb.New(event.EndTime),
		Rewards:     event.Rewards,
	}, nil
}

// DeleteEvent implements the gRPC DeleteEvent method
func (s *GRPCServer) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	// Authenticate request
	user, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// Check permission
	if err := s.checkPermission(user, "delete"); err != nil {
		return nil, err
	}

	// Delete event
	err = s.eventService.DeleteEvent(req.Id)
	if err != nil {
		if err == models.ErrEventNotFound {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
