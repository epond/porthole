package main

import (
	"log"
	"time"
	"path"
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

func LatestAdditions(musicFolder string, knownReleasesFile string) []string {
	foldersToScan := []FolderToScan{
		//{path.Join(musicFolder, "flac"), 3},
		//{path.Join(musicFolder, "flac-cd"), 3},
		//{path.Join(musicFolder, "flac-add"), 2},
		{path.Join(musicFolder, "flac-vorbis320"), 2},
		{path.Join(musicFolder, "mp3", "main"), 2},
	}
	return UpdateKnownReleases(ScanFolders(foldersToScan), knownReleasesFile, 3)
}