package main

import (
	"log"
	"time"
)

type Status struct {
	GitCommit string
	Counter int
	LastRequest time.Time
	LatestAdditions []string
}

type StatusCoordinator struct {
	status *Status
	recordCollectionAdditions RecordCollectionAdditions
}

type RecordCollectionAdditions interface {
	FetchLatestAdditions() []string
}

func NewStatusCoordinator(gitCommit string, recordCollectionAdditions RecordCollectionAdditions, clock <-chan time.Time) *StatusCoordinator {
	status := &Status{
		GitCommit: gitCommit,
		Counter: 0,
		LastRequest: time.Now(),
		LatestAdditions: []string{},
	}
	statusCoordinator := &StatusCoordinator{
		status,
		recordCollectionAdditions,
	}

	go func() {
		statusCoordinator.doWork()
		for _ = range clock {
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