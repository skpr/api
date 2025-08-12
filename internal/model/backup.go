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
}

func NewBackup() *Backup {
	return &Backup{
		Id:        gofakeit.UUID(),
		StartTime: time.Now().Round(time.Second),
		Duration:  120 * time.Second,
	}
}

func (b *Backup) Status() pb.BackupStatus_Phase {
	status := pb.BackupStatus_Completed
	if b.StartTime.Add(b.Duration).After(time.Now()) {
		status = pb.BackupStatus_InProgress
	}
	return status
}
