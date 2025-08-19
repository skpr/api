package mysql

import (
	"context"
	"time"

	"github.com/skpr/api/internal/model"
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "backup" definition.
type Server struct {
	Model *model.Model
	pb.UnimplementedMysqlServer
}

func (s *Server) ImageList(ctx context.Context, req *pb.ImageListRequest) (*pb.ImageListResponse, error) {
	environment, err := s.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

	resp := &pb.ImageListResponse{}
	for _, value := range environment.Mysql {
		image := &pb.ImageStatus{
			ID:        value.Id,
			Phase:     value.Status(),
			StartTime: value.StartTime.Format(time.RFC3339),
			Duration:  value.Duration.String(),
		}
		if value.Status() == pb.ImageStatus_Completed {
			image.CompletionTime = value.StartTime.Add(value.Duration).Format(time.RFC3339)
		}
		resp.List = append(resp.List, image)
	}
	return resp, nil
}
