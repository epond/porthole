package main

import "time"

// ClockTicker ticks regularly at the specified interval
type ClockTicker struct {
	interval time.Duration
}

// NewClockTicker creates a new ClockTicker
func NewClockTicker(interval time.Duration) *ClockTicker {
	return &ClockTicker{interval}
}

// NewClock creates a new channel that ticks the time
func (c *ClockTicker) NewClock() <-chan time.Time {
	return time.Tick(c.interval)
}
