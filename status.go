package main

import (
	"log"
	"time"
)

type Status struct {
	GitCommit       string
	Counter         int
	LastRequest     time.Time
	LatestAdditions []string
}

type StatusCoordinator struct {
	status             *Status
	statusUpdateWorker StatusUpdateWorker
}

type StatusUpdateWorker interface {
	UpdateStatus(status *Status)
}

func NewStatusCoordinator(
	gitCommit string,
	statusUpdateWorker StatusUpdateWorker,
	clock <-chan time.Time) *StatusCoordinator {

	status := &Status{
		GitCommit:       gitCommit,
		Counter:         0,
		LastRequest:     time.Now(),
		LatestAdditions: []string{},
	}
	statusCoordinator := &StatusCoordinator{
		status,
		statusUpdateWorker,
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
	s.statusUpdateWorker.UpdateStatus(s.status)
	s.status.Counter = s.status.Counter + 1
	log.Printf("Status counter:%v, additions:%v", s.status.Counter, s.status.LatestAdditions)
}
