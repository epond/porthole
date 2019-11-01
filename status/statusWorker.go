package status

import (
	"time"
)

// RecordCollectionAdditions gets an array of new additions as strings
type RecordCollectionAdditions interface {
	FetchLatestAdditions() []string
}

// MusicStatusWorker uses RecordCollectionAdditions to update Status
type MusicStatusWorker struct {
	RecordCollectionAdditions RecordCollectionAdditions
}

// UpdateStatus updates Status using RecordCollectionAdditions
func (m *MusicStatusWorker) UpdateStatus(timestamp time.Time, status *Status) {
	status.LatestAdditions = m.RecordCollectionAdditions.FetchLatestAdditions()
}
