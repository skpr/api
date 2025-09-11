package trace

import (
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "trace" definition.
type Server struct {
	pb.UnimplementedTraceServer
}
