package trace

import (
	"context"

	"github.com/skpr/api/pb"
)

// Resume tracing for a specific environment.
func (s *Server) Resume(_ context.Context, req *pb.TraceResumeRequest) (*pb.TraceResumeResponse, error) {
	s.SetSuspended(req.Environment, false)

	return &pb.TraceResumeResponse{}, nil
}
