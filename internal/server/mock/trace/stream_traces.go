package trace

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/skpr/api/pb"
)

// StreamTraces streams traces from a specific environment.
func (s *Server) StreamTraces(_ *pb.StreamTracesRequest, server pb.Trace_StreamTracesServer) error {
	for {
		// Simulate some processing delay.
		time.Sleep(time.Second)

		for {
			resp := &pb.StreamTracesResponse{
				Traces: []*pb.Trace{
					{
						Metadata: &pb.TraceMetadata{
							RequestId: gofakeit.UUID(),
							Method:    http.MethodGet,
							Uri:       "/sites/default/files/styles/scale_crop_7_3_wide/public/veggie-pasta-bake-hero-umami.jpg.webp?itok=CYsHBUlX",
							StartTime: 11479712402527,
							EndTime:   11480550685871,
						},
						FunctionCalls: []*pb.TraceFunctionCall{
							{
								Name:        "PDOStatement::execute",
								StartTime:   11479719656578,
								ElapsedTime: 5999966,
							},
							{
								Name:        "Drupal\\Core\\Database\\StatementPrefetchIterator::execute",
								StartTime:   11479719664878,
								ElapsedTime: 5999966,
							},
							{
								Name:        "Drupal\\sqlite\\Driver\\Database\\sqlite\\Statement::execute",
								StartTime:   11479719666498,
								ElapsedTime: 5999966,
							},
							{
								Name:        "Drupal\\Core\\Database\\Query\\Upsert::execute",
								StartTime:   11479719668488,
								ElapsedTime: 5999966,
							},
						},
					},
					{
						Metadata: &pb.TraceMetadata{
							RequestId: gofakeit.UUID(),
							Method:    http.MethodGet,
							Uri:       "/sites/default/files/styles/scale_crop_7_3_wide/public/veggie-pasta-bake-hero-umami.jpg.webp?itok=CYsHBUlX",
							StartTime: 11479712402527,
							EndTime:   11480550685871,
						},
						FunctionCalls: []*pb.TraceFunctionCall{
							{
								Name:        "PDOStatement::execute",
								StartTime:   11479719656578,
								ElapsedTime: 5999966,
							},
							{
								Name:        "Drupal\\Core\\Database\\StatementPrefetchIterator::execute",
								StartTime:   11479719664878,
								ElapsedTime: 5999966,
							},
							{
								Name:        "Drupal\\sqlite\\Driver\\Database\\sqlite\\Statement::execute",
								StartTime:   11479719666498,
								ElapsedTime: 5999966,
							},
							{
								Name:        "Drupal\\Core\\Database\\Query\\Upsert::execute",
								StartTime:   11479719668488,
								ElapsedTime: 5999966,
							},
						},
					},
					{
						Metadata: &pb.TraceMetadata{
							RequestId: gofakeit.UUID(),
							Uri:       "/sites/default/files/styles/scale_crop_7_3_wide/public/veggie-pasta-bake-hero-umami.jpg.webp?itok=CYsHBUlX",
							StartTime: 11479712402527,
							EndTime:   11480550685871,
						},
						FunctionCalls: []*pb.TraceFunctionCall{
							{
								Name:        "PDOStatement::execute",
								StartTime:   11479719656578,
								ElapsedTime: 5999966,
							},
							{
								Name:        "Drupal\\Core\\Database\\StatementPrefetchIterator::execute",
								StartTime:   11479719664878,
								ElapsedTime: 5999966,
							},
							{
								Name:        "Drupal\\sqlite\\Driver\\Database\\sqlite\\Statement::execute",
								StartTime:   11479719666498,
								ElapsedTime: 5999966,
							},
							{
								Name:        "Drupal\\Core\\Database\\Query\\Upsert::execute",
								StartTime:   11479719668488,
								ElapsedTime: 5999966,
							},
						},
					},
				},
			}

			err := server.Send(resp)
			if err != nil {
				return fmt.Errorf("stopping log stream for: %w", err)
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}
