package main

type RecordCollectionAdditions interface {
	FetchLatestAdditions() []string
}

type MusicStatusWorker struct {
	RecordCollectionAdditions RecordCollectionAdditions
}

func (m *MusicStatusWorker) UpdateStatus(status *Status) {
	status.LatestAdditions = m.RecordCollectionAdditions.FetchLatestAdditions()
}