package metrics

import (
	"context"
	"fmt"
	"hash/fnv"
	"maps"
	"slices"
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

func deterministicRange(t time.Time, minVal, maxVal, seconds int64, key string) int64 {
	h := fnv.New32a()

	bucketKey := fmt.Sprintf("%d-%s", t.Unix()/seconds, key)
	_, _ = h.Write([]byte(bucketKey))
	hashVal := int64(h.Sum32())

	rangeSize := maxVal - minVal + 1
	return minVal + (hashVal % rangeSize)
}

func metricMappings(metricType pb.MetricType) map[string][]int64 {
	data := map[pb.MetricType]map[string][]int64{
		pb.MetricType_CLUSTER: {
			"requests":            {250_000, 4_000_000},
			"httpcode_target_200": {500, 1_000},
			"httpcode_target_300": {25, 50},
			"httpcode_target_400": {10, 25},
			"httpcode_target_500": {0, 10},
		},
		pb.MetricType_ENVIRONMENT: {
			"requests":              {25_000, 200_000},
			"cpu":                   {25, 100},
			"memory":                {512, 4096},
			"replicas":              {2, 8},
			"php_active":            {4, 48},
			"php_idle":              {2, 12},
			"php_queued":            {0, 8},
			"cache_hit_rate":        {60, 99},
			"invalidation_paths":    {10, 50},
			"invalidation_requests": {5, 20},
			"origin_errors":         {0, 10},
			"httpcode_target_200":   {200, 500},
			"httpcode_target_300":   {25, 50},
			"httpcode_target_400":   {10, 25},
			"httpcode_target_500":   {0, 10},
			"response_times_avg":    {100, 250},
			"response_times_p95":    {2_000, 5_000},
			"response_times_p99":    {10_000, 20_000},
		},
	}
	return data[metricType]
}

// AvailableMetrics lists all available metrics.
func (s *Server) AvailableMetrics(ctx context.Context, req *pb.AvailableMetricsRequest) (*pb.AvailableMetricsResponse, error) {
	mappings := metricMappings(req.Type)
	keys := slices.Collect(maps.Keys(mappings))

	metrics := make([]*pb.MetricDefinition, len(keys))
	for i, key := range keys {
		metrics[i] = &pb.MetricDefinition{
			Name: key,
			Type: req.Type,
		}
	}

	return &pb.AvailableMetricsResponse{
		Metrics: metrics,
	}, nil
}

// AbsoluteRange gets a metric for a given timestamp range.
func (s *Server) AbsoluteRange(ctx context.Context, req *pb.AbsoluteRangeRequest) (*pb.AbsoluteRangeResponse, error) {
	mappings := metricMappings(req.Type)
	if _, ok := mappings[req.Metric]; !ok {
		return nil, status.Errorf(codes.NotFound, "metric %s not found", req.Metric)
	}
	metricMin, metricMax := mappings[req.Metric][0], mappings[req.Metric][1]

	output := []*pb.MetricValue{}
	metricTime := req.StartTime.AsTime()
	for metricTime.Before(req.EndTime.AsTime()) {
		cacheKey := fmt.Sprintf("%s_%s", req.Type, req.Metric)
		metric := pb.MetricValue{
			Timestamp: timestamppb.New(metricTime),
			Value:     deterministicRange(metricTime, metricMin, metricMax, 60, cacheKey),
		}
		output = append(output, &metric)
		metricTime = metricTime.Add(time.Minute)
	}

	return &pb.AbsoluteRangeResponse{
		Metrics: output,
	}, nil
}
