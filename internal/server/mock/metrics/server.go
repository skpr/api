package metrics

import (
	"context"
	"fmt"
	"hash/fnv"
	"maps"
	"slices"
	"time"

	"github.com/skpr/api/pb"
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

func (s *Server) ClusterRequests(ctx context.Context, req *pb.ClusterRequestsRequest) (*pb.ClusterRequestsResponse, error) {
	return &pb.ClusterRequestsResponse{
		Metrics: []*pb.MetricClusterRequests{
			{Date: "2024-04-01 00:00", Requests: 1583870},
			{Date: "2024-04-01 00:02", Requests: 2563450},
			{Date: "2024-04-01 00:03", Requests: 1473176},
			{Date: "2024-04-01 00:04", Requests: 3200394},
			{Date: "2024-04-01 00:05", Requests: 4317112},
			{Date: "2024-04-01 00:06", Requests: 1949708},
			{Date: "2024-04-01 00:07", Requests: 2225637},
			{Date: "2024-04-01 00:08", Requests: 4044486},
			{Date: "2024-04-01 00:09", Requests: 1378839},
			{Date: "2024-04-01 00:10", Requests: 3708293},
			{Date: "2024-04-01 00:11", Requests: 1353318},
			{Date: "2024-04-01 00:12", Requests: 2809815},
			{Date: "2024-04-01 00:13", Requests: 2128851},
			{Date: "2024-04-01 00:14", Requests: 4526252},
			{Date: "2024-04-01 00:15", Requests: 4145685},
			{Date: "2024-04-01 00:16", Requests: 3494504},
			{Date: "2024-04-01 00:17", Requests: 4935480},
			{Date: "2024-04-01 00:18", Requests: 2863479},
			{Date: "2024-04-01 00:19", Requests: 4802475},
			{Date: "2024-04-01 00:20", Requests: 1027111},
			{Date: "2024-04-01 00:21", Requests: 4880666},
			{Date: "2024-04-01 00:22", Requests: 3752164},
			{Date: "2024-04-01 00:23", Requests: 2670388},
			{Date: "2024-04-01 00:24", Requests: 1073939},
			{Date: "2024-04-01 00:25", Requests: 2191549},
			{Date: "2024-04-01 00:26", Requests: 1485730},
			{Date: "2024-04-01 00:27", Requests: 2413299},
			{Date: "2024-04-01 00:28", Requests: 2185148},
			{Date: "2024-04-01 00:29", Requests: 1183772},
			{Date: "2024-04-01 00:30", Requests: 3227069},
		},
	}, nil
}

// ClusterResponseCodes mocked function call.
func (s *Server) ClusterResponseCodes(ctx context.Context, req *pb.ClusterResponseCodesRequest) (*pb.ClusterResponseCodesResponse, error) {
	return &pb.ClusterResponseCodesResponse{
		Metrics: []*pb.MetricClusterResponseCodes{
			{Date: "2024-04-01 00:00", Server: 1, Client: 15, Redirection: 17, Successful: 120},
			{Date: "2024-04-01 00:02", Server: 2, Client: 14, Redirection: 16, Successful: 124},
			{Date: "2024-04-01 00:03", Server: 3, Client: 13, Redirection: 15, Successful: 125},
			{Date: "2024-04-01 00:04", Server: 1, Client: 12, Redirection: 14, Successful: 120},
			{Date: "2024-04-01 00:05", Server: 2, Client: 15, Redirection: 17, Successful: 135},
			{Date: "2024-04-01 00:06", Server: 4, Client: 13, Redirection: 15, Successful: 250},
			{Date: "2024-04-01 00:07", Server: 9, Client: 12, Redirection: 14, Successful: 240},
			{Date: "2024-04-01 00:08", Server: 1, Client: 14, Redirection: 17, Successful: 230},
			{Date: "2024-04-01 00:09", Server: 0, Client: 11, Redirection: 12, Successful: 150},
			{Date: "2024-04-01 00:10", Server: 0, Client: 10, Redirection: 12, Successful: 300},
			{Date: "2024-04-01 00:11", Server: 2, Client: 12, Redirection: 11, Successful: 280},
			{Date: "2024-04-01 00:12", Server: 3, Client: 11, Redirection: 10, Successful: 287},
			{Date: "2024-04-01 00:13", Server: 1, Client: 14, Redirection: 13, Successful: 250},
			{Date: "2024-04-01 00:14", Server: 3, Client: 13, Redirection: 12, Successful: 155},
			{Date: "2024-04-01 00:15", Server: 2, Client: 12, Redirection: 11, Successful: 145},
			{Date: "2024-04-01 00:16", Server: 1, Client: 22, Redirection: 22, Successful: 120},
			{Date: "2024-04-01 00:17", Server: 5, Client: 21, Redirection: 18, Successful: 110},
			{Date: "2024-04-01 00:18", Server: 1, Client: 34, Redirection: 40, Successful: 130},
			{Date: "2024-04-01 00:19", Server: 6, Client: 51, Redirection: 12, Successful: 122},
			{Date: "2024-04-01 00:20", Server: 6, Client: 45, Redirection: 14, Successful: 144},
			{Date: "2024-04-01 00:21", Server: 1, Client: 12, Redirection: 13, Successful: 133},
			{Date: "2024-04-01 00:22", Server: 2, Client: 34, Redirection: 22, Successful: 98},
			{Date: "2024-04-01 00:23", Server: 0, Client: 41, Redirection: 12, Successful: 95},
			{Date: "2024-04-01 00:24", Server: 1, Client: 55, Redirection: 34, Successful: 97},
			{Date: "2024-04-01 00:25", Server: 9, Client: 11, Redirection: 15, Successful: 98},
			{Date: "2024-04-01 00:26", Server: 8, Client: 13, Redirection: 64, Successful: 99},
			{Date: "2024-04-01 00:27", Server: 5, Client: 10, Redirection: 55, Successful: 101},
			{Date: "2024-04-01 00:28", Server: 3, Client: 12, Redirection: 45, Successful: 100},
			{Date: "2024-04-01 00:29", Server: 2, Client: 13, Redirection: 35, Successful: 120},
			{Date: "2024-04-01 00:30", Server: 2, Client: 14, Redirection: 11, Successful: 115},
		},
	}, nil
}

// Requests mocked function call.
func (s *Server) Requests(ctx context.Context, req *pb.RequestsRequest) (*pb.RequestsResponse, error) {
	return &pb.RequestsResponse{
		Metrics: []*pb.MetricRequests{
			{Date: "2024-04-01 00:00", Requests: 1583870},
			{Date: "2024-04-01 00:02", Requests: 2563450},
			{Date: "2024-04-01 00:03", Requests: 1473176},
			{Date: "2024-04-01 00:04", Requests: 3200394},
			{Date: "2024-04-01 00:05", Requests: 4317112},
			{Date: "2024-04-01 00:06", Requests: 1949708},
			{Date: "2024-04-01 00:07", Requests: 2225637},
			{Date: "2024-04-01 00:08", Requests: 4044486},
			{Date: "2024-04-01 00:09", Requests: 1378839},
			{Date: "2024-04-01 00:10", Requests: 3708293},
			{Date: "2024-04-01 00:11", Requests: 1353318},
			{Date: "2024-04-01 00:12", Requests: 2809815},
			{Date: "2024-04-01 00:13", Requests: 2128851},
			{Date: "2024-04-01 00:14", Requests: 4526252},
			{Date: "2024-04-01 00:15", Requests: 4145685},
			{Date: "2024-04-01 00:16", Requests: 3494504},
			{Date: "2024-04-01 00:17", Requests: 4935480},
			{Date: "2024-04-01 00:18", Requests: 2863479},
			{Date: "2024-04-01 00:19", Requests: 4802475},
			{Date: "2024-04-01 00:20", Requests: 1027111},
			{Date: "2024-04-01 00:21", Requests: 4880666},
			{Date: "2024-04-01 00:22", Requests: 3752164},
			{Date: "2024-04-01 00:23", Requests: 2670388},
			{Date: "2024-04-01 00:24", Requests: 1073939},
			{Date: "2024-04-01 00:25", Requests: 2191549},
			{Date: "2024-04-01 00:26", Requests: 1485730},
			{Date: "2024-04-01 00:27", Requests: 2413299},
			{Date: "2024-04-01 00:28", Requests: 2185148},
			{Date: "2024-04-01 00:29", Requests: 1183772},
			{Date: "2024-04-01 00:30", Requests: 3227069},
		},
	}, nil
}

// ResponseCodes mocked function call.
func (s *Server) ResponseCodes(ctx context.Context, req *pb.ResponseCodesRequest) (*pb.ResponseCodesResponse, error) {
	return &pb.ResponseCodesResponse{
		Metrics: []*pb.MetricResponseCodes{
			{Date: "2024-04-01 00:00", Server: 1, Client: 15, Redirection: 17, Successful: 120},
			{Date: "2024-04-01 00:02", Server: 2, Client: 14, Redirection: 16, Successful: 124},
			{Date: "2024-04-01 00:03", Server: 3, Client: 13, Redirection: 15, Successful: 125},
			{Date: "2024-04-01 00:04", Server: 1, Client: 12, Redirection: 14, Successful: 120},
			{Date: "2024-04-01 00:05", Server: 2, Client: 15, Redirection: 17, Successful: 135},
			{Date: "2024-04-01 00:06", Server: 4, Client: 13, Redirection: 15, Successful: 250},
			{Date: "2024-04-01 00:07", Server: 9, Client: 12, Redirection: 14, Successful: 240},
			{Date: "2024-04-01 00:08", Server: 1, Client: 14, Redirection: 17, Successful: 230},
			{Date: "2024-04-01 00:09", Server: 0, Client: 11, Redirection: 12, Successful: 150},
			{Date: "2024-04-01 00:10", Server: 0, Client: 10, Redirection: 12, Successful: 300},
			{Date: "2024-04-01 00:11", Server: 2, Client: 12, Redirection: 11, Successful: 280},
			{Date: "2024-04-01 00:12", Server: 3, Client: 11, Redirection: 10, Successful: 287},
			{Date: "2024-04-01 00:13", Server: 1, Client: 14, Redirection: 13, Successful: 250},
			{Date: "2024-04-01 00:14", Server: 3, Client: 13, Redirection: 12, Successful: 155},
			{Date: "2024-04-01 00:15", Server: 2, Client: 12, Redirection: 11, Successful: 145},
			{Date: "2024-04-01 00:16", Server: 1, Client: 22, Redirection: 22, Successful: 120},
			{Date: "2024-04-01 00:17", Server: 5, Client: 21, Redirection: 18, Successful: 110},
			{Date: "2024-04-01 00:18", Server: 1, Client: 34, Redirection: 40, Successful: 130},
			{Date: "2024-04-01 00:19", Server: 6, Client: 51, Redirection: 12, Successful: 122},
			{Date: "2024-04-01 00:20", Server: 6, Client: 45, Redirection: 14, Successful: 144},
			{Date: "2024-04-01 00:21", Server: 1, Client: 12, Redirection: 13, Successful: 133},
			{Date: "2024-04-01 00:22", Server: 2, Client: 34, Redirection: 22, Successful: 98},
			{Date: "2024-04-01 00:23", Server: 0, Client: 41, Redirection: 12, Successful: 95},
			{Date: "2024-04-01 00:24", Server: 1, Client: 55, Redirection: 34, Successful: 97},
			{Date: "2024-04-01 00:25", Server: 9, Client: 11, Redirection: 15, Successful: 98},
			{Date: "2024-04-01 00:26", Server: 8, Client: 13, Redirection: 64, Successful: 99},
			{Date: "2024-04-01 00:27", Server: 5, Client: 10, Redirection: 55, Successful: 101},
			{Date: "2024-04-01 00:28", Server: 3, Client: 12, Redirection: 45, Successful: 100},
			{Date: "2024-04-01 00:29", Server: 2, Client: 13, Redirection: 35, Successful: 120},
			{Date: "2024-04-01 00:30", Server: 2, Client: 14, Redirection: 11, Successful: 115},
		},
	}, nil
}

// ResponseTimes mocked function call.
func (s *Server) ResponseTimes(ctx context.Context, req *pb.ResponseTimesRequest) (*pb.ResponseTimesResponse, error) {
	return &pb.ResponseTimesResponse{
		Metrics: []*pb.MetricResponseTimes{
			{Date: "2024-04-01 00:00", Average: 190, P95: 480, P99: 1600},
			{Date: "2024-04-01 00:02", Average: 205, P95: 500, P99: 2250},
			{Date: "2024-04-01 00:03", Average: 198, P95: 520, P99: 1350},
			{Date: "2024-04-01 00:04", Average: 210, P95: 560, P99: 3000},
			{Date: "2024-04-01 00:05", Average: 195, P95: 530, P99: 2000},
			{Date: "2024-04-01 00:06", Average: 200, P95: 550, P99: 3750},
			{Date: "2024-04-01 00:07", Average: 185, P95: 490, P99: 1150},
			{Date: "2024-04-01 00:08", Average: 215, P95: 580, P99: 4600},
			{Date: "2024-04-01 00:09", Average: 192, P95: 470, P99: 1550},
			{Date: "2024-04-01 00:10", Average: 205, P95: 510, P99: 2650},
			{Date: "2024-04-01 00:11", Average: 208, P95: 540, P99: 3550},
			{Date: "2024-04-01 00:12", Average: 195, P95: 520, P99: 2250},
			{Date: "2024-04-01 00:13", Average: 212, P95: 560, P99: 4100},
			{Date: "2024-04-01 00:14", Average: 187, P95: 490, P99: 1350},
			{Date: "2024-04-01 00:15", Average: 200, P95: 500, P99: 1650},
			{Date: "2024-04-01 00:16", Average: 193, P95: 510, P99: 1450},
			{Date: "2024-04-01 00:17", Average: 218, P95: 600, P99: 4900},
			{Date: "2024-04-01 00:18", Average: 210, P95: 570, P99: 4250},
			{Date: "2024-04-01 00:19", Average: 190, P95: 500, P99: 1700},
			{Date: "2024-04-01 00:20", Average: 185, P95: 470, P99: 1300},
			{Date: "2024-04-01 00:21", Average: 198, P95: 500, P99: 2100},
			{Date: "2024-04-01 00:22", Average: 192, P95: 480, P99: 1850},
			{Date: "2024-04-01 00:23", Average: 205, P95: 530, P99: 2300},
			{Date: "2024-04-01 00:24", Average: 215, P95: 590, P99: 4550},
			{Date: "2024-04-01 00:25", Average: 200, P95: 520, P99: 2500},
			{Date: "2024-04-01 00:26", Average: 195, P95: 480, P99: 1600},
			{Date: "2024-04-01 00:27", Average: 220, P95: 610, P99: 4800},
			{Date: "2024-04-01 00:28", Average: 190, P95: 490, P99: 1700},
			{Date: "2024-04-01 00:29", Average: 208, P95: 550, P99: 3550},
			{Date: "2024-04-01 00:30", Average: 215, P95: 580, P99: 4350},
		},
	}, nil
}

// CacheRatio mocked function call.
func (s *Server) CacheRatio(ctx context.Context, req *pb.CacheRatioRequest) (*pb.CacheRatioResponse, error) {
	return &pb.CacheRatioResponse{
		Metrics: []*pb.MetricCacheRatio{
			{Date: "2024-04-01 00:00", Hit: 65, Miss: 35},
			{Date: "2024-04-01 00:02", Hit: 84, Miss: 16},
			{Date: "2024-04-01 00:03", Hit: 82, Miss: 18},
			{Date: "2024-04-01 00:04", Hit: 82, Miss: 18},
			{Date: "2024-04-01 00:05", Hit: 77, Miss: 23},
			{Date: "2024-04-01 00:06", Hit: 72, Miss: 28},
			{Date: "2024-04-01 00:07", Hit: 63, Miss: 37},
			{Date: "2024-04-01 00:08", Hit: 68, Miss: 32},
			{Date: "2024-04-01 00:09", Hit: 86, Miss: 14},
			{Date: "2024-04-01 00:10", Hit: 69, Miss: 31},
			{Date: "2024-04-01 00:11", Hit: 94, Miss: 6},
			{Date: "2024-04-01 00:12", Hit: 95, Miss: 5},
			{Date: "2024-04-01 00:13", Hit: 75, Miss: 25},
			{Date: "2024-04-01 00:14", Hit: 69, Miss: 31},
			{Date: "2024-04-01 00:15", Hit: 92, Miss: 8},
			{Date: "2024-04-01 00:16", Hit: 80, Miss: 20},
			{Date: "2024-04-01 00:17", Hit: 90, Miss: 10},
			{Date: "2024-04-01 00:18", Hit: 98, Miss: 2},
			{Date: "2024-04-01 00:19", Hit: 81, Miss: 19},
			{Date: "2024-04-01 00:20", Hit: 78, Miss: 22},
			{Date: "2024-04-01 00:21", Hit: 94, Miss: 6},
			{Date: "2024-04-01 00:22", Hit: 72, Miss: 28},
			{Date: "2024-04-01 00:23", Hit: 98, Miss: 2},
			{Date: "2024-04-01 00:24", Hit: 80, Miss: 20},
			{Date: "2024-04-01 00:25", Hit: 89, Miss: 11},
			{Date: "2024-04-01 00:26", Hit: 64, Miss: 36},
			{Date: "2024-04-01 00:27", Hit: 97, Miss: 3},
			{Date: "2024-04-01 00:28", Hit: 99, Miss: 1},
			{Date: "2024-04-01 00:29", Hit: 65, Miss: 35},
			{Date: "2024-04-01 00:30", Hit: 92, Miss: 8},
		},
	}, nil
}

// OriginErrors mocked function call.
func (s *Server) OriginErrors(ctx context.Context, req *pb.OriginErrorsRequest) (*pb.OriginErrorsResponse, error) {
	return &pb.OriginErrorsResponse{
		Metrics: []*pb.MetricOriginErrors{
			{Date: "2024-04-01 00:00", Errors: 15},
			{Date: "2024-04-01 00:02", Errors: 0},
			{Date: "2024-04-01 00:03", Errors: 19},
			{Date: "2024-04-01 00:04", Errors: 7},
			{Date: "2024-04-01 00:05", Errors: 19},
			{Date: "2024-04-01 00:06", Errors: 0},
			{Date: "2024-04-01 00:07", Errors: 0},
			{Date: "2024-04-01 00:08", Errors: 8},
			{Date: "2024-04-01 00:09", Errors: 20},
			{Date: "2024-04-01 00:10", Errors: 3},
			{Date: "2024-04-01 00:11", Errors: 18},
			{Date: "2024-04-01 00:12", Errors: 4},
			{Date: "2024-04-01 00:13", Errors: 2},
			{Date: "2024-04-01 00:14", Errors: 10},
			{Date: "2024-04-01 00:15", Errors: 18},
			{Date: "2024-04-01 00:16", Errors: 19},
			{Date: "2024-04-01 00:17", Errors: 12},
			{Date: "2024-04-01 00:18", Errors: 10},
			{Date: "2024-04-01 00:19", Errors: 15},
			{Date: "2024-04-01 00:20", Errors: 5},
			{Date: "2024-04-01 00:21", Errors: 1},
			{Date: "2024-04-01 00:22", Errors: 8},
			{Date: "2024-04-01 00:23", Errors: 6},
			{Date: "2024-04-01 00:24", Errors: 6},
			{Date: "2024-04-01 00:25", Errors: 7},
			{Date: "2024-04-01 00:26", Errors: 11},
			{Date: "2024-04-01 00:27", Errors: 18},
			{Date: "2024-04-01 00:28", Errors: 18},
			{Date: "2024-04-01 00:29", Errors: 11},
			{Date: "2024-04-01 00:30", Errors: 3},
		},
	}, nil
}

// InvalidationRequests mocked function call.
func (s *Server) InvalidationRequests(ctx context.Context, req *pb.InvalidationRequestsRequest) (*pb.InvalidationRequestsResponse, error) {
	return &pb.InvalidationRequestsResponse{
		Metrics: []*pb.MetricInvalidationRequests{
			{Date: "2024-04-01 00:00", Requests: 15},
			{Date: "2024-04-01 00:02", Requests: 0},
			{Date: "2024-04-01 00:03", Requests: 19},
			{Date: "2024-04-01 00:04", Requests: 7},
			{Date: "2024-04-01 00:05", Requests: 19},
			{Date: "2024-04-01 00:06", Requests: 0},
			{Date: "2024-04-01 00:07", Requests: 0},
			{Date: "2024-04-01 00:08", Requests: 8},
			{Date: "2024-04-01 00:09", Requests: 20},
			{Date: "2024-04-01 00:10", Requests: 3},
			{Date: "2024-04-01 00:11", Requests: 18},
			{Date: "2024-04-01 00:12", Requests: 4},
			{Date: "2024-04-01 00:13", Requests: 2},
			{Date: "2024-04-01 00:14", Requests: 10},
			{Date: "2024-04-01 00:15", Requests: 18},
			{Date: "2024-04-01 00:16", Requests: 19},
			{Date: "2024-04-01 00:17", Requests: 12},
			{Date: "2024-04-01 00:18", Requests: 10},
			{Date: "2024-04-01 00:19", Requests: 15},
			{Date: "2024-04-01 00:20", Requests: 5},
			{Date: "2024-04-01 00:21", Requests: 1},
			{Date: "2024-04-01 00:22", Requests: 8},
			{Date: "2024-04-01 00:23", Requests: 6},
			{Date: "2024-04-01 00:24", Requests: 6},
			{Date: "2024-04-01 00:25", Requests: 7},
			{Date: "2024-04-01 00:26", Requests: 11},
			{Date: "2024-04-01 00:27", Requests: 18},
			{Date: "2024-04-01 00:28", Requests: 18},
			{Date: "2024-04-01 00:29", Requests: 11},
			{Date: "2024-04-01 00:30", Requests: 3},
		},
	}, nil
}

// InvalidationPaths mocked function call.
func (s *Server) InvalidationPaths(ctx context.Context, req *pb.InvalidationPathsRequest) (*pb.InvalidationPathsResponse, error) {
	return &pb.InvalidationPathsResponse{
		Metrics: []*pb.MetricInvalidationPaths{
			{Date: "2024-04-01 00:00", Paths: 15},
			{Date: "2024-04-01 00:02", Paths: 0},
			{Date: "2024-04-01 00:03", Paths: 19},
			{Date: "2024-04-01 00:04", Paths: 7},
			{Date: "2024-04-01 00:05", Paths: 19},
			{Date: "2024-04-01 00:06", Paths: 0},
			{Date: "2024-04-01 00:07", Paths: 0},
			{Date: "2024-04-01 00:08", Paths: 8},
			{Date: "2024-04-01 00:09", Paths: 20},
			{Date: "2024-04-01 00:10", Paths: 3},
			{Date: "2024-04-01 00:11", Paths: 18},
			{Date: "2024-04-01 00:12", Paths: 4},
			{Date: "2024-04-01 00:13", Paths: 2},
			{Date: "2024-04-01 00:14", Paths: 10},
			{Date: "2024-04-01 00:15", Paths: 18},
			{Date: "2024-04-01 00:16", Paths: 19},
			{Date: "2024-04-01 00:17", Paths: 12},
			{Date: "2024-04-01 00:18", Paths: 10},
			{Date: "2024-04-01 00:19", Paths: 15},
			{Date: "2024-04-01 00:20", Paths: 5},
			{Date: "2024-04-01 00:21", Paths: 1},
			{Date: "2024-04-01 00:22", Paths: 8},
			{Date: "2024-04-01 00:23", Paths: 6},
			{Date: "2024-04-01 00:24", Paths: 6},
			{Date: "2024-04-01 00:25", Paths: 7},
			{Date: "2024-04-01 00:26", Paths: 11},
			{Date: "2024-04-01 00:27", Paths: 18},
			{Date: "2024-04-01 00:28", Paths: 18},
			{Date: "2024-04-01 00:29", Paths: 11},
			{Date: "2024-04-01 00:30", Paths: 3},
		},
	}, nil
}

// ResourceUsage mocked function call.
func (s *Server) ResourceUsage(ctx context.Context, req *pb.ResourceUsageRequest) (*pb.ResourceUsageResponse, error) {
	return &pb.ResourceUsageResponse{
		Metrics: []*pb.MetricResourceUsage{
			{Date: "2024-04-01 00:00", CPU: 500, Memory: 1024, Replicas: 3, ActiveProcesses: 10, IdleProcesses: 2, ListenQueue: 0},
			{Date: "2024-04-01 00:05", CPU: 550, Memory: 1034, Replicas: 4, ActiveProcesses: 11, IdleProcesses: 3, ListenQueue: 1},
			{Date: "2024-04-01 00:10", CPU: 600, Memory: 1044, Replicas: 5, ActiveProcesses: 12, IdleProcesses: 4, ListenQueue: 2},
			{Date: "2024-04-01 00:15", CPU: 650, Memory: 1054, Replicas: 3, ActiveProcesses: 13, IdleProcesses: 2, ListenQueue: 3},
			{Date: "2024-04-01 00:20", CPU: 700, Memory: 1064, Replicas: 4, ActiveProcesses: 14, IdleProcesses: 3, ListenQueue: 0},
			{Date: "2024-04-01 00:25", CPU: 750, Memory: 1074, Replicas: 5, ActiveProcesses: 10, IdleProcesses: 4, ListenQueue: 1},
			{Date: "2024-04-01 00:30", CPU: 800, Memory: 1084, Replicas: 3, ActiveProcesses: 11, IdleProcesses: 2, ListenQueue: 2},
			{Date: "2024-04-01 00:35", CPU: 850, Memory: 1094, Replicas: 4, ActiveProcesses: 12, IdleProcesses: 3, ListenQueue: 3},
			{Date: "2024-04-01 00:40", CPU: 900, Memory: 1104, Replicas: 5, ActiveProcesses: 13, IdleProcesses: 4, ListenQueue: 0},
			{Date: "2024-04-01 00:45", CPU: 950, Memory: 1114, Replicas: 3, ActiveProcesses: 14, IdleProcesses: 2, ListenQueue: 1},
			{Date: "2024-04-01 00:50", CPU: 500, Memory: 1124, Replicas: 4, ActiveProcesses: 10, IdleProcesses: 3, ListenQueue: 2},
			{Date: "2024-04-01 00:55", CPU: 550, Memory: 1134, Replicas: 5, ActiveProcesses: 11, IdleProcesses: 4, ListenQueue: 3},
			{Date: "2024-04-01 01:00", CPU: 600, Memory: 1144, Replicas: 3, ActiveProcesses: 12, IdleProcesses: 2, ListenQueue: 0},
			{Date: "2024-04-01 01:05", CPU: 650, Memory: 1154, Replicas: 4, ActiveProcesses: 13, IdleProcesses: 3, ListenQueue: 1},
			{Date: "2024-04-01 01:10", CPU: 700, Memory: 1164, Replicas: 5, ActiveProcesses: 14, IdleProcesses: 4, ListenQueue: 2},
			{Date: "2024-04-01 01:15", CPU: 750, Memory: 1174, Replicas: 3, ActiveProcesses: 10, IdleProcesses: 2, ListenQueue: 3},
			{Date: "2024-04-01 01:20", CPU: 800, Memory: 1184, Replicas: 4, ActiveProcesses: 11, IdleProcesses: 3, ListenQueue: 0},
			{Date: "2024-04-01 01:25", CPU: 850, Memory: 1194, Replicas: 5, ActiveProcesses: 12, IdleProcesses: 4, ListenQueue: 1},
			{Date: "2024-04-01 01:30", CPU: 900, Memory: 1204, Replicas: 3, ActiveProcesses: 13, IdleProcesses: 2, ListenQueue: 2},
			{Date: "2024-04-01 01:35", CPU: 950, Memory: 1214, Replicas: 4, ActiveProcesses: 14, IdleProcesses: 3, ListenQueue: 3},
			{Date: "2024-04-01 01:40", CPU: 500, Memory: 1224, Replicas: 5, ActiveProcesses: 10, IdleProcesses: 4, ListenQueue: 0},
			{Date: "2024-04-01 01:45", CPU: 550, Memory: 1234, Replicas: 3, ActiveProcesses: 11, IdleProcesses: 2, ListenQueue: 1},
			{Date: "2024-04-01 01:50", CPU: 600, Memory: 1244, Replicas: 4, ActiveProcesses: 12, IdleProcesses: 3, ListenQueue: 2},
			{Date: "2024-04-01 01:55", CPU: 650, Memory: 1254, Replicas: 5, ActiveProcesses: 13, IdleProcesses: 4, ListenQueue: 3},
			{Date: "2024-04-01 02:00", CPU: 700, Memory: 1264, Replicas: 3, ActiveProcesses: 14, IdleProcesses: 2, ListenQueue: 0},
			{Date: "2024-04-01 02:05", CPU: 750, Memory: 1274, Replicas: 4, ActiveProcesses: 10, IdleProcesses: 3, ListenQueue: 1},
			{Date: "2024-04-01 02:10", CPU: 800, Memory: 1284, Replicas: 5, ActiveProcesses: 11, IdleProcesses: 4, ListenQueue: 2},
			{Date: "2024-04-01 02:15", CPU: 850, Memory: 1294, Replicas: 3, ActiveProcesses: 12, IdleProcesses: 2, ListenQueue: 3},
			{Date: "2024-04-01 02:20", CPU: 900, Memory: 1304, Replicas: 4, ActiveProcesses: 13, IdleProcesses: 3, ListenQueue: 0},
			{Date: "2024-04-01 02:25", CPU: 950, Memory: 1314, Replicas: 5, ActiveProcesses: 14, IdleProcesses: 4, ListenQueue: 1},
			{Date: "2024-04-01 02:30", CPU: 500, Memory: 1324, Replicas: 3, ActiveProcesses: 10, IdleProcesses: 2, ListenQueue: 2},
			{Date: "2024-04-01 02:35", CPU: 550, Memory: 1334, Replicas: 4, ActiveProcesses: 11, IdleProcesses: 3, ListenQueue: 3},
			{Date: "2024-04-01 02:40", CPU: 600, Memory: 1344, Replicas: 5, ActiveProcesses: 12, IdleProcesses: 4, ListenQueue: 0},
			{Date: "2024-04-01 02:45", CPU: 650, Memory: 1354, Replicas: 3, ActiveProcesses: 13, IdleProcesses: 2, ListenQueue: 1},
			{Date: "2024-04-01 02:50", CPU: 700, Memory: 1364, Replicas: 4, ActiveProcesses: 14, IdleProcesses: 3, ListenQueue: 2},
			{Date: "2024-04-01 02:55", CPU: 750, Memory: 1374, Replicas: 5, ActiveProcesses: 10, IdleProcesses: 4, ListenQueue: 3},
		},
	}, nil
}
