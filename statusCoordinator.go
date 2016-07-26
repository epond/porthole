package main

import (
	"log"
	"time"
)

type StatusCoordinator struct {}

func NewStatusCoordinator(status *Status, fetchInterval int) *StatusCoordinator {
	statusCoordinator := &StatusCoordinator{}

	go func() {
		c := time.Tick(time.Duration(fetchInterval) * time.Second)
		statusCoordinator.doWork(status, time.Now())
		for now := range c {
			statusCoordinator.doWork(status, now)
		}
	}()

	return statusCoordinator
}

func (e *StatusCoordinator) doWork(status *Status, now time.Time) {
	status.Counter = status.Counter + 1
	log.Printf("Status counter:%v", status.Counter)
}
