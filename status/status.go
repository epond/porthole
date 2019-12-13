package status

import (
	"log"
	"time"

	"github.com/epond/porthole/shared"
)

// Status represents the current spplication status
type Status struct {
	GitCommit       string
	LastRequest     time.Time
	LastFetch       string
	LatestAdditions []shared.Album
}

// Coordinator knows about application status and how to update it
type Coordinator struct {
	Status             *Status
	statusUpdateWorker UpdateWorker
	sleepAfter         time.Duration
}

// UpdateWorker knows how to update application status
type UpdateWorker interface {
	UpdateStatus(timestamp time.Time, status *Status)
}

// NewCoordinator constructs a new Coordinator
func NewCoordinator(
	gitCommit string,
	statusUpdateWorker UpdateWorker,
	clock <-chan time.Time,
	sleepAfter time.Duration) *Coordinator {

	status := &Status{
		GitCommit:       gitCommit,
		LastRequest:     time.Now(),
		LastFetch:       "",
		LatestAdditions: []shared.Album{},
	}
	statusCoordinator := &Coordinator{
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

func (s *Coordinator) doWork(tick time.Time) {
	if tick.Before(s.Status.LastRequest.Add(s.sleepAfter)) {
		log.Println("Working")
		s.Status.LastFetch = "in progress..."
		s.statusUpdateWorker.UpdateStatus(tick, s.Status)
		s.Status.LastFetch = tick.Format(time.ANSIC)
	}
}
