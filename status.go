package main

import (
	"log"
	"time"
)

type Status struct {
	GitCommit string
	Counter int
	LatestAdditions []string
}

type StatusCoordinator struct {
	status *Status
	recordCollectionAdditions RecordCollectionAdditions
}

type RecordCollectionAdditions interface {
	FetchLatestAdditions() []string
}

func NewStatusCoordinator(gitCommit string, fetchInterval int, recordCollectionAdditions RecordCollectionAdditions) *StatusCoordinator {
	status := &Status{
		GitCommit: gitCommit,
		Counter: 0,
		LatestAdditions: []string{},
	}
	statusCoordinator := &StatusCoordinator{
		status,
		recordCollectionAdditions,
	}

	go func() {
		c := time.Tick(time.Duration(fetchInterval) * time.Millisecond)
		statusCoordinator.doWork()
		for _ = range c {
			statusCoordinator.doWork()
		}
	}()

	return statusCoordinator
}

func (s *StatusCoordinator) doWork() {
	s.status.LatestAdditions = s.recordCollectionAdditions.FetchLatestAdditions()
	s.status.Counter = s.status.Counter + 1
	log.Printf("Status counter:%v, additions:%v", s.status.Counter, s.status.LatestAdditions)
}