package model

import (
	"time"

	"github.com/skpr/api/internal/random"
)

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
