package shared

import (
	"time"
)

// Scanning scans for folder information
type Scanning interface {
	ScanFolders() []Album
}

// Persistence provides access to the persisted known albums
type Persistence interface {
	ReadKnownAlbums() (albums []Album, lineMap map[string]int)
	AppendNewAlbums(knownAlbums []Album, newAlbums []Album)
}

// Clock defines how often the clock ticks
type Clock interface {
	NewClock() <-chan time.Time
}

// Album has a one line string representing an album in the record collection
type Album struct {
	Text string
}
