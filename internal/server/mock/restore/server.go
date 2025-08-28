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

	var backup *model.Backup
	exists := false
	for _, environment := range s.Model.GetEnvironments() {
		backup, exists = environment.Backup[req.Backup]
		if exists {
			break
		}
	}
	if !exists {
		return nil, fmt.Errorf("backup does not exist")
	}
	if backup.Status() != pb.BackupStatus_Completed {
		return nil, fmt.Errorf("backup is not available for restore")
	}

	environment, err := s.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
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

	restoreResponse, err := buildRestore(s.Model, req.ID)
	if err != nil {
		return nil, err
	}

	resp := &pb.RestoreGetResponse{
		Restore: restoreResponse,
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
		summary, _ := buildRestore(s.Model, value.Id)
		resp.List = append(resp.List, summary)
	}

	sort.Slice(resp.List, func(i, j int) bool {
		return resp.List[i].StartTime > resp.List[j].StartTime
	})

	return resp, nil
}

func buildRestore(model *model.Model, id string) (*pb.RestoreStatus, error) {
	restore, err := model.GetRestore(id)
	if err != nil {
		return nil, err
	}

	restoreResponse := &pb.RestoreStatus{
		Name:      restore.Id,
		Backup:    restore.BackupId,
		Phase:     restore.Status(),
		StartTime: restore.StartTime.Format(time.RFC3339),
		Duration:  restore.Duration.String(),
		Databases: []string{"default"},
		Volumes:   []string{"public", "private"},
	}
	if restore.Status() != pb.RestoreStatus_InProgress {
		restoreResponse.CompletionTime = restore.StartTime.Add(restore.Duration).Format(time.RFC3339)
	}

	return restoreResponse, nil
}
