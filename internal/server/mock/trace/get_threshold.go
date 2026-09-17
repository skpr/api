package trace

import (
	"context"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/skpr/api/pb"
)

// GetThreshold returns the minimum call duration which is traced for a specific environment.
func (s *Server) GetThreshold(_ context.Context, req *pb.TraceGetThresholdRequest) (*pb.TraceGetThresholdResponse, error) {
	return &pb.TraceGetThresholdResponse{
		Threshold: durationpb.New(s.getThreshold(req.Environment)),
	}, nil
}
