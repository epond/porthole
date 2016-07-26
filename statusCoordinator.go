package main

import (
	"log"
	"time"
)

type StatusCoordinator struct {}

func NewStatusCoordinator(fetchInterval int) *StatusCoordinator {
	enquiryCoordinator := &StatusCoordinator{}

	go func() {
		c := time.Tick(time.Duration(fetchInterval) * time.Second)
		enquiryCoordinator.doWork(time.Now())
		for now := range c {
			enquiryCoordinator.doWork(now)
		}
	}()

	return enquiryCoordinator
}

func (e *StatusCoordinator) doWork(now time.Time) {
	log.Printf("Doing work at %v", now)
}
