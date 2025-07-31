package model

import (
	"fmt"
	"time"

	"github.com/skpr/api/pb"
)

type Model struct {
	storage map[string]*EnvironmentModel
}

func NewModel() *Model {
	return &Model{
		storage: make(map[string]*EnvironmentModel),
	}
}

func (s *Model) GetEnvironment(name string) (*EnvironmentModel, error) {
	if name == "" {
		return nil, fmt.Errorf("environment not provided")
	}

	value, exists := s.storage[name]
	if !exists {
		return nil, fmt.Errorf("environment not valid")
	}

	return value, nil
}

func (s *Model) GetEnvironments() []*EnvironmentModel {
	var response []*EnvironmentModel
	for _, value := range s.storage {
		response = append(response, value)
	}
	return response
}

func (s *Model) CreateEnvironment(name string, size int) {
	environment := &pb.Environment{
		Name:       name,
		Version:    "v0.0.1",
		Phase:      "Deployed",
		Production: name == "prod",
		Ingress: &pb.Ingress{
			Routes: []string{
				fmt.Sprintf("%s.mock.local.skpr.dev", name),
			},
		},
		Resources: &pb.EnvironmentResources{
			CPU: &pb.EnvironmentResourcesCPU{
				Current: 10 * int64(size),
				Limit:   1000 * int64(size),
			},
			Memory: &pb.EnvironmentResourcesMemory{
				Current: 128 * int64(size),
				Limit:   512 * int64(size),
			},
			Replicas: &pb.EnvironmentResourcesReplicas{
				Current: 1 * int32(size),
				Min:     1 * int32(size),
				Max:     5 * int32(size),
			},
		},
		Cron: []*pb.Cron{
			{
				Name:     "drush",
				Command:  "drush cron",
				Schedule: "* * * * *",
			},
			{
				Name:     "search-api-index",
				Command:  "drush search-api:index example",
				Schedule: "*/6 * * * *",
			},
			{
				Name:     "queue-run",
				Command:  "drush queue:run example",
				Schedule: "* * * * *",
			},
		},
	}
	if name == "prod" {
		environment.Ingress.Routes = append(environment.Ingress.Routes, "example.com", "www.example.com")
	}

	// Duplicated for now but considering a refactor for later.
	cron := map[string]*pb.CronDetail{
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

	s.storage[name] = &EnvironmentModel{
		Environment: environment,
		CronDetail:  cron,
	}
}

func (s *Model) DeleteEnvironment(name string) {
	delete(s.storage, name)
}

type EnvironmentModel struct {
	Environment *pb.Environment
	CronDetail  map[string]*pb.CronDetail
}
