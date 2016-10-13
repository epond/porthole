package main

import "time"

type RecordCollectionAdditions interface {
	FetchLatestAdditions() []string
}

type MusicStatusWorker struct {
	RecordCollectionAdditions RecordCollectionAdditions
}

func (m *MusicStatusWorker) UpdateStatus(timestamp time.Time, status *Status) {
	status.LatestAdditions = m.RecordCollectionAdditions.FetchLatestAdditions()
}
