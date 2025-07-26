package compass

import (
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "compass" definition.
type Server struct {
	pb.UnimplementedCompassServer
}
