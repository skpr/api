package restore

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/skpr/api/internal/model"
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "backup" definition.
type Server struct {
	Model *model.Model
	pb.UnimplementedRestoreServer
}

func (s *Server) Create(ctx context.Context, req *pb.RestoreCreateRequest) (*pb.RestoreCreateResponse, error) {
	if req.Backup == "" {
		return nil, fmt.Errorf("backup not provided")
	}

	environment, err := s.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

	backup, exists := environment.Backup[req.Backup]
	if !exists {
		return nil, fmt.Errorf("backup does not exist")
	}
	if backup.Status() == pb.BackupStatus_InProgress {
		return nil, fmt.Errorf("backup is not yet complete")
	}

	restore := model.NewRestore(environment.Environment.Name, req.Backup)
	environment.AddRestore(restore)
	resp := &pb.RestoreCreateResponse{
		ID: restore.Id,
	}
	return resp, nil
}

func (s *Server) Get(ctx context.Context, req *pb.RestoreGetRequest) (*pb.RestoreGetResponse, error) {
	if req.ID == "" {
		return nil, fmt.Errorf("id not provided")
	}

	restore, err := s.Model.GetRestore(req.ID)
	if err != nil {
		return nil, err
	}

	resp := &pb.RestoreGetResponse{
		Restore: &pb.RestoreStatus{
			Name:      restore.Id,
			Backup:    restore.BackupId,
			Phase:     restore.Status(),
			StartTime: restore.StartTime.Format(time.RFC3339),
			Duration:  restore.Duration.String(),
			Databases: []string{"default"},
			Volumes:   []string{"public", "private"},
		},
	}
	return resp, nil
}

func (s *Server) List(ctx context.Context, req *pb.RestoreListRequest) (*pb.RestoreListResponse, error) {
	environment, err := s.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

	resp := &pb.RestoreListResponse{}

	for _, value := range environment.Restore {
		summary := &pb.RestoreStatus{
			Name:      value.Id,
			Backup:    value.BackupId,
			Phase:     value.Status(),
			StartTime: value.StartTime.Format(time.RFC3339),
			Duration:  value.Duration.String(),
			Databases: []string{"default"},
			Volumes:   []string{"public", "private"},
		}
		resp.List = append(resp.List, summary)
	}

	sort.Slice(resp.List, func(i, j int) bool {
		return resp.List[i].StartTime > resp.List[j].StartTime
	})

	return resp, nil
}
