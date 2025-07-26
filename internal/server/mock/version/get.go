package version

import (
	"context"

	"github.com/skpr/api/pb"
)

// Get the api version from the server.
func (s *VersionServer) Get(ctx context.Context, req *pb.VersionGetRequest) (*pb.VersionGetResponse, error) {
	resp := &pb.VersionGetResponse{
		Version:   "0.34.5",
		BuildDate: "2025-05-06",
	}

	return resp, nil
}
