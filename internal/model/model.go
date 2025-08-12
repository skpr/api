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

	backups := []*Backup{
		{
			Id:        "b84dd996-8113-4cd3-8dfe-018c990f5f1a",
			StartTime: time.Now().Add(-12 * time.Hour).Round(time.Second),
			Duration:  115 * time.Second,
		},
		{
			Id:        "aabb7bee-ab04-4ae1-ba3f-1142aae1353f",
			StartTime: time.Now().Add(-24 * time.Hour).Round(time.Second),
			Duration:  134 * time.Second,
		},
		{
			Id:        "7d419959-a3b7-483b-8fa7-769a0977e46b",
			StartTime: time.Now().Add(-36 * time.Hour).Round(time.Second),
			Duration:  123 * time.Second,
		},
	}

	config := map[string]*pb.Config{
		"mysql.default.database": {
			Key:   "mysql.default.database",
			Value: "heuu87bueoa8euua",
			Type:  pb.ConfigType_System,
		},
		"mysql.default.password": {
			Key:    "mysql.default.password",
			Value:  "Passw0rd",
			Secret: true,
			Type:   pb.ConfigType_System,
		},
		"my.personal.key": {
			Key:   "my.personal.key",
			Value: "A_value",
			Type:  pb.ConfigType_User,
		},
	}

	s.storage[name] = &Environment{
		Environment: environment,
		Config:      config,
		Cron:        cron,
		Purge:       purge,
		Backup:      backups,
	}
}

func (s *Model) DeleteEnvironment(name string) {
	delete(s.storage, name)
}

type Environment struct {
	Environment *pb.Environment
	Config      map[string]*pb.Config
	Cron        map[string]*Cron
	Purge       []*Purge
	Backup      []*Backup
}

type Cron struct {
	Suspended bool
}

<<<<<<< HEAD
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

func (m *Environment) AppendBackup(backup *Backup) {
	m.Backup = append(m.Backup, backup)
}

func (m *Environment) AddConfig(config *pb.Config) {
	m.Config[config.Key] = config
}

func (m *Environment) GetConfig(key string) (*pb.Config, error) {
	value, exists := m.Config[key]
	if !exists {
		return nil, fmt.Errorf("config key does not exist")
	}

	return value, nil
}

func (m *Environment) DeleteConfig(key string) error {
	_, exists := m.Config[key]
	if !exists {
		return fmt.Errorf("config key does not exist")
	}

	if m.Config[key].Type != pb.ConfigType_System {
		delete(m.Config, key)
	}
	return nil
}
