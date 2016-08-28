package main

import (
	"log"
	"time"
	"strings"
	"path"
	"strconv"
)

type StatusCoordinator struct {
	status *Status
	musicFolder string
	knownReleasesFile string
	foldersToScan string
}

// TODO decouple status coordinator from record collection additions
func NewStatusCoordinator(status *Status, musicFolder string, knownReleasesFile string, foldersToScan string, fetchInterval int) *StatusCoordinator {
	statusCoordinator := &StatusCoordinator{status, musicFolder, knownReleasesFile, foldersToScan}

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
	foldersToScan := parseFoldersToScan(s.musicFolder, s.foldersToScan)
	s.status.LatestAdditions = UpdateKnownReleases(ScanFolders(foldersToScan), s.knownReleasesFile, 3)
	s.status.Counter = s.status.Counter + 1
	log.Printf("Status counter:%v, additions:%v", s.status.Counter, s.status.LatestAdditions)
}

func parseFoldersToScan(musicFolder string, folders string) []FolderToScan {
	var foldersToScan []FolderToScan
	folderPairs := strings.Split(folders, ",")
	for _, pair := range folderPairs {
		pairArray := strings.Split(pair, ":")
		depth, _ := strconv.Atoi(pairArray[1])
		foldersToScan = append(foldersToScan, FolderToScan{
			path.Join(musicFolder, pairArray[0]),
			depth})
	}
	return foldersToScan
}