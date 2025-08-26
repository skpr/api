package model

type Project struct {
	Id      string
	Name    string
	Tags    []string
	Contact string
	Size    string
}

func NewProject(id string, name string, tags []string, size string) *Project {
	return &Project{
		Id:      id,
		Name:    name,
		Tags:    tags,
		Contact: "admin@example.com",
		Size:    size,
	}
}
