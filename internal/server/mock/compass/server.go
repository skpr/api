package compass

import (
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "mysql" definition.
type Server struct {
	pb.UnimplementedCompassServer
}
