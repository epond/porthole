package shared

import (
	"time"

	"github.com/epond/porthole/status"
)

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
