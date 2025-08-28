package project

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/skpr/api/internal/model"
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "environments" definition.
type Server struct {
	Model *model.Model
	pb.UnimplementedProjectServer
}

func (s *Server) List(ctx context.Context, req *pb.ProjectListRequest) (*pb.ProjectListResponse, error) {
	resp := &pb.ProjectListResponse{}

	for _, project := range s.Model.GetProjects() {
		respProject, err := buildProject(s.Model, project.Id)
		if err == nil {
			resp.Projects = append(resp.Projects, respProject)
		}
	}

	return resp, nil
}

func (s *Server) Get(ctx context.Context, req *pb.ProjectGetRequest) (*pb.ProjectGetResponse, error) {
	projectId, err := getMetadata(ctx)
	if err != nil {
		return nil, err
	}

	project, err := buildProject(s.Model, projectId)
	if err != nil {
		return nil, err
	}

	return &pb.ProjectGetResponse{
		Project: project,
	}, nil
}

func (s *Server) SetTags(ctx context.Context, req *pb.SetTagsRequest) (*pb.SetTagsResponse, error) {
	projectId, err := getMetadata(ctx)
	if err != nil {
		return nil, err
	}

	project, err := s.Model.GetProject(projectId)
	if err != nil {
		return nil, err
	}

	project.Tags = req.Tags

	return &pb.SetTagsResponse{}, nil
}

func (s *Server) SetContact(ctx context.Context, req *pb.SetContactRequest) (*pb.SetContactResponse, error) {
	projectId, err := getMetadata(ctx)
	if err != nil {
		return nil, err
	}

	project, err := s.Model.GetProject(projectId)
	if err != nil {
		return nil, err
	}

	project.Contact = req.Contact

	return &pb.SetContactResponse{}, nil
}

func getMetadata(ctx context.Context) (string, error) {
	data, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("gRPC requset contained no metadata")
	}

	// Access metadata values
	projectMetadata := data.Get("project")
	if len(projectMetadata) == 0 {
		return "", fmt.Errorf("project metadata not set")
	}

	return projectMetadata[0], nil
}

func buildProject(model *model.Model, id string) (*pb.Project, error) {
	project, err := model.GetProject(id)
	if err != nil {
		return nil, err
	}

	respProject := &pb.Project{
		ID:      project.Id,
		Name:    project.Name,
		Tags:    project.Tags,
		Contact: project.Contact,
		Size:    project.Size,
		ResourceTotals: &pb.ProjectResourceTotals{
			CPU:      0,
			Memory:   0,
			Replicas: 0,
		},
		Environments: &pb.ProjectEnvironments{
			NonProd: []string{},
		},
	}

	environments := model.GetEnvironments()
	for _, env := range environments {
		respProject.ResourceTotals.CPU = respProject.ResourceTotals.CPU + env.Environment.Resources.CPU.Limit
		respProject.ResourceTotals.Memory = respProject.ResourceTotals.Memory + env.Environment.Resources.Memory.Limit
		respProject.ResourceTotals.Replicas = respProject.ResourceTotals.Replicas + env.Environment.Resources.Replicas.Max

		if env.Environment.Production {
			respProject.Environments.Prod = env.Environment.Name
			respProject.Version = env.Environment.Version
		} else {
			respProject.Environments.NonProd = append(respProject.Environments.NonProd, env.Environment.Name)
		}
	}

	return respProject, nil
}
