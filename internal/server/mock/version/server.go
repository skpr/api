package version

import (
	"context"

	"github.com/skpr/api/pb"
)

// Server implements the GRPC "version" definition.
type Server struct {
	pb.UnimplementedVersionServer
}

// Get the api version from the server.
func (s *Server) Get(ctx context.Context, req *pb.VersionGetRequest) (*pb.VersionGetResponse, error) {
	resp := &pb.VersionGetResponse{
		Version:   "0.34.5-mock",
		BuildDate: "2025-05-06",
	}

	return resp, nil
}
