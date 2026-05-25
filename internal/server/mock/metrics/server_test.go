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

// TestPickStep exercises every tier boundary of pickStep, which computes
// duration as (end - 1m) - start to mirror the platform's
// NewTimeSeriesAutoStep.  The table lists the requested (end - start) span
// and the resulting effective duration_used = span - 1m.
func TestPickStep(t *testing.T) {
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		span time.Duration // end - start supplied to pickStep
		want time.Duration
	}{
		// <= 30m tier (15s)
		{"10m → 15s", 10 * time.Minute, 15 * time.Second},
		{"31m → 15s (boundary)", 31 * time.Minute, 15 * time.Second},
		// <= 1h tier (30s)
		{"32m → 30s", 32 * time.Minute, 30 * time.Second},
		{"1h1m → 30s (boundary)", 1*time.Hour + time.Minute, 30 * time.Second},
		// <= 3h tier (1m)
		{"1h2m → 1m", 1*time.Hour + 2*time.Minute, time.Minute},
		{"3h1m → 1m (boundary)", 3*time.Hour + time.Minute, time.Minute},
		// <= 12h tier (5m)
		{"3h2m → 5m", 3*time.Hour + 2*time.Minute, 5 * time.Minute},
		{"12h1m → 5m (boundary)", 12*time.Hour + time.Minute, 5 * time.Minute},
		// <= 24h tier (10m)
		{"12h2m → 10m", 12*time.Hour + 2*time.Minute, 10 * time.Minute},
		{"24h1m → 10m (boundary)", 24*time.Hour + time.Minute, 10 * time.Minute},
		// <= 72h tier (30m)
		{"24h2m → 30m", 24*time.Hour + 2*time.Minute, 30 * time.Minute},
		{"72h1m → 30m (boundary)", 72*time.Hour + time.Minute, 30 * time.Minute},
		// > 72h fallthrough (5m — matches platform)
		{"72h2m → 5m (fallthrough)", 72*time.Hour + 2*time.Minute, 5 * time.Minute},
		{"30d → 5m (fallthrough)", 30 * 24 * time.Hour, 5 * time.Minute},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := pickStep(ref, ref.Add(tc.span))
			if got != tc.want {
				t.Errorf("pickStep(span=%v) = %v, want %v", tc.span, got, tc.want)
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
// the expected count given the adaptive step.  For a 1-hour range the
// effective duration_used = 59m <= 1h, so the step is 30s.  The loop is
// inclusive of the end boundary, so we expect exactly
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

	// duration_used = 59m <= 1h → step = 30s → 120 intervals → 121 points (inclusive)
	wantStep := 30 * time.Second
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

	// duration_used = 71h59m <= 72h → step = 30m → 144 intervals → 145 points
	wantStep := 30 * time.Minute
	wantPoints := int(end.Sub(start)/wantStep) + 1
	if len(resp.Metrics) != wantPoints {
		t.Errorf("got %d points, want %d (step=%v)", len(resp.Metrics), wantPoints, wantStep)
	}
}

func TestAbsoluteRange_PointCount_30d(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(30 * 24 * time.Hour)

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

	// duration_used = 30d - 1m > 72h → fallthrough → step = 5m → 8640 intervals → 8641 points
	wantStep := 5 * time.Minute
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
