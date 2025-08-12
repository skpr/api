package config

import (
	"context"

	"github.com/skpr/api/internal/model"
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "backup" definition.
type Server struct {
	Model *model.Model
	pb.UnimplementedConfigServer
}

func (s *Server) List(ctx context.Context, req *pb.ConfigListRequest) (*pb.ConfigListResponse, error) {
	environment, err := s.Model.GetEnvironment(req.Name)
	if err != nil {
		return nil, err
	}

	// TODO Hide secrets if not asked for

	resp := &pb.ConfigListResponse{}
	for _, value := range environment.Config {
		resp.List = append(resp.List, value)
	}
	return resp, nil
}

func (s *Server) Get(ctx context.Context, req *pb.ConfigGetRequest) (*pb.ConfigGetResponse, error) {
	environment, err := s.Model.GetEnvironment(req.Name)
	if err != nil {
		return nil, err
	}

	value, err := environment.GetConfig(req.Key)
	if err != nil {
		return nil, err
	}

	// TODO Hide secrets if not asked for

	return &pb.ConfigGetResponse{
		Config: value,
	}, nil
}

func (s *Server) Set(ctx context.Context, req *pb.ConfigSetRequest) (*pb.ConfigSetResponse, error) {
	environment, err := s.Model.GetEnvironment(req.Name)
	if err != nil {
		return nil, err
	}

	_, err = environment.GetConfig(req.Config.Key)
	environment.AddConfig(req.Config)

	return &pb.ConfigSetResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.ConfigDeleteRequest) (*pb.ConfigDeleteResponse, error) {
	environment, err := s.Model.GetEnvironment(req.Name)
	if err != nil {
		return nil, err
	}

	err = environment.DeleteConfig(req.Key)
	if err != nil {
		return nil, err
	}

	return &pb.ConfigDeleteResponse{}, nil
}
