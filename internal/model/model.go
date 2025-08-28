package model

import (
	"fmt"
	"time"

	"github.com/skpr/api/pb"
)

type Model struct {
	Environments map[string]*Environment
	Projects     map[string]*Project
}

func NewModel() *Model {
	return &Model{
		Environments: make(map[string]*Environment),
		Projects:     make(map[string]*Project),
	}
}

func (s *Model) GetEnvironment(name string) (*Environment, error) {
	if name == "" {
		return nil, fmt.Errorf("environment not provided")
	}

	value, exists := s.Environments[name]
	if !exists {
		return nil, fmt.Errorf("environment not valid")
	}

	return value, nil
}

func (s *Model) GetEnvironments() []*Environment {
	var response []*Environment
	for _, value := range s.Environments {
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
				Current: 100000 * int64(size),
				Limit:   1000000 * int64(size),
			},
			Memory: &pb.EnvironmentResourcesMemory{
				// In GB.
				Current: 128 * int64(size) * 1024,
				Limit:   512 * int64(size) * 1024,
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
			{
				Name:     "data-pipeline-refresh",
				Command:  "drush data-pipelines:reindex my_source_data_from_json_url && drush dn:invalidate-paths /path/to/page/that/uses/data",
				Schedule: "*/30 * * * *",
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

	backups := map[string]*Backup{
		name + "-b84dd996-8113-4cd3-8dfe-018c990f5f1a": {
			Id:        name + "-b84dd996-8113-4cd3-8dfe-018c990f5f1a",
			StartTime: time.Now().Add(-12 * time.Hour).Round(time.Second),
			Duration:  115 * time.Second,
			Failed:    false,
		},
		name + "-aabb7bee-ab04-4ae1-ba3f-1142aae1353f": {
			Id:        name + "-aabb7bee-ab04-4ae1-ba3f-1142aae1353f",
			StartTime: time.Now().Add(-24 * time.Hour).Round(time.Second),
			Duration:  134 * time.Second,
			Failed:    true,
		},
		name + "-7d419959-a3b7-483b-8fa7-769a0977e46b": {
			Id:        name + "-7d419959-a3b7-483b-8fa7-769a0977e46b",
			StartTime: time.Now().Add(-36 * time.Hour).Round(time.Second),
			Duration:  123 * time.Second,
			Failed:    false,
		},
	}

	restores := map[string]*Restore{
		name + "-ef394df8-fbf6-4a3d-9b4b-6c1011a8afa5": {
			Id:        name + "-ef394df8-fbf6-4a3d-9b4b-6c1011a8afa5",
			BackupId:  name + "-7d419959-a3b7-483b-8fa7-769a0977e46b",
			StartTime: time.Now().Add(-12 * time.Hour).Round(time.Second),
			Duration:  145 * time.Second,
			Failed:    false,
		},
		name + "-39b11746-eca6-40ac-9aed-094d38c1e444": {
			Id:        name + "-39b11746-eca6-40ac-9aed-094d38c1e444",
			BackupId:  name + "-7d419959-a3b7-483b-8fa7-769a0977e46b",
			StartTime: time.Now().Add(-12 * time.Hour).Round(time.Second),
			Duration:  145 * time.Second,
			Failed:    true,
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

	mysql := map[string]*Mysql{
		name + "-7GH59AB92AJF": {
			Id:        name + "-7GH59AB92AJF",
			StartTime: time.Now().Add(-12 * time.Hour).Round(time.Second),
			Duration:  115 * time.Second,
		},
		name + "-GR87QH51ALQ5": {
			Id:        name + "-GR87QH51ALQ5",
			StartTime: time.Now().Add(-24 * time.Hour).Round(time.Second),
			Duration:  134 * time.Second,
		},
	}

	s.Environments[name] = &Environment{
		Environment: environment,
		Config:      config,
		Cron:        cron,
		Purge:       purge,
		Backup:      backups,
		Restore:     restores,
		Mysql:       mysql,
	}
}

func (s *Model) DeleteEnvironment(name string) {
	delete(s.Environments, name)
}

func (m *Model) GetBackup(id string) (*Backup, error) {
	for _, value := range m.GetEnvironments() {
		backup, exists := value.Backup[id]
		if exists {
			return backup, nil
		}
	}
	return nil, fmt.Errorf("backup with id not found")
}

func (m *Model) GetRestore(id string) (*Restore, error) {
	for _, value := range m.GetEnvironments() {
		restore, exists := value.Restore[id]
		if exists {
			return restore, nil
		}
	}
	return nil, fmt.Errorf("restore with id not found")
}

type Environment struct {
	Environment *pb.Environment
	Config      map[string]*pb.Config
	Cron        map[string]*Cron
	Purge       []*Purge
	Backup      map[string]*Backup
	Restore     map[string]*Restore
	Mysql       map[string]*Mysql
}

type Cron struct {
	Suspended bool
}

func (m *Environment) AddPurge(purge *Purge) {
	m.Purge = append(m.Purge, purge)
}

func (m *Environment) AddBackup(backup *Backup) {
	m.Backup[backup.Id] = backup
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

	if m.Config[key].Type == pb.ConfigType_Overridden {
		m.Config[key].Type = pb.ConfigType_System
		m.Config[key].Value = configDefaults[key]
	} else if m.Config[key].Type != pb.ConfigType_System {
		delete(m.Config, key)
	}
	return nil
}

func (m *Environment) AddRestore(restore *Restore) {
	m.Restore[restore.Id] = restore
}

func (m *Model) AddProject(project *Project) {
	m.Projects[project.Id] = project
}

func (s *Model) GetProjects() []*Project {
	var response []*Project
	for _, value := range s.Projects {
		response = append(response, value)
	}
	return response
}

func (s *Model) GetProject(name string) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("project not provided")
	}

	value, exists := s.Projects[name]
	if !exists {
		return nil, fmt.Errorf("project not valid")
	}

	return value, nil
}
