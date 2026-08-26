package trace

import (
	"sync"

	"github.com/skpr/api/pb"
)

// Server implements the GRPC "trace" definition.
type Server struct {
	pb.UnimplementedTraceServer

	lock      sync.Mutex
	suspended map[string]bool
}

// SetSuspended marks tracing as suspended, or resumed, for an environment.
func (s *Server) SetSuspended(environment string, suspended bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.suspended == nil {
		s.suspended = make(map[string]bool)
	}

	s.suspended[environment] = suspended
}

// IsSuspended determines if tracing has been suspended for an environment.
func (s *Server) IsSuspended(environment string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.suspended[environment]
}
