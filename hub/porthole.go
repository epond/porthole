package hub

import (
	"time"

	"github.com/epond/porthole/status"
)

// Porthole is the core business logic of the porthole application
type Porthole struct {
	// Can we also have a port for the Web UI?
	scanner   Scanning
	persister Persistence
	ticker    Clock
	config    Config
}

// Scanning scans for folder information
type Scanning interface {
	ScanFolders() []status.Album
}

// Persistence provides access to the persisted known albums
type Persistence interface {
	ReadKnownAlbums() (albums []status.Album, lineMap map[string]int)
	AppendNewAlbums(knownAlbums []status.Album, newAlbums []status.Album)
}

// Clock defines how often the clock ticks
type Clock interface {
	NewClock() <-chan time.Time
}

// NewPorthole creates a new Porthole
func NewPorthole(
	scanner Scanning,
	persister Persistence,
	ticker Clock,
	config Config,
) *Porthole {
	return &Porthole{
		scanner,
		persister,
		ticker,
		config,
	}
}
