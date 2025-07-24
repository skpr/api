package compass

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/skpr/api/pb"
)

// StreamTraces streams Compass traces from a specific environment.
func (s *Server) StreamTraces(ctx context.Context, req *pb.StreamTracesRequest) (*pb.StreamTracesResponse, error) {
	for {
		// Simulate some processing delay.
		time.Sleep(time.Second)

		resp := &pb.StreamTracesResponse{
			Traces: []*pb.CompassTrace{
				{
					Metadata: &pb.CompassTraceMetadata{
						RequestId: gofakeit.UUID(),
						Uri:       "/sites/default/files/styles/scale_crop_7_3_wide/public/veggie-pasta-bake-hero-umami.jpg.webp?itok=CYsHBUlX",
						StartTime: 11479712402527,
						EndTime:   11480550685871,
					},
					FunctionCalls: []*pb.CompassTraceFunctionCall{
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
					Metadata: &pb.CompassTraceMetadata{
						RequestId: gofakeit.UUID(),
						Uri:       "/sites/default/files/styles/scale_crop_7_3_wide/public/veggie-pasta-bake-hero-umami.jpg.webp?itok=CYsHBUlX",
						StartTime: 11479712402527,
						EndTime:   11480550685871,
					},
					FunctionCalls: []*pb.CompassTraceFunctionCall{
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
					Metadata: &pb.CompassTraceMetadata{
						RequestId: gofakeit.UUID(),
						Uri:       "/sites/default/files/styles/scale_crop_7_3_wide/public/veggie-pasta-bake-hero-umami.jpg.webp?itok=CYsHBUlX",
						StartTime: 11479712402527,
						EndTime:   11480550685871,
					},
					FunctionCalls: []*pb.CompassTraceFunctionCall{
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

		// Simulate sending the response back to the client.
		return resp, nil
	}
}
