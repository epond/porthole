package status

import (
	"time"
)

// Album has a one line string representing an album in the record collection
type Album struct {
	Text string
}

// AlbumAdditions gets an array of new additions as strings
type AlbumAdditions interface {
	FetchLatestAdditions() []Album
}

// MusicStatusWorker uses AlbumAdditions to update Status
type MusicStatusWorker struct {
	AlbumAdditions AlbumAdditions
}

// UpdateStatus updates Status using AlbumAdditions
func (m *MusicStatusWorker) UpdateStatus(timestamp time.Time, status *Status) {
	status.LatestAdditions = m.AlbumAdditions.FetchLatestAdditions()
}
