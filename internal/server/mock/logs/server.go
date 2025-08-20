package logs

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/skpr/api/pb"
)

const (
	StreamNginx = "nginx"
	StreamFPM   = "fpm"
)

// Server implements the GRPC "events" definition.
type Server struct {
	pb.UnimplementedLogsServer
}

func (s *Server) ListStreams(ctx context.Context, req *pb.LogListStreamsRequest) (*pb.LogListStreamsResponse, error) {
	return &pb.LogListStreamsResponse{
		Default: StreamFPM,
		Streams: []string{
			StreamNginx,
			StreamFPM,
		},
	}, nil
}

func (s *Server) Tail(req *pb.LogTailRequest, server pb.Logs_TailServer) error {
	if req.Environment == "" {
		return fmt.Errorf("environment not provided")
	}

	if req.Stream == "" {
		return fmt.Errorf("stream not provided")
	}

	if !slices.Contains([]string{StreamNginx, StreamFPM}, req.Stream) {
		return fmt.Errorf("stream not found")
	}

	fmt.Println("Starting log stream for", req.Stream)

	for {
		message := ""
		switch req.Stream {
		case StreamNginx:
			message = `{ "body_bytes_sent": "37", "http_forward": "10.0.39.194", "http_header": "-", "http_referrer": "-", "http_user_agent": "ELB-HealthChecker/2.0", "http_x_amzn_trace_id": "-", "remote_addr": "10.0.39.194", "remote_user": "-", "request": "GET /readyz HTTP/1.1", "request_id": "6a06b1bf387b54f5f88cd1cac8c75de1", "request_method": "GET", "request_time": "0.001", "request_uri": "/readyz", "request_uri_query": "-", "server_name": "", "status": "200", "timestamp": "2025-03-25T01:21:50+00:00", "upstream_addr": "127.0.0.1:9000", "upstream_cache_status": "-", "upstream_http_x_drupal_cache": "-", "upstream_http_x_drupal_dynamic_cache": "-", "upstream_response_length": "23", "upstream_response_time": "0.002", "upstream_status": "200" }`
		case StreamFPM:
			message = `{ "body_bytes_sent": "0", "client_ip": "-", "cpu": "0.00", "headers": { "Cache-Control": "no-cache, no-store, must-revalidate, max-age=0" }, "http_referrer": "", "http_user_agent": "kube-probe/1.31+", "memory": "2097152", "remote_addr": "127.0.0.1", "remote_user": "", "request_id": "eaff1b73b352356ea0a321a250c1d591", "request_time": "0.001", "request_uri": "/readyz", "skpr_component": "fpm", "status": "200", "timestamp": "2025-03-25T01:21:26+0000" }`
		}
		pbMessage := &pb.LogTailResponse{
			Message: message,
		}
		sendError := server.Send(pbMessage)
		if sendError != nil {
			fmt.Println("Stopping log stream for", req.Stream)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
