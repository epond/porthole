package main

import (
	"log"
	"time"
)

type StatusCoordinator struct {
	status *Status
	recordCollectionAdditions RecordCollectionAdditions
}

type RecordCollectionAdditions interface {
	latestAdditions() []string
}

func NewStatusCoordinator(status *Status, fetchInterval int, recordCollectionAdditions RecordCollectionAdditions) *StatusCoordinator {
	statusCoordinator := &StatusCoordinator{status, recordCollectionAdditions}

	go func() {
		c := time.Tick(time.Duration(fetchInterval) * time.Second)
		statusCoordinator.doWork()
		for _ = range c {
			statusCoordinator.doWork()
		}
	}()

	return statusCoordinator
}

func (s *StatusCoordinator) doWork() {
	s.status.LatestAdditions = s.recordCollectionAdditions.latestAdditions()
	s.status.Counter = s.status.Counter + 1
	log.Printf("Status counter:%v, additions:%v", s.status.Counter, s.status.LatestAdditions)
}