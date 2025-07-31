package cron

import (
	"context"
	"time"

	"github.com/skpr/api/internal/model"
	"github.com/skpr/api/pb"
)

// Server implements the GRPC "cron" definition.
type Server struct {
	Model *model.Model
	pb.UnimplementedCronServer
}

// Suspend all the cron jobs
func (c *Server) Suspend(ctx context.Context, req *pb.CronSuspendRequest) (*pb.CronSuspendResponse, error) {
	environment, err := c.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

	jobs := environment.CronDetail
	for key := range jobs {
		jobs[key].Suspended = true
	}

	return nil, nil
}

// Resume all the cron jobs
func (c *Server) Resume(ctx context.Context, req *pb.CronResumeRequest) (*pb.CronResumeResponse, error) {
	environment, err := c.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

	jobs := environment.CronDetail
	for key := range jobs {
		jobs[key].Suspended = false
	}

	return nil, nil
}

// List of mocked cron jobs.
func (c *Server) List(ctx context.Context, req *pb.CronListRequest) (*pb.CronListResponse, error) {
	environment, err := c.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

	resp := &pb.CronListResponse{}

	jobs := environment.CronDetail
	for _, value := range jobs {
		resp.List = append(resp.List, value)
	}

	return resp, nil
}

// JobList about when cron jobs last ran.
func (c *Server) JobList(ctx context.Context, req *pb.CronJobListRequest) (*pb.CronJobListResponse, error) {
	_, err := c.Model.GetEnvironment(req.Environment)
	if err != nil {
		return nil, err
	}

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
