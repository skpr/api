package environment

import (
	"context"

	"github.com/skpr/api/internal/model"
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "environments" definition.
type Server struct {
	Model *model.Model
	pb.UnimplementedEnvironmentServer
}

func (c *Server) Get(ctx context.Context, req *pb.EnvironmentGetRequest) (*pb.EnvironmentGetResponse, error) {
	environment, err := c.Model.GetEnvironment(req.Name)
	if err != nil {
		return nil, err
	}

	resp := &pb.EnvironmentGetResponse{
		Environment: environment.Environment,
	}
	return resp, nil
}

func (c *Server) List(ctx context.Context, req *pb.EnvironmentListRequest) (*pb.EnvironmentListResponse, error) {
	resp := &pb.EnvironmentListResponse{}

	environments := c.Model.GetEnvironments()
	for _, value := range environments {
		resp.Environments = append(resp.Environments, value.Environment)
	}

	return resp, nil
}

func (c *Server) Delete(ctx context.Context, req *pb.EnvironmentDeleteRequest) (*pb.EnvironmentDeleteResponse, error) {
	_, err := c.Model.GetEnvironment(req.Name)
	if err != nil {
		return nil, err
	}

	c.Model.DeleteEnvironment(req.Name)

	resp := &pb.EnvironmentDeleteResponse{}
	return resp, nil
}

func (c *Server) Suspend(ctx context.Context, req *pb.EnvironmentSuspendRequest) (*pb.EnvironmentSuspendResponse, error) {
	environment, err := c.Model.GetEnvironment(req.Name)
	if err != nil {
		return nil, err
	}

	environment.Suspended = true

	resp := &pb.EnvironmentSuspendResponse{}
	return resp, nil
}

func (c *Server) Resume(ctx context.Context, req *pb.EnvironmentResumeRequest) (*pb.EnvironmentResumeResponse, error) {
	environment, err := c.Model.GetEnvironment(req.Name)
	if err != nil {
		return nil, err
	}

	environment.Suspended = false

	resp := &pb.EnvironmentResumeResponse{}
	return resp, nil
}
