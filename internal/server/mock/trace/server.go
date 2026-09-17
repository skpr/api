package trace

import (
	"sync"
	"time"

	"github.com/skpr/api/pb"
)

// DefaultThreshold is the tracing threshold an environment has until one is set
// for it, matching the default which Compass ships with.
const DefaultThreshold = time.Millisecond

// Server implements the GRPC "trace" definition.
type Server struct {
	pb.UnimplementedTraceServer

	lock       sync.Mutex
	suspended  map[string]bool
	thresholds map[string]time.Duration
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

// setThreshold records the minimum call duration which is traced for an environment.
func (s *Server) setThreshold(environment string, threshold time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.thresholds == nil {
		s.thresholds = make(map[string]time.Duration)
	}

	s.thresholds[environment] = threshold
}

// getThreshold returns the minimum call duration which is traced for an
// environment, falling back to the default for an environment which has not had
// one set.
func (s *Server) getThreshold(environment string) time.Duration {
	s.lock.Lock()
	defer s.lock.Unlock()

	threshold, ok := s.thresholds[environment]
	if !ok {
		return DefaultThreshold
	}

	return threshold
}
