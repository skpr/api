package model

import (
	"fmt"
	"time"

	"github.com/skpr/api/internal/random"
	"github.com/skpr/api/pb"
)

type Model struct {
	storage map[string]*Environment
}

func NewModel() *Model {
	return &Model{
		storage: make(map[string]*Environment),
	}
}

func (s *Model) GetEnvironment(name string) (*Environment, error) {
	if name == "" {
		return nil, fmt.Errorf("environment not provided")
	}

	value, exists := s.storage[name]
	if !exists {
		return nil, fmt.Errorf("environment not valid")
	}

	return value, nil
}

func (s *Model) GetEnvironments() []*Environment {
	var response []*Environment
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
				Schedule: "0 * * * *",
			},
			{
				Name:     "search-api-index",
				Command:  "drush search-api:index example",
				Schedule: "*/5 * * * *",
			},
			{
				Name:     "queue-run",
				Command:  "drush queue:run example",
				Schedule: "*/15 * * * *",
			},
		},
	}
	if name == "prod" {
		environment.Ingress.Routes = append(environment.Ingress.Routes, "example.com", "www.example.com")
	}

	cron := make(map[string]*Cron)
	for _, value := range environment.Cron {
		cron[value.Name] = &Cron{
			Suspended: false,
		}
	}

	purge := []*Purge{
		{
			Id:      "YT7A8K9YQI5RNST2N0RY",
			Created: time.Now().Round(time.Second),
			Paths:   []string{"/example/path/1"},
		},
		{
			Id:      "JAA4IVVQCIZXV39BXVHT",
			Created: time.Now().Add(-2 * time.Hour).Round(time.Second),
			Paths:   []string{"/example/path"},
		},
		{
			Id:      "B3Y8OOXCY9D6UFSE1NN3",
			Created: time.Now().Add(-3 * time.Hour).Round(time.Second),
			Paths: []string{
				"/example/*",
				"/test/lower/path",
			},
		},
	}

	s.storage[name] = &Environment{
		Environment: environment,
		Cron:        cron,
		Purge:       purge,
	}
}

func (s *Model) DeleteEnvironment(name string) {
	delete(s.storage, name)
}

type Environment struct {
	Environment *pb.Environment
	Cron        map[string]*Cron
	Purge       []*Purge
}

type Cron struct {
	Suspended bool
}

func (m *Environment) AppendPurge(purge *Purge) {
	m.Purge = append(m.Purge, purge)
}

type Purge struct {
	Id      string
	Created time.Time
	Paths   []string
}

func NewPurge(paths []string) *Purge {
	return &Purge{
		Id:      random.StringOfLength(20),
		Created: time.Now().Round(time.Second),
		Paths:   paths,
	}
}
