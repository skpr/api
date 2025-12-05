package events

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/skpr/api/pb"
)

// Server implements the GRPC "events" definition.
type Server struct {
	pb.UnimplementedEventsServer
}

var mockEvents = []*pb.Event{
	{
		Timestamp: timestamppb.New(time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC).Round(time.Minute)),
		ID:        "ABCDXXXX",
		Severity:  pb.Event_Info,
		Type:      "ConfigSet",
		Message:   "A config was set: api.key",
	},
	{
		Timestamp: timestamppb.New(time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC).Round(time.Minute)),
		ID:        "ABCDYYYY",
		Severity:  pb.Event_Warning,
		Type:      "ErrorRate",
		Message:   "Elevated error rate has been detected",
	},
	{
		Timestamp: timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Round(time.Minute)),
		ID:        "ABCDZZZZ",
		Severity:  pb.Event_Error,
		Type:      "BackupFailed",
		Message:   "The following backup failed with the ID: xxxxxxxxxxxxxxxxx",
	},
	{
		Timestamp: timestamppb.New(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Round(time.Minute)),
		ID:        "ABCDWWWW",
		Severity:  pb.Event_Critical,
		Type:      "EnvironmentDown",
		Message:   "The environment is down due to a critical failure",
	},
}

// Get the list of events from the server.
func (s *Server) List(ctx context.Context, req *pb.EventsListRequest) (*pb.EventsListResponse, error) {
	if req.Environment == "" {
		return nil, fmt.Errorf("environment not provided")
	}

	resp := &pb.EventsListResponse{
		Events: mockEvents,
	}

	return resp, nil
}
