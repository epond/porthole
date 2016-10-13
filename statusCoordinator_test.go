package main

import (
	"github.com/epond/porthole/test"
	"testing"
	"time"
)

func TestStatusCoordinatorDoesWorkBeforeWaitingForFirstClockTick(t *testing.T) {
	rca := &DummyWorker{workCount: 0}
	NewStatusCoordinator("commit", rca, make(chan time.Time, 0))
	time.Sleep(5 * time.Millisecond)
	test.ExpectInt(t, "workCount", 1, rca.workCount)
}

func TestStatusCoordinatorDoesWorkEveryClockTick(t *testing.T) {
	worker := &DummyWorker{workCount: 0, ticks: []time.Time{}}
	tick1 := time.Now().Add(1 * time.Minute)
	tick2 := time.Now().Add(1 * time.Minute)
	tick3 := time.Now().Add(1 * time.Minute)
	dummyClock := make(chan time.Time, 3)
	dummyClock <- tick1
	dummyClock <- tick2
	dummyClock <- tick3
	NewStatusCoordinator("commit", worker, dummyClock)
	time.Sleep(5 * time.Millisecond)
	test.ExpectInt(t, "workCount", 4, worker.workCount)
	test.Expect(t, "tick1", tick1.Format(time.RFC822), worker.ticks[1].Format(time.RFC822))
	test.Expect(t, "tick1", tick2.Format(time.RFC822), worker.ticks[2].Format(time.RFC822))
	test.Expect(t, "tick1", tick3.Format(time.RFC822), worker.ticks[3].Format(time.RFC822))
}

func TestStatusCoordinatorDoesNoWorkWhenLastRequestWasALongTimeAgo(t *testing.T) {}

type DummyWorker struct {
	workCount int
	ticks []time.Time
}

func (d *DummyWorker) UpdateStatus(timestamp time.Time, status *Status) {
	d.workCount += 1
	d.ticks = append(d.ticks, timestamp)
}
