package project

import (
	"context"
	"fmt"
	
	"google.golang.org/grpc/metadata"

	"github.com/skpr/api/pb"
)

// Projects that we are using for mock data.
// todo Move to modelling.
var Projects = map[string]*pb.Project{
	"project1": {
		ID:   "project1",
		Name: "Project One",
		Tags: []string{
			"group1",
			"group3",
		},
		Contact: "admin@example.com",
		Version: "1.0.0",
		Size:    "small",
		ResourceTotals: &pb.ProjectResourceTotals{
			CPU:      10,
			Memory:   4096,
			Replicas: 3,
		},
		Environments: &pb.ProjectEnvironments{
			Prod: "prod",
			NonProd: []string{
				"dev",
				"stg",
			},
		},
	},
	"project2": {
		ID:   "project2",
		Name: "Project Two",
		Tags: []string{
			"group1",
		},
		Contact: "admin@example.com",
		Version: "2.1.0",
		Size:    "medium",
		ResourceTotals: &pb.ProjectResourceTotals{
			CPU:      100,
			Memory:   2048,
			Replicas: 2,
		},
		Environments: &pb.ProjectEnvironments{
			Prod: "prod",
			NonProd: []string{
				"dev",
				"stg",
			},
		},
	},
	"project3": {
		ID:   "project3",
		Name: "Project Three",
		Tags: []string{
			"group2",
			"group4",
		},
		Contact: "admin@example.com",
		Version: "1.2.3",
		Size:    "large",
		ResourceTotals: &pb.ProjectResourceTotals{
			CPU:      3000,
			Memory:   24576,
			Replicas: 5,
		},
		Environments: &pb.ProjectEnvironments{
			Prod: "prod",
			NonProd: []string{
				"dev",
				"stg",
			},
		},
	},
}

// Server implements the GRPC "environments" definition.
type Server struct {
	pb.UnimplementedProjectServer
}

func (s *Server) List(ctx context.Context, req *pb.ProjectListRequest) (*pb.ProjectListResponse, error) {
	resp := &pb.ProjectListResponse{}

	for _, project := range Projects {
		resp.Projects = append(resp.Projects, project)
	}

	return resp, nil
}

func (s *Server) Get(ctx context.Context, req *pb.ProjectGetRequest) (*pb.ProjectGetResponse, error) {
	data, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("gRPC requset contained no metadata")
	}

	// Access metadata values
	projectMetadata := data.Get("project")
	if len(projectMetadata) == 0 {
		return nil, fmt.Errorf("project metadata not set")
	}

	project, ok := Projects[projectMetadata[0]]
	if !ok {
		return nil, fmt.Errorf("project not found")
	}

	return &pb.ProjectGetResponse{
		Project: project,
	}, nil
}
