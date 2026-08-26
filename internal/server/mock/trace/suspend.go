package trace

import (
	"context"

	"github.com/skpr/api/pb"
)

// Suspend tracing for a specific environment.
func (s *Server) Suspend(_ context.Context, req *pb.TraceSuspendRequest) (*pb.TraceSuspendResponse, error) {
	s.SetSuspended(req.Environment, true)

	return &pb.TraceSuspendResponse{}, nil
}
