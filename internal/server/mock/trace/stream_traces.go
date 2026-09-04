package trace

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/skpr/api/pb"
)

// StreamTraces streams traces from a specific environment.
func (s *Server) StreamTraces(req *pb.StreamTracesRequest, server pb.Trace_StreamTracesServer) error {
	for {
		// Tracing has been suspended for this environment, so there is nothing to
		// send. Fail fast instead of holding open a stream which sends nothing,
		// which the client cannot tell apart from an environment with no traffic.
		if s.IsSuspended(req.Environment) {
			return status.Errorf(codes.FailedPrecondition, "tracing is suspended for environment: %s", req.Environment)
		}

		now := time.Now()

		// ~6ms realistic function execution time
		latency := 6 * time.Millisecond

		// Function calls generator, placed relative to the start of the request.
		makeFunctionCalls := func() []*pb.TraceFunctionCall {
			return []*pb.TraceFunctionCall{
				{
					Name:    "PDOStatement::execute",
					Offset:  durationpb.New(0),
					Elapsed: durationpb.New(latency),
					Memory:  1048576,
				},
				{
					Name:    "Drupal\\Core\\Database\\StatementPrefetchIterator::execute",
					Offset:  durationpb.New(500 * time.Microsecond),
					Elapsed: durationpb.New(latency),
					Memory:  2097152,
				},
				{
					Name:    "Drupal\\sqlite\\Driver\\Database\\sqlite\\Statement::execute",
					Offset:  durationpb.New(1 * time.Millisecond),
					Elapsed: durationpb.New(latency),
					Memory:  524288,
				},
				{
					Name:    "Drupal\\Core\\Database\\Query\\Upsert::execute",
					Offset:  durationpb.New(1500 * time.Microsecond),
					Elapsed: durationpb.New(latency),
					Memory:  786432,
				},
			}
		}

		makeDrupal := func() *pb.TraceDrupal {
			return &pb.TraceDrupal{
				CacheEvents: []*pb.TraceDrupalCacheEvent{
					{
						Origin:   pb.TraceDrupalCacheOrigin_TRACE_DRUPAL_CACHE_ORIGIN_RENDER_ARRAY,
						Caller:   "Drupal\\Core\\Render\\Renderer::doRender",
						MaxAge:   -1,
						Tags:     []string{"node:1", "node_list"},
						Contexts: []string{"url.path", "user.permissions"},
						Offset:   durationpb.New(2 * time.Millisecond),
						Calls:    12,
					},
					{
						Origin:     pb.TraceDrupalCacheOrigin_TRACE_DRUPAL_CACHE_ORIGIN_OBJECT,
						Caller:     "Drupal\\Core\\Cache\\CacheableMetadata::createFromObject",
						ObjectType: "Drupal\\node\\Entity\\Node",
						MaxAge:     3600,
						Tags:       []string{"node:1"},
						Contexts:   []string{"user.roles"},
						Offset:     durationpb.New(4 * time.Millisecond),
						Calls:      3,
					},
				},
				CacheEventsDropped: 0,
			}
		}

		// Create 3 traces with increasing offsets
		traces := make([]*pb.Trace, 0, 3)

		for i := 0; i < 3; i++ {
			start := now.Add(time.Duration(i*250) * time.Millisecond)
			end := start.Add(50 * time.Millisecond)

			traces = append(traces, &pb.Trace{
				Metadata: &pb.TraceMetadata{
					RequestId: gofakeit.UUID(),
					StartTime: timestamppb.New(start),
					EndTime:   timestamppb.New(end),
					Source:    pb.TraceSource_TRACE_SOURCE_HTTP,
					Runtime:   pb.TraceRuntime_TRACE_RUNTIME_PHP,
					Http: &pb.TraceMetadataHTTP{
						Method: http.MethodGet,
						Uri:    "/sites/default/files/styles/scale_crop_7_3_wide/public/veggie-pasta-bake-hero-umami.jpg.webp?itok=CYsHBUlX",
					},
				},
				FunctionCalls:        makeFunctionCalls(),
				FunctionCallsDropped: 0,
				ResourceUtilisation: &pb.TraceResourceUtilisation{
					MaxMemory: 33554432,
				},
				Drupal: makeDrupal(),
			})
		}

		resp := &pb.StreamTracesResponse{Traces: traces}

		if err := server.Send(resp); err != nil {
			return fmt.Errorf("stopping log stream for: %w", err)
		}

		time.Sleep(time.Second)
	}
}
