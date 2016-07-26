package main

import (
	"log"
	"time"
)

type StatusCoordinator struct {}

func NewStatusCoordinator(status *Status, fetchInterval int) *StatusCoordinator {
	enquiryCoordinator := &StatusCoordinator{}

	go func() {
		c := time.Tick(time.Duration(fetchInterval) * time.Second)
		enquiryCoordinator.doWork(status, time.Now())
		for now := range c {
			enquiryCoordinator.doWork(status, now)
		}
	}()

	return enquiryCoordinator
}

func (e *StatusCoordinator) doWork(status *Status, now time.Time) {
	status.Counter = status.Counter + 1
	log.Printf("Status counter:%v", status.Counter)
}
