package purge

import (
	"context"
	"sort"
	"time"

	"github.com/skpr/api/internal/model"
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "cron" definition.
type Server struct {
	Model *model.Model
	pb.UnimplementedPurgeServer
}

func (s *Server) List(ctx context.Context, req *pb.PurgeListRequest) (*pb.PurgeListResponse, error) {
	environment, err := s.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

	resp := &pb.PurgeListResponse{}

	for _, value := range environment.Purge {
		status := "completed"
		if time.Now().Add(-2 * time.Minute).Before(value.Created) {
			status = "in-progress"
		}
		summary := &pb.RequestSummary{
			ID:      value.Id,
			Created: value.Created.Format(time.RFC3339),
			Paths:   value.Paths,
			Status:  status,
		}
		resp.Requests = append(resp.Requests, summary)
	}

	sort.Slice(resp.Requests, func(i, j int) bool {
		return resp.Requests[i].Created > resp.Requests[j].Created
	})

	return resp, nil
}

func (s *Server) Create(ctx context.Context, req *pb.PurgeCreateRequest) (*pb.PurgeCreateResponse, error) {
	environment, err := s.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

	purge := model.NewPurge(req.Paths)
	environment.AddPurge(purge)
	resp := &pb.PurgeCreateResponse{
		ID: purge.Id,
	}
	return resp, nil
}
