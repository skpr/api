package trace

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/skpr/api/pb"
)

// --- SetThreshold ---

func TestSetThreshold(t *testing.T) {
	cases := []struct {
		name      string
		threshold *durationpb.Duration
		wantCode  codes.Code
		want      time.Duration
	}{
		{"sets the threshold", durationpb.New(5 * time.Millisecond), codes.OK, 5 * time.Millisecond},
		{"zero traces every call", durationpb.New(0), codes.OK, 0},
		{"omitted is rejected", nil, codes.InvalidArgument, DefaultThreshold},
		{"negative is rejected", durationpb.New(-time.Millisecond), codes.InvalidArgument, DefaultThreshold},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := &Server{}

			_, err := srv.SetThreshold(context.TODO(), &pb.TraceSetThresholdRequest{
				Environment: "dev",
				Threshold:   c.threshold,
			})

			if got := status.Code(err); got != c.wantCode {
				t.Errorf("got code %s, want %s", got, c.wantCode)
			}

			if got := getThreshold(t, srv, "dev"); got != c.want {
				t.Errorf("got threshold %s, want %s", got, c.want)
			}
		})
	}
}

// An environment which has never had a threshold set still has to stream, so it
// falls back to the default rather than to zero, which would trace every call.
func TestGetThresholdDefaults(t *testing.T) {
	srv := &Server{}

	if got := getThreshold(t, srv, "dev"); got != DefaultThreshold {
		t.Errorf("got threshold %s, want %s", got, DefaultThreshold)
	}
}

// Thresholds are held per environment, so lowering one does not quietly start
// tracing every call on another.
func TestSetThresholdIsPerEnvironment(t *testing.T) {
	srv := &Server{}

	if _, err := srv.SetThreshold(context.TODO(), &pb.TraceSetThresholdRequest{
		Environment: "dev",
		Threshold:   durationpb.New(0),
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := getThreshold(t, srv, "prod"); got != DefaultThreshold {
		t.Errorf("got threshold %s for prod, want %s", got, DefaultThreshold)
	}
}

// --- GetSuspended ---

// Suspending and resuming are only useful if a client can find out which of the
// two an environment is in, without opening a stream to see if it fails.
func TestGetSuspended(t *testing.T) {
	srv := &Server{}

	if getSuspended(t, srv, "dev") {
		t.Error("got suspended for an environment which was never suspended")
	}

	if _, err := srv.Suspend(context.TODO(), &pb.TraceSuspendRequest{Environment: "dev"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !getSuspended(t, srv, "dev") {
		t.Error("got resumed after Suspend")
	}

	// Suspending one environment does not suspend the rest.
	if getSuspended(t, srv, "prod") {
		t.Error("got suspended for prod after suspending dev")
	}

	if _, err := srv.Resume(context.TODO(), &pb.TraceResumeRequest{Environment: "dev"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if getSuspended(t, srv, "dev") {
		t.Error("got suspended after Resume")
	}
}

// --- aboveThreshold ---

func TestAboveThreshold(t *testing.T) {
	spans := []*pb.TraceSpan{
		{Name: "slow", Elapsed: durationpb.New(6 * time.Millisecond), Calls: 3},
		{Name: "exactly", Elapsed: durationpb.New(time.Millisecond), Calls: 2},
		{Name: "fast", Elapsed: durationpb.New(250 * time.Microsecond), Calls: 5},
	}

	cases := []struct {
		name      string
		threshold time.Duration
		want      []string
	}{
		{"everything is traced at zero", 0, []string{"slow", "exactly", "fast"}},
		{"the threshold itself is traced", time.Millisecond, []string{"slow", "exactly"}},
		{"a coarse threshold keeps the slowest", 5 * time.Millisecond, []string{"slow"}},
		{"nothing reaches an impossible threshold", time.Hour, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			retained, calls := aboveThreshold(spans, c.threshold)

			if len(retained) != len(c.want) {
				t.Fatalf("got %d spans, want %d", len(retained), len(c.want))
			}

			for i, name := range c.want {
				if retained[i].Name != name {
					t.Errorf("got span %q at %d, want %q", retained[i].Name, i, name)
				}
			}

			// The request made these calls whether or not a span survived the
			// threshold, so the count is the same at every threshold.
			if calls != 10 {
				t.Errorf("got %d calls, want 10", calls)
			}
		})
	}
}

// --- helpers ---

func getThreshold(t *testing.T, srv *Server, environment string) time.Duration {
	t.Helper()

	resp, err := srv.GetThreshold(context.TODO(), &pb.TraceGetThresholdRequest{Environment: environment})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return resp.Threshold.AsDuration()
}

func getSuspended(t *testing.T, srv *Server, environment string) bool {
	t.Helper()

	resp, err := srv.GetSuspended(context.TODO(), &pb.TraceGetSuspendedRequest{Environment: environment})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return resp.Suspended
}
