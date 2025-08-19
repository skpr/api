package project

import (
	"context"

	"github.com/skpr/api/pb"
)

// Projects that we are using for mock data.
// todo Move to modelling.
var Projects = []*pb.Project{
	{
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
	{
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
	{
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
	return &pb.ProjectListResponse{
		Projects: Projects,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *pb.ProjectGetRequest) (*pb.ProjectGetResponse, error) {
	resp := &pb.ProjectGetResponse{
		Project: Projects[0],
	}
	return resp, nil
}
