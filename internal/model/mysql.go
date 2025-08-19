package model

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/skpr/api/pb"
)

type Mysql struct {
	Id        string
	StartTime time.Time
	Duration  time.Duration
}

func NewMysql(environment string) *Mysql {
	return &Mysql{
		Id:        environment + "-" + gofakeit.UUID(),
		StartTime: time.Now().Round(time.Second),
		Duration:  150 * time.Second,
	}
}

func (b *Mysql) Status() pb.ImageStatus_Phase {
	status := pb.ImageStatus_Completed
	if b.StartTime.Add(b.Duration).After(time.Now()) {
		status = pb.ImageStatus_InProgress
	}
	return status
}
