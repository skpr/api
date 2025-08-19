package model

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/skpr/api/pb"
)

type Restore struct {
	Id        string
	BackupId  string
	StartTime time.Time
	Duration  time.Duration
}

func NewRestore(environment string, backupId string) *Restore {
	return &Restore{
		Id:        environment + "-" + gofakeit.UUID(),
		BackupId:  backupId,
		StartTime: time.Now().Round(time.Second),
		Duration:  150 * time.Second,
	}
}

func (b *Restore) Status() pb.RestoreStatus_Phase {
	status := pb.RestoreStatus_Completed
	if b.StartTime.Add(b.Duration).After(time.Now()) {
		status = pb.RestoreStatus_InProgress
	}
	return status
}
