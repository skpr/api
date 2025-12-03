package trace

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/skpr/api/pb"
)

// StreamTraces streams traces from a specific environment.
func (s *Server) StreamTraces(_ *pb.StreamTracesRequest, server pb.Trace_StreamTracesServer) error {
	for {
		now := time.Now()

		// ~6ms realistic function execution time
		latency := 6 * time.Millisecond

		// Function calls generator with Duration + timestamp
		makeFunctionCalls := func(base time.Time) []*pb.TraceFunctionCall {
			return []*pb.TraceFunctionCall{
				{
					Name:      "PDOStatement::execute",
					StartTime: timestamppb.New(base),
					Elapsed:   durationpb.New(latency),
				},
				{
					Name:      "Drupal\\Core\\Database\\StatementPrefetchIterator::execute",
					StartTime: timestamppb.New(base.Add(500 * time.Microsecond)),
					Elapsed:   durationpb.New(latency),
				},
				{
					Name:      "Drupal\\sqlite\\Driver\\Database\\sqlite\\Statement::execute",
					StartTime: timestamppb.New(base.Add(1 * time.Millisecond)),
					Elapsed:   durationpb.New(latency),
				},
				{
					Name:      "Drupal\\Core\\Database\\Query\\Upsert::execute",
					StartTime: timestamppb.New(base.Add(1500 * time.Microsecond)),
					Elapsed:   durationpb.New(latency),
				},
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
					Method:    http.MethodGet,
					Uri:       "/sites/default/files/styles/scale_crop_7_3_wide/public/veggie-pasta-bake-hero-umami.jpg.webp?itok=CYsHBUlX",
					StartTime: timestamppb.New(start),
					EndTime:   timestamppb.New(end),
				},
				FunctionCalls: makeFunctionCalls(start),
			})
		}

		resp := &pb.StreamTracesResponse{Traces: traces}

		if err := server.Send(resp); err != nil {
			return fmt.Errorf("stopping log stream for: %w", err)
		}

		time.Sleep(time.Second)
	}
}
