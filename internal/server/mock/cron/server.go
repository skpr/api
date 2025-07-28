package cron

import (
	"context"
	"time"

	"github.com/skpr/api/pb"
)

var state = map[string]*pb.CronDetail{
	"drush": {
		Name:               "drush",
		Schedule:           "* * * * *",
		Command:            "drush cron",
		LastScheduleTime:   time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		LastSuccessfulTime: time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
	},
	"search-api-index": {
		Name:               "search-api-index",
		Schedule:           "*/6 * * * *",
		Command:            "drush search-api:index example",
		LastScheduleTime:   time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		LastSuccessfulTime: time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
	},
	"queue-run": {
		Name:               "queue-run",
		Schedule:           "* * * * *",
		Command:            "drush queue:run example",
		LastScheduleTime:   time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		LastSuccessfulTime: time.Now().Add(-2 * time.Minute).Format(time.RFC3339),
		Suspended:          true,
	},
}

// Server implements the GRPC "version" definition.
type Server struct {
	pb.UnimplementedCronServer
}

// Suspend all the cron jobs
func (c *Server) Suspend(ctx context.Context, req *pb.CronSuspendRequest) (*pb.CronSuspendResponse, error) {
	for key := range state {
		state[key].Suspended = true
	}

	return nil, nil
}

// Resume all the cron jobs
func (c *Server) Resume(ctx context.Context, req *pb.CronResumeRequest) (*pb.CronResumeResponse, error) {
	for key := range state {
		state[key].Suspended = false
	}

	return nil, nil
}

// List of mocked cron jobs.
func (c *Server) List(ctx context.Context, req *pb.CronListRequest) (*pb.CronListResponse, error) {
	var values = []*pb.CronDetail{}
	for _, value := range state {
		values = append(values, value)
	}
	resp := &pb.CronListResponse{
		List: values,
	}

	return resp, nil
}

// JobList about when cron jobs last ran.
func (c *Server) JobList(ctx context.Context, req *pb.CronJobListRequest) (*pb.CronJobListResponse, error) {
	resp := &pb.CronJobListResponse{
		List: []*pb.CronJobDetail{
			{
				Name:      "drush",
				Phase:     pb.CronJobDetail_Running,
				StartTime: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
				Duration:  (10 * time.Second).String(),
			},
			{
				Name:      "drush",
				Phase:     pb.CronJobDetail_Succeeded,
				StartTime: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
				Duration:  (10 * time.Second).String(),
			},

			{
				Name:      "drush",
				Phase:     pb.CronJobDetail_Succeeded,
				StartTime: time.Now().Add(-3 * time.Hour).Format(time.RFC3339),
				Duration:  (10 * time.Second).String(),
			},

			{
				Name:      "drush",
				Phase:     pb.CronJobDetail_Succeeded,
				StartTime: time.Now().Add(-4 * time.Hour).Format(time.RFC3339),
				Duration:  (10 * time.Second).String(),
			},

			{
				Name:      "drush",
				Phase:     pb.CronJobDetail_Succeeded,
				StartTime: time.Now().Add(-5 * time.Hour).Format(time.RFC3339),
				Duration:  (10 * time.Second).String(),
			},

			{
				Name:      "drush",
				Phase:     pb.CronJobDetail_Succeeded,
				StartTime: time.Now().Add(-6 * time.Hour).Format(time.RFC3339),
				Duration:  (10 * time.Second).String(),
			},
			{
				Name:      "search-api-index",
				Phase:     pb.CronJobDetail_Succeeded,
				StartTime: time.Now().Add(-7 * time.Hour).Format(time.RFC3339),
				Duration:  (5 * time.Minute).String(),
			},
			{
				Name:      "queue-run",
				Phase:     pb.CronJobDetail_Succeeded,
				StartTime: time.Now().Add(-8 * time.Hour).Format(time.RFC3339),
				Duration:  time.Minute.String(),
			},
		},
	}

	return resp, nil
}
