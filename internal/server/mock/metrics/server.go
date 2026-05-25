package metrics

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/skpr/api/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements the GRPC "events" definition.
type Server struct {
	pb.UnimplementedMetricsServer
}

type MockMetric struct {
	Min    float64
	Max    float64
	Source pb.MetricSource
}

// metricData is built once at package initialization
var metricData = map[pb.MetricType]map[string]MockMetric{
	pb.MetricType_CLUSTER: {
		"requests":            {250_000, 4_000_000, pb.MetricSource_SYSTEM},
		"httpcode_target_200": {500, 1_000, pb.MetricSource_SYSTEM},
		"httpcode_target_300": {25, 50, pb.MetricSource_SYSTEM},
		"httpcode_target_400": {10, 25, pb.MetricSource_SYSTEM},
		"httpcode_target_500": {0, 10, pb.MetricSource_SYSTEM},
	},
	pb.MetricType_ENVIRONMENT: {
		"requests":              {25_000, 200_000, pb.MetricSource_SYSTEM},
		"cpu":                   {25, 100, pb.MetricSource_SYSTEM},
		"memory":                {512_000, 4_294_967_296, pb.MetricSource_SYSTEM}, // 512KB to 4GB
		"replicas":              {2, 8, pb.MetricSource_SYSTEM},
		"php_active":            {4, 48, pb.MetricSource_SYSTEM},
		"php_idle":              {2, 12, pb.MetricSource_SYSTEM},
		"php_queued":            {0, 8, pb.MetricSource_SYSTEM},
		"cache_hit_rate":        {60, 99, pb.MetricSource_SYSTEM},
		"cdn_download":          {50_000_000, 300_000_000, pb.MetricSource_SYSTEM},
		"cdn_upload":            {0, 20_000, pb.MetricSource_SYSTEM},
		"invalidation_paths":    {10, 50, pb.MetricSource_SYSTEM},
		"invalidation_requests": {5, 20, pb.MetricSource_SYSTEM},
		"origin_errors":         {0, 10, pb.MetricSource_SYSTEM},
		"httpcode_target_200":   {200, 500, pb.MetricSource_SYSTEM},
		"httpcode_target_300":   {25, 50, pb.MetricSource_SYSTEM},
		"httpcode_target_400":   {10, 25, pb.MetricSource_SYSTEM},
		"httpcode_target_500":   {0, 10, pb.MetricSource_SYSTEM},
		"response_times_avg":    {0.1, 0.25, pb.MetricSource_SYSTEM},
		"response_times_p95":    {2.0, 5.0, pb.MetricSource_SYSTEM},
		"response_times_p99":    {10.0, 20.0, pb.MetricSource_SYSTEM},
		"mock_application_1":    {0, 100, pb.MetricSource_APPLICATION},
		"mock_application_2":    {0.3, 1.4, pb.MetricSource_APPLICATION},
		"mock_application_3":    {0, 100_000_000, pb.MetricSource_APPLICATION},
	},
}

// deterministicRange returns a specific value consistently for a point in a time series.
func deterministicRange(t time.Time, minVal, maxVal float64, seconds int64, key string) float64 {
	h := fnv.New32a()

	bucketKey := fmt.Sprintf("%d-%s", t.Unix()/seconds, key)
	_, _ = h.Write([]byte(bucketKey))
	hashVal := float64(h.Sum32()) / float64(^uint32(0))

	rangeSize := maxVal - minVal
	return minVal + (hashVal * rangeSize)
}

func metricMappings(metricType pb.MetricType) map[string]MockMetric {
	return metricData[metricType]
}

// AvailableMetrics lists all available metrics.
func (s *Server) AvailableMetrics(ctx context.Context, req *pb.AvailableMetricsRequest) (*pb.AvailableMetricsResponse, error) {
	mappings := metricMappings(req.Type)

	availableMetrics := []*pb.MetricDefinition{}
	for key, metric := range mappings {
		// Hide application metrics for the prod environment.
		if req.Environment != nil && *req.Environment == "prod" && metric.Source == pb.MetricSource_APPLICATION {
			continue
		}

		availableMetrics = append(availableMetrics, &pb.MetricDefinition{
			Name:   key,
			Type:   req.Type,
			Title:  strings.ReplaceAll(key, "_", " "),
			Source: metric.Source,
		})
	}

	return &pb.AvailableMetricsResponse{
		Metrics: availableMetrics,
	}, nil
}

// pickStep mirrors the platform's NewTimeSeriesAutoStep, aiming for ~200
// points per period on a sensible boundary.
func pickStep(start, end time.Time) time.Duration {
	duration := end.Sub(start)

	var step time.Duration
	switch {
	case duration <= 45*time.Minute:
		step = 15 * time.Second
	case duration <= 90*time.Minute:
		step = 30 * time.Second
	case duration <= 3*time.Hour:
		step = time.Minute
	case duration <= 8*time.Hour:
		step = 2 * time.Minute
	case duration <= 16*time.Hour:
		step = 5 * time.Minute
	case duration <= 2*24*time.Hour:
		step = 10 * time.Minute
	case duration <= 4*24*time.Hour:
		step = 30 * time.Minute
	case duration <= 8*24*time.Hour:
		step = time.Hour
	case duration <= 21*24*time.Hour:
		step = 2 * time.Hour
	case duration <= 42*24*time.Hour:
		step = 4 * time.Hour
	default:
		step = 6 * time.Hour
	}

	return step
}

// AbsoluteRange gets a metric for a given timestamp range.
func (s *Server) AbsoluteRange(ctx context.Context, req *pb.AbsoluteRangeRequest) (*pb.AbsoluteRangeResponse, error) {
	// Validate start_time and end_time.
	if req.StartTime == nil || req.EndTime == nil {
		return nil, status.Errorf(codes.InvalidArgument, "start_time and end_time are required")
	}
	start := req.StartTime.AsTime()
	end := req.EndTime.AsTime()
	if start.IsZero() || end.IsZero() {
		return nil, status.Errorf(codes.InvalidArgument, "start_time and end_time are required")
	}
	if !end.After(start) {
		return nil, status.Errorf(codes.InvalidArgument, "end_time must be after start_time")
	}

	mappings := metricMappings(req.Type)
	if _, ok := mappings[req.Metric]; !ok {
		return nil, status.Errorf(codes.NotFound, "metric %s not found", req.Metric)
	}
	metricMin, metricMax := mappings[req.Metric].Min, mappings[req.Metric].Max

	step := pickStep(start, end)
	seedKey := fmt.Sprintf("%s_%s", req.Type, req.Metric)

	// Align the starting point to a clean step boundary so that bucket keys
	// remain stable regardless of the exact second the caller provides.
	metricTime := start.Truncate(step)

	output := []*pb.MetricValue{}
	for !metricTime.After(end) {
		v := deterministicRange(metricTime, metricMin, metricMax, int64(step.Seconds()), seedKey)
		output = append(output, &pb.MetricValue{
			Timestamp: timestamppb.New(metricTime),
			Value:     &v,
		})
		metricTime = metricTime.Add(step)
	}

	return &pb.AbsoluteRangeResponse{
		Metrics: output,
	}, nil
}
