package trace

import (
	"context"

	"github.com/skpr/api/pb"
)

// GetSuspended returns whether tracing is currently suspended for a specific environment.
func (s *Server) GetSuspended(_ context.Context, req *pb.TraceGetSuspendedRequest) (*pb.TraceGetSuspendedResponse, error) {
	return &pb.TraceGetSuspendedResponse{
		Suspended: s.IsSuspended(req.Environment),
	}, nil
}
