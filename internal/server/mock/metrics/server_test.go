package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/skpr/api/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var srv = &Server{}

// --- pickStep ---

func TestPickStep(t *testing.T) {
	cases := []struct {
		name string
		d    time.Duration
		want time.Duration
	}{
		{"30m → 1m", 30 * time.Minute, time.Minute},
		{"2h exact → 1m", 2 * time.Hour, time.Minute},
		{"2h1s → 5m", 2*time.Hour + time.Second, 5 * time.Minute},
		{"12h exact → 5m", 12 * time.Hour, 5 * time.Minute},
		{"12h1s → 15m", 12*time.Hour + time.Second, 15 * time.Minute},
		{"48h exact → 15m", 48 * time.Hour, 15 * time.Minute},
		{"48h1s → 1h", 48*time.Hour + time.Second, time.Hour},
		{"7d exact → 1h", 7 * 24 * time.Hour, time.Hour},
		{"8d → 6h", 8 * 24 * time.Hour, 6 * time.Hour},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := pickStep(tc.d)
			if got != tc.want {
				t.Errorf("pickStep(%v) = %v, want %v", tc.d, got, tc.want)
			}
		})
	}
}

// --- AbsoluteRange validation ---

func TestAbsoluteRange_NilStartTime(t *testing.T) {
	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_ENVIRONMENT,
		Metric:    "cpu",
		StartTime: nil,
		EndTime:   timestamppb.New(time.Now()),
	}
	_, err := srv.AbsoluteRange(context.Background(), req)
	assertCode(t, err, codes.InvalidArgument)
}

func TestAbsoluteRange_NilEndTime(t *testing.T) {
	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_ENVIRONMENT,
		Metric:    "cpu",
		StartTime: timestamppb.New(time.Now().Add(-time.Hour)),
		EndTime:   nil,
	}
	_, err := srv.AbsoluteRange(context.Background(), req)
	assertCode(t, err, codes.InvalidArgument)
}

func TestAbsoluteRange_InvertedRange(t *testing.T) {
	now := time.Now()
	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_ENVIRONMENT,
		Metric:    "cpu",
		StartTime: timestamppb.New(now),
		EndTime:   timestamppb.New(now.Add(-time.Hour)),
	}
	_, err := srv.AbsoluteRange(context.Background(), req)
	assertCode(t, err, codes.InvalidArgument)
}

func TestAbsoluteRange_EqualStartEnd(t *testing.T) {
	now := time.Now()
	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_ENVIRONMENT,
		Metric:    "cpu",
		StartTime: timestamppb.New(now),
		EndTime:   timestamppb.New(now),
	}
	_, err := srv.AbsoluteRange(context.Background(), req)
	assertCode(t, err, codes.InvalidArgument)
}

func TestAbsoluteRange_UnknownMetric(t *testing.T) {
	now := time.Now()
	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_ENVIRONMENT,
		Metric:    "does_not_exist",
		StartTime: timestamppb.New(now.Add(-time.Hour)),
		EndTime:   timestamppb.New(now),
	}
	_, err := srv.AbsoluteRange(context.Background(), req)
	assertCode(t, err, codes.NotFound)
}

// --- AbsoluteRange happy-path / data correctness ---

// TestAbsoluteRange_PointCount checks that the number of returned points matches
// the expected count given the adaptive step.  For a 1-hour range the step is
// 1 minute.  The loop is inclusive of the end boundary, so we expect exactly
// floor((end-start)/step) + 1 points after truncation.
func TestAbsoluteRange_PointCount_1h(t *testing.T) {
	// Use a fixed, clean reference time so truncation is deterministic.
	start := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_ENVIRONMENT,
		Metric:    "cpu",
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	resp, err := srv.AbsoluteRange(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// step = 1m → 60 intervals → 61 points (inclusive)
	wantStep := time.Minute
	wantPoints := int(end.Sub(start)/wantStep) + 1
	if len(resp.Metrics) != wantPoints {
		t.Errorf("got %d points, want %d (step=%v)", len(resp.Metrics), wantPoints, wantStep)
	}
}

func TestAbsoluteRange_PointCount_6h(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(6 * time.Hour)

	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_CLUSTER,
		Metric:    "requests",
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	resp, err := srv.AbsoluteRange(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 6h is <= 12h → step = 5m → 72 intervals → 73 points
	wantStep := 5 * time.Minute
	wantPoints := int(end.Sub(start)/wantStep) + 1
	if len(resp.Metrics) != wantPoints {
		t.Errorf("got %d points, want %d (step=%v)", len(resp.Metrics), wantPoints, wantStep)
	}
}

func TestAbsoluteRange_PointCount_3d(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(3 * 24 * time.Hour)

	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_ENVIRONMENT,
		Metric:    "memory",
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	resp, err := srv.AbsoluteRange(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3d is <= 7d → step = 1h → 72 intervals → 73 points
	wantStep := time.Hour
	wantPoints := int(end.Sub(start)/wantStep) + 1
	if len(resp.Metrics) != wantPoints {
		t.Errorf("got %d points, want %d (step=%v)", len(resp.Metrics), wantPoints, wantStep)
	}
}

// TestAbsoluteRange_ValuesWithinRange checks that all returned values fall
// within [min, max] for the requested metric.
func TestAbsoluteRange_ValuesWithinRange(t *testing.T) {
	start := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(30 * time.Minute)

	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_ENVIRONMENT,
		Metric:    "cache_hit_rate",
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	resp, err := srv.AbsoluteRange(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Metrics) == 0 {
		t.Fatal("expected at least one metric point")
	}

	meta := metricData[pb.MetricType_ENVIRONMENT]["cache_hit_rate"]
	for i, m := range resp.Metrics {
		if m.Value == nil {
			t.Errorf("point %d: value is nil", i)
			continue
		}
		v := *m.Value
		if v < meta.Min || v > meta.Max {
			t.Errorf("point %d: value %v out of range [%v, %v]", i, v, meta.Min, meta.Max)
		}
	}
}

// TestAbsoluteRange_Deterministic verifies that two identical requests return
// identical values (no random drift).
func TestAbsoluteRange_Deterministic(t *testing.T) {
	start := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)

	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_ENVIRONMENT,
		Metric:    "cpu",
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}

	resp1, err := srv.AbsoluteRange(context.Background(), req)
	if err != nil {
		t.Fatalf("first call error: %v", err)
	}
	resp2, err := srv.AbsoluteRange(context.Background(), req)
	if err != nil {
		t.Fatalf("second call error: %v", err)
	}

	if len(resp1.Metrics) != len(resp2.Metrics) {
		t.Fatalf("length mismatch: %d vs %d", len(resp1.Metrics), len(resp2.Metrics))
	}
	for i := range resp1.Metrics {
		if *resp1.Metrics[i].Value != *resp2.Metrics[i].Value {
			t.Errorf("point %d: value %v != %v", i, *resp1.Metrics[i].Value, *resp2.Metrics[i].Value)
		}
	}
}

// TestAbsoluteRange_TimestampsMonotonic checks that timestamps in the response
// are strictly increasing.
func TestAbsoluteRange_TimestampsMonotonic(t *testing.T) {
	start := time.Date(2024, 1, 10, 6, 0, 0, 0, time.UTC)
	end := start.Add(90 * time.Minute)

	req := &pb.AbsoluteRangeRequest{
		Type:      pb.MetricType_ENVIRONMENT,
		Metric:    "requests",
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	resp, err := srv.AbsoluteRange(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 1; i < len(resp.Metrics); i++ {
		prev := resp.Metrics[i-1].Timestamp.AsTime()
		curr := resp.Metrics[i].Timestamp.AsTime()
		if !curr.After(prev) {
			t.Errorf("timestamps not strictly increasing at index %d: %v >= %v", i, curr, prev)
		}
	}
}

// --- helpers ---

func assertCode(t *testing.T, err error, want codes.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %v, got nil", want)
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got: %v", err)
	}
	if st.Code() != want {
		t.Errorf("got code %v, want %v (msg: %s)", st.Code(), want, st.Message())
	}
}
