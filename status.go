package main

import (
	"log"
	"time"
)

type Status struct {
	GitCommit       string
	PreviousState   string
	LastRequest     time.Time
	LastFetch       string
	LatestAdditions []string
}

type StatusCoordinator struct {
	status             *Status
	statusUpdateWorker StatusUpdateWorker
	sleepAfter         time.Duration
}

type StatusUpdateWorker interface {
	UpdateStatus(timestamp time.Time, status *Status)
}

func NewStatusCoordinator(
	gitCommit string,
	statusUpdateWorker StatusUpdateWorker,
	clock <-chan time.Time,
	sleepAfter time.Duration) *StatusCoordinator {

	status := &Status{
		GitCommit:       gitCommit,
		PreviousState:   "",
		LastRequest:     time.Now(),
		LastFetch:       "",
		LatestAdditions: []string{},
	}
	statusCoordinator := &StatusCoordinator{
		status,
		statusUpdateWorker,
		sleepAfter,
	}

	go func() {
		statusCoordinator.doWork(time.Now())
		for tick := range clock {
			statusCoordinator.doWork(tick)
		}
	}()

	return statusCoordinator
}

func (s *StatusCoordinator) doWork(tick time.Time) {
	if tick.Before(s.status.LastRequest.Add(s.sleepAfter)) {
		log.Println("Working")
		s.statusUpdateWorker.UpdateStatus(tick, s.status)
		s.status.LastFetch = tick.Format(time.ANSIC)
		s.status.PreviousState = "work"
	} else {
		if s.status.PreviousState != "sleep" {
			log.Println("Sleeping")
			s.status.PreviousState = "sleep"
		}
	}
}
