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
	Min         float64
	Max         float64
	Application bool
}

// metricData is built once at package initialization
var metricData = map[pb.MetricType]map[string]MockMetric{
	pb.MetricType_CLUSTER: {
		"requests":            {250_000, 4_000_000, false},
		"httpcode_target_200": {500, 1_000, false},
		"httpcode_target_300": {25, 50, false},
		"httpcode_target_400": {10, 25, false},
		"httpcode_target_500": {0, 10, false},
	},
	pb.MetricType_ENVIRONMENT: {
		"requests":              {25_000, 200_000, false},
		"cpu":                   {25, 100, false},
		"memory":                {512_000, 4_294_967_296, false}, // 512KB to 4GB
		"replicas":              {2, 8, false},
		"php_active":            {4, 48, false},
		"php_idle":              {2, 12, false},
		"php_queued":            {0, 8, false},
		"cache_hit_rate":        {60, 99, false},
		"invalidation_paths":    {10, 50, false},
		"invalidation_requests": {5, 20, false},
		"origin_errors":         {0, 10, false},
		"httpcode_target_200":   {200, 500, false},
		"httpcode_target_300":   {25, 50, false},
		"httpcode_target_400":   {10, 25, false},
		"httpcode_target_500":   {0, 10, false},
		"response_times_avg":    {0.1, 0.25, false},
		"response_times_p95":    {2.0, 5.0, false},
		"response_times_p99":    {10.0, 20.0, false},
		"mock_application_1":    {0, 100, true},
		"mock_application_2":    {0.3, 1.4, true},
		"mock_application_3":    {0, 100_000_000, true},
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
		availableMetrics = append(availableMetrics, &pb.MetricDefinition{
			Name:        key,
			Type:        req.Type,
			Title:       strings.ReplaceAll(key, "_", " "),
			Application: metric.Application,
		})
	}

	return &pb.AvailableMetricsResponse{
		Metrics: availableMetrics,
	}, nil
}

// AbsoluteRange gets a metric for a given timestamp range.
func (s *Server) AbsoluteRange(ctx context.Context, req *pb.AbsoluteRangeRequest) (*pb.AbsoluteRangeResponse, error) {
	mappings := metricMappings(req.Type)
	if _, ok := mappings[req.Metric]; !ok {
		return nil, status.Errorf(codes.NotFound, "metric %s not found", req.Metric)
	}
	metricMin, metricMax := mappings[req.Metric].Min, mappings[req.Metric].Max

	output := []*pb.MetricValue{}
	seedKey := fmt.Sprintf("%s_%s", req.Type, req.Metric)
	metricTime := req.StartTime.AsTime()
	for metricTime.Before(req.EndTime.AsTime()) {
		metric := pb.MetricValue{
			Timestamp: timestamppb.New(metricTime),
			Value:     new(deterministicRange(metricTime, metricMin, metricMax, 60, seedKey)),
		}
		output = append(output, &metric)
		metricTime = metricTime.Add(time.Minute)
	}

	return &pb.AbsoluteRangeResponse{
		Metrics: output,
	}, nil
}
