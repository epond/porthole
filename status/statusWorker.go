package status

import (
	"time"
)

// AlbumAdditions gets an array of new additions as strings
type AlbumAdditions interface {
	FetchLatestAdditions() []string
}

// MusicStatusWorker uses AlbumAdditions to update Status
type MusicStatusWorker struct {
	AlbumAdditions AlbumAdditions
}

// UpdateStatus updates Status using AlbumAdditions
func (m *MusicStatusWorker) UpdateStatus(timestamp time.Time, status *Status) {
	status.LatestAdditions = m.AlbumAdditions.FetchLatestAdditions()
}
