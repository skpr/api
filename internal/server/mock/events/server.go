package events

import (
	"context"
	"fmt"
	"time"

	"github.com/skpr/api/pb"
)

// Server implements the GRPC "events" definition.
type Server struct {
	pb.UnimplementedEventsServer
}

var mockEvents = []*pb.Event{
	{
		Date:    time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC).Round(time.Minute).Format(time.RFC3339),
		ID:      "ABCD1234",
		Type:    "ShellEvent",
		Source:  "skpr-project-1",
		Details: "Shell event triggered",
	},
	{
		Date:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC).Round(time.Minute).Format(time.RFC3339),
		ID:      "ABCD1234",
		Type:    "ShellEvent",
		Source:  "skpr-project-1",
		Details: "Shell event triggered",
	},
	{
		Date:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Round(time.Minute).Format(time.RFC3339),
		ID:      "ABCD1234",
		Type:    "ShellEvent",
		Source:  "skpr-project-1",
		Details: "Shell event triggered",
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
