package main

import (
	"github.com/epond/porthole/test"
	"testing"
	"time"
)

func TestStatusCoordinatorDoesWorkBeforeWaitingForFirstClockTick(t *testing.T) {
	rca := &DummyWorker{workCount: 0}
	coordinator := NewStatusCoordinator("commit", rca, make(chan time.Time, 0), 10 * time.Minute)
	coordinator.status.LastRequest = time.Now()
	time.Sleep(5 * time.Millisecond)
	test.ExpectInt(t, "workCount", 1, rca.workCount)
}

func TestStatusCoordinatorDoesWorkEveryClockTick(t *testing.T) {
	worker := &DummyWorker{workCount: 0, ticks: []time.Time{}}
	now := time.Now()
	tick1 := now.Add(1 * time.Minute)
	tick2 := now.Add(1 * time.Minute)
	tick3 := now.Add(1 * time.Minute)
	dummyClock := make(chan time.Time, 3)
	coordinator := NewStatusCoordinator("commit", worker, dummyClock, 10 * time.Minute)
	coordinator.status.LastRequest = time.Now()
	dummyClock <- tick1
	dummyClock <- tick2
	dummyClock <- tick3
	time.Sleep(5 * time.Millisecond)
	test.ExpectInt(t, "workCount", 4, worker.workCount)
	test.Expect(t, "tick1", tick1.Format(time.RFC822), worker.ticks[1].Format(time.RFC822))
	test.Expect(t, "tick2", tick2.Format(time.RFC822), worker.ticks[2].Format(time.RFC822))
	test.Expect(t, "tick3", tick3.Format(time.RFC822), worker.ticks[3].Format(time.RFC822))
}

func TestStatusCoordinatorDoesNoWorkWhenLastRequestWasALongTimeAgo(t *testing.T) {
	worker := &DummyWorker{workCount: 0, ticks: []time.Time{}}
	now := time.Now()
	tick1 := now.Add(1 * time.Minute)
	tick2 := now.Add(9 * time.Minute)
	tick3 := now.Add(11 * time.Minute)
	tick4 := now.Add(15 * time.Minute)
	tick5 := now.Add(101 * time.Minute)
	dummyClock := make(chan time.Time, 5)
	coordinator := NewStatusCoordinator("commit", worker, dummyClock, 10 * time.Minute)
	coordinator.status.LastRequest = time.Now()
	dummyClock <- tick1
	dummyClock <- tick2
	dummyClock <- tick3
	dummyClock <- tick4
	dummyClock <- tick5
	time.Sleep(5 * time.Millisecond)
	test.ExpectInt(t, "workCount", 3, worker.workCount)
	test.Expect(t, "tick1", tick1.Format(time.RFC822), worker.ticks[1].Format(time.RFC822))
	test.Expect(t, "tick2", tick2.Format(time.RFC822), worker.ticks[2].Format(time.RFC822))
}

func TestStatusCoordinatorResumesWorkWhenRequestsResume(t *testing.T) {
	worker := &DummyWorker{workCount: 0, ticks: []time.Time{}}
	now := time.Now()
	tick1 := now.Add(1 * time.Minute)
	tick2 := now.Add(9 * time.Minute)
	tick3 := now.Add(11 * time.Minute)
	tick4 := now.Add(15 * time.Minute)
	tick5 := now.Add(101 * time.Minute)
	tick6 := now.Add(109 * time.Minute)
	dummyClock := make(chan time.Time, 10)
	coordinator := NewStatusCoordinator("commit", worker, dummyClock, 10 * time.Minute)
	coordinator.status.LastRequest = time.Now()
	dummyClock <- tick1
	dummyClock <- tick2
	dummyClock <- tick3
	dummyClock <- tick4
	time.Sleep(5 * time.Millisecond)
	coordinator.status.LastRequest = now.Add(100 * time.Minute)
	dummyClock <- tick5
	dummyClock <- tick6
	time.Sleep(5 * time.Millisecond)
	test.ExpectInt(t, "workCount", 5, worker.workCount)
	test.Expect(t, "tick1", tick1.Format(time.RFC822), worker.ticks[1].Format(time.RFC822))
	test.Expect(t, "tick2", tick2.Format(time.RFC822), worker.ticks[2].Format(time.RFC822))
	test.Expect(t, "tick5", tick5.Format(time.RFC822), worker.ticks[3].Format(time.RFC822))
	test.Expect(t, "tick6", tick6.Format(time.RFC822), worker.ticks[4].Format(time.RFC822))
}

func TestLastFetchIsTimeOfLastTickThatDidWork(t *testing.T) {
	worker := &DummyWorker{workCount: 0, ticks: []time.Time{}}
	now := time.Now()
	tick1 := now.Add(9 * time.Minute)
	tick2 := now.Add(11 * time.Minute)
	dummyClock := make(chan time.Time, 2)
	coordinator := NewStatusCoordinator("commit", worker, dummyClock, 10 * time.Minute)
	coordinator.status.LastRequest = time.Now()
	dummyClock <- tick1
	dummyClock <- tick2
	time.Sleep(5 * time.Millisecond)
	test.Expect(t, "LastFetch", tick1.Format(time.ANSIC), coordinator.status.LastFetch)
}

type DummyWorker struct {
	workCount int
	ticks []time.Time
}

func (d *DummyWorker) UpdateStatus(timestamp time.Time, status *Status) {
	d.workCount += 1
	d.ticks = append(d.ticks, timestamp)
}
