package cluster

import (
	"context"

	"github.com/skpr/api/pb"
)

// Server implements the GRPC "cluster" definition.
type Server struct {
	pb.UnimplementedClusterServer
}

// Get the cluster status from the server.
func (s *Server) Get(ctx context.Context, req *pb.ClusterGetRequest) (*pb.ClusterGetResponse, error) {
	resp := &pb.ClusterGetResponse{
		Status:    pb.ClusterGetResponse_Operational,
		Version:   "0.34.5-mock",
		BuildDate: "2025-05-06",
		Endpoints: &pb.ClusterEndpoints{
			Console: "https://console.mock.localhost",
			SSH:     "console.mock.localhost:5000",
		},
	}

	return resp, nil
}
