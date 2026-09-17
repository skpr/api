package trace

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/skpr/api/pb"
)

// SetThreshold sets the minimum call duration which is traced for a specific environment.
func (s *Server) SetThreshold(_ context.Context, req *pb.TraceSetThresholdRequest) (*pb.TraceSetThresholdResponse, error) {
	// An omitted threshold is not the same as a zero one: zero traces every call,
	// which is expensive enough that a caller has to ask for it rather than get
	// it by leaving the field out.
	if req.Threshold == nil {
		return nil, status.Error(codes.InvalidArgument, "threshold is required")
	}

	threshold := req.Threshold.AsDuration()

	// A negative threshold would trace nothing at all, which is not what a caller
	// asking for less tracing means. Suspend is how tracing gets turned off.
	if threshold < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "threshold cannot be negative: %s", threshold)
	}

	s.setThreshold(req.Environment, threshold)

	return &pb.TraceSetThresholdResponse{}, nil
}
