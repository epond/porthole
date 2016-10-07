package main

import (
	"testing"
	"github.com/epond/porthole/test"
	"time"
)

func TestStatusCoordinatorDoesWorkBeforeWaitingForFirstClockTick(t *testing.T) {
	rca := &DummyRCA{workCount: 0}
	NewStatusCoordinator("commit", 1, rca, make(chan time.Time, 0))
	time.Sleep(5 * time.Millisecond)
	test.ExpectInt(t, "workCount", 1, rca.workCount)
}

func TestStatusCoordinatorDoesWorkEveryClockTick(t *testing.T) {
	rca := &DummyRCA{workCount: 0}
	dummyClock := make(chan time.Time, 3)
	dummyClock <- time.Now()
	dummyClock <- time.Now()
	dummyClock <- time.Now()
	NewStatusCoordinator("commit", 1, rca, dummyClock)
	time.Sleep(5 * time.Millisecond)
	test.ExpectInt(t, "workCount", 4, rca.workCount)
}

func TestStatusCoordinatorDoesNoWorkWhenLastRequestWasALongTimeAgo(t *testing.T) {}

type DummyRCA struct {
	workCount int
}

func (d *DummyRCA) FetchLatestAdditions() []string {
	d.workCount += 1
	return []string{}
}
