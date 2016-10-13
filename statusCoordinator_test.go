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
	rca := &DummyWorker{workCount: 0}
	dummyClock := make(chan time.Time, 3)
	dummyClock <- time.Now()
	dummyClock <- time.Now()
	dummyClock <- time.Now()
	NewStatusCoordinator("commit", rca, dummyClock)
	time.Sleep(5 * time.Millisecond)
	test.ExpectInt(t, "workCount", 4, rca.workCount)
}

func TestStatusCoordinatorDoesNoWorkWhenLastRequestWasALongTimeAgo(t *testing.T) {}

type DummyWorker struct {
	workCount int
}

func (d *DummyWorker) UpdateStatus(status *Status) {
	d.workCount += 1
}
