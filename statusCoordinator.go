package main

import (
	"log"
	"time"
	"path"
	"fmt"
)

type StatusCoordinator struct {
	status *Status
	musicFolder string
	knownReleasesFile string
}

func NewStatusCoordinator(status *Status, musicFolder string, knownReleasesFile string, fetchInterval int) *StatusCoordinator {
	statusCoordinator := &StatusCoordinator{status, musicFolder, knownReleasesFile}

	go func() {
		c := time.Tick(time.Duration(fetchInterval) * time.Second)
		statusCoordinator.doWork()
		for _ = range c {
			statusCoordinator.doWork()
		}
	}()

	return statusCoordinator
}

func (s *StatusCoordinator) doWork() {
	s.status.LatestAdditions = LatestAdditions(s.musicFolder, s.knownReleasesFile)
	s.status.Counter = s.status.Counter + 1
	log.Printf("Status counter:%v, additions:%v", s.status.Counter, s.status.LatestAdditions)
}

func LatestAdditions(musicFolder string, knownReleasesFile string) string {
	foldersToScan := []FolderToScan{
		{path.Join(musicFolder, "flac-add"), 2},
		{path.Join(musicFolder, "flac-vorbis320"), 2}}
	folderScanList := ScanFolders(foldersToScan)
	latestReleases := UpdateKnownReleases(folderScanList, knownReleasesFile, 3)
	return fmt.Sprintf("%v, %v, %v", latestReleases[0], latestReleases[1], latestReleases[2])
}