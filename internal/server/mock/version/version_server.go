package version

import (
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "version" definition.
type VersionServer struct {
	pb.UnimplementedVersionServer
}
