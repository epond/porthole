package hub

import (
	"time"

	"github.com/epond/porthole/foldermusic"
	"github.com/epond/porthole/shared"
	"github.com/epond/porthole/status"
)

// Porthole is the core business logic of the porthole application
type Porthole struct {
	ticker      shared.Clock
	config      *Config
	coordinator *status.Coordinator
}

// NewPorthole creates a new Porthole
func NewPorthole(
	scanner shared.Scanning,
	persister shared.Persistence,
	ticker shared.Clock,
	config *Config,
) *Porthole {
	albumAdditions := foldermusic.NewAdditions(
		scanner,
		persister,
		config.LatestAdditionsLimit)
	statusCoordinator := status.NewCoordinator(
		config.GitCommit,
		&status.MusicStatusWorker{albumAdditions},
		ticker.NewClock(),
		time.Duration(config.SleepAfter)*time.Millisecond)
	return &Porthole{
		ticker,
		config,
		statusCoordinator,
	}
}

// GetStatus gets the current set of information returned by the application
func (p *Porthole) GetStatus() *status.Status {
	// copy current status values into a new Status instance
	return &status.Status{
		GitCommit:       p.coordinator.Status.GitCommit,
		LastRequest:     p.coordinator.Status.LastRequest,
		LastFetch:       p.coordinator.Status.LastFetch,
		LatestAdditions: p.coordinator.Status.LatestAdditions,
	}
}

// RequestScan requests that a new scan be performed
func (p *Porthole) RequestScan() {
	p.coordinator.Status.LastRequest = time.Now()
}
