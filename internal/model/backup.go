package model

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/skpr/api/pb"
)

type Backup struct {
	Id        string
	StartTime time.Time
	Duration  time.Duration
	Failed    bool
}

func NewBackup(environment string) *Backup {
	return &Backup{
		Id:        environment + "-" + gofakeit.UUID(),
		StartTime: time.Now().Round(time.Second),
		Duration:  120 * time.Second,
		Failed:    false,
	}
}

func (b *Backup) Status() pb.BackupStatus_Phase {
	status := pb.BackupStatus_Completed
	if b.Failed {
		status = pb.BackupStatus_Failed
	} else if b.StartTime.Add(b.Duration).After(time.Now()) {
		status = pb.BackupStatus_InProgress
	}
	return status
}
