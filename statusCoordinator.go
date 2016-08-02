package main

import (
	"log"
	"time"
)

type StatusCoordinator struct {
	status *Status
	musicFolder string
}

func NewStatusCoordinator(status *Status, musicFolder string, fetchInterval int) *StatusCoordinator {
	statusCoordinator := &StatusCoordinator{status, musicFolder}

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
	s.status.LatestAdditions = LatestAdditions(s.musicFolder)
	s.status.Counter = s.status.Counter + 1
	log.Printf("Status counter:%v", s.status.Counter)
}
