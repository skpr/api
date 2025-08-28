package backup

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
	pb.UnimplementedBackupServer
}

func (s *Server) Create(ctx context.Context, req *pb.BackupCreateRequest) (*pb.BackupCreateResponse, error) {
	environment, err := s.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

	backup := model.NewBackup(environment.Environment.Name)
	environment.AddBackup(backup)
	resp := &pb.BackupCreateResponse{
		ID: backup.Id,
	}
	return resp, nil
}

func (s *Server) Get(ctx context.Context, req *pb.BackupGetRequest) (*pb.BackupGetResponse, error) {
	if req.ID == "" {
		return nil, fmt.Errorf("id not provided")
	}

	backupResponse, err := buildBackup(s.Model, req.ID)
	if err != nil {
		return nil, err
	}

	resp := &pb.BackupGetResponse{
		Backup: backupResponse,
	}
	return resp, nil
}

func (s *Server) List(ctx context.Context, req *pb.BackupListRequest) (*pb.BackupListResponse, error) {
	environment, err := s.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

	resp := &pb.BackupListResponse{}

	for _, value := range environment.Backup {
		summary, _ := buildBackup(s.Model, value.Id)
		resp.List = append(resp.List, summary)
	}

	sort.Slice(resp.List, func(i, j int) bool {
		return resp.List[i].StartTime > resp.List[j].StartTime
	})

	return resp, nil
}

func buildBackup(model *model.Model, id string) (*pb.BackupStatus, error) {
	backup, err := model.GetBackup(id)
	if err != nil {
		return nil, err
	}

	backupResponse := &pb.BackupStatus{
		Name:      backup.Id,
		Phase:     backup.Status(),
		StartTime: backup.StartTime.Format(time.RFC3339),
		Duration:  backup.Duration.String(),
		Databases: []string{"default"},
		Volumes:   []string{"public", "private"},
	}
	if backup.Status() != pb.BackupStatus_InProgress {
		backupResponse.CompletionTime = backup.StartTime.Add(backup.Duration).Format(time.RFC3339)
	}

	return backupResponse, nil
}
