package logs

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/skpr/api/pb"
)

const (
	StreamNginx = "nginx"
	StreamFPM   = "fpm"
)

// Raw JSON payloads reused from the Tail mock, kept here so fixture events
// look realistic without duplicating the strings across methods.
const (
	rawNginx = `{ "body_bytes_sent": "37", "http_forward": "10.0.39.194", "http_header": "-", "http_referrer": "-", "http_user_agent": "ELB-HealthChecker/2.0", "http_x_amzn_trace_id": "-", "remote_addr": "10.0.39.194", "remote_user": "-", "request": "GET /readyz HTTP/1.1", "request_id": "6a06b1bf387b54f5f88cd1cac8c75de1", "request_method": "GET", "request_time": "0.001", "request_uri": "/readyz", "request_uri_query": "-", "server_name": "", "status": "200", "timestamp": "2025-03-25T01:21:50+00:00", "upstream_addr": "127.0.0.1:9000", "upstream_cache_status": "-", "upstream_http_x_drupal_cache": "-", "upstream_http_x_drupal_dynamic_cache": "-", "upstream_response_length": "23", "upstream_response_time": "0.002", "upstream_status": "200" }`
	rawFPM   = `{ "body_bytes_sent": "0", "client_ip": "-", "cpu": "0.00", "headers": { "Cache-Control": "no-cache, no-store, must-revalidate, max-age=0" }, "http_referrer": "", "http_user_agent": "kube-probe/1.31+", "memory": "2097152", "remote_addr": "127.0.0.1", "remote_user": "", "request_id": "eaff1b73b352356ea0a321a250c1d591", "request_time": "0.001", "request_uri": "/readyz", "skpr_component": "fpm", "status": "200", "timestamp": "2025-03-25T01:21:26+0000" }`
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
			message = rawNginx
		case StreamFPM:
			message = rawFPM
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

// buildMockEvents returns a fresh fixture with timestamps relative to now.
// Built per request so events always appear recent to the caller.
func buildMockEvents() []*pb.LogEvent {
	now := time.Now()
	return []*pb.LogEvent{
		{
			Timestamp: timestamppb.New(now.Add(-1 * time.Minute)),
			Stream:    StreamNginx,
			Message:   rawNginx,
		},
		{
			Timestamp: timestamppb.New(now.Add(-3 * time.Minute)),
			Stream:    StreamFPM,
			Message:   rawFPM,
		},
		{
			Timestamp: timestamppb.New(now.Add(-7 * time.Minute)),
			Stream:    StreamNginx,
			Message:   rawNginx,
		},
		{
			Timestamp: timestamppb.New(now.Add(-12 * time.Minute)),
			Stream:    StreamFPM,
			Message:   rawFPM,
		},
		{
			Timestamp: timestamppb.New(now.Add(-25 * time.Minute)),
			Stream:    StreamNginx,
			Message:   rawNginx,
		},
		{
			Timestamp: timestamppb.New(now.Add(-45 * time.Minute)),
			Stream:    StreamFPM,
			Message:   rawFPM,
		},
	}
}

// Query streams matching log events followed by a terminal metadata message.
// The Window oneof (Timeframe or TimeRange) on the filter is accepted but ignored by the mock.
func (s *Server) Query(req *pb.LogQueryRequest, stream pb.Logs_QueryServer) error {
	if req.Filter == nil {
		return fmt.Errorf("filter not provided")
	}

	if req.Filter.Environment == "" {
		return fmt.Errorf("environment not provided")
	}

	events := buildMockEvents()
	scanned := int64(len(events))

	// Filter by stream selection if provided.
	if len(req.Filter.Streams) > 0 {
		filtered := events[:0:0]
		for _, evt := range events {
			if slices.Contains(req.Filter.Streams, evt.Stream) {
				filtered = append(filtered, evt)
			}
		}
		events = filtered
	}

	// Filter by substring queries if provided. Each entry must be satisfied:
	// include entries require the substring to be present, exclude entries
	// require it to be absent. Empty Value entries are skipped.
	if len(req.Filter.Contains) > 0 {
		filtered := events[:0:0]
		for _, evt := range events {
			ok := true
			for _, f := range req.Filter.Contains {
				if f == nil || f.Value == "" {
					continue
				}
				has := strings.Contains(evt.Message, f.Value)
				if f.Exclude && has {
					ok = false
					break
				}
				if !f.Exclude && !has {
					ok = false
					break
				}
			}
			if ok {
				filtered = append(filtered, evt)
			}
		}
		events = filtered
	}

	// Apply result cap.
	if req.Limit > 0 && int(req.Limit) < len(events) {
		events = events[:req.Limit]
	}

	for _, evt := range events {
		resp := &pb.LogQueryResponse{
			Body: &pb.LogQueryResponse_Event{Event: evt},
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}

	meta := &pb.LogQueryResponse{
		Body: &pb.LogQueryResponse_Meta{
			Meta: &pb.LogQueryMeta{
				Scanned: scanned,
				RanAt:   timestamppb.Now(),
			},
		},
	}
	return stream.Send(meta)
}

// Summarise returns a canned AI-style summary of the requested log window.
// The Prompt field is ignored by the mock.
func (s *Server) Summarise(ctx context.Context, req *pb.LogSummariseRequest) (*pb.LogSummariseResponse, error) {
	if req.Filter == nil {
		return nil, fmt.Errorf("filter not provided")
	}

	if req.Filter.Environment == "" {
		return nil, fmt.Errorf("environment not provided")
	}

	return &pb.LogSummariseResponse{
		Overview: "Traffic is nominal with a small number of 5xx errors originating from the fpm stream; nginx is healthy.",
		Bullets: []string{
			"99.2% of requests returned 2xx",
			"3 elevated-error windows detected in fpm between 01:20-01:40 UTC",
			"No WAF blocks observed in the window",
		},
		SuggestedActions: []string{
			"Inspect fpm error logs around 01:30 UTC for root cause",
			"Consider increasing the fpm worker count if load continues to rise",
		},
	}, nil
}
