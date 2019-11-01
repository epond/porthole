package status

import (
	"log"
	"time"
)

type Status struct {
	GitCommit       string
	LastRequest     time.Time
	LastFetch       string
	LatestAdditions []string
}

type StatusCoordinator struct {
	Status             *Status
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
	if tick.Before(s.Status.LastRequest.Add(s.sleepAfter)) {
		log.Println("Working")
		s.Status.LastFetch = "in progress..."
		s.statusUpdateWorker.UpdateStatus(tick, s.Status)
		s.Status.LastFetch = tick.Format(time.ANSIC)
	}
}
