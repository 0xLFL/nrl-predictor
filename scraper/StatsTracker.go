package main

import (
	"sync/atomic"
)

type StatsTracker struct {
	started  int32
	finished int32
}

func (s *StatsTracker) Start() {
	atomic.AddInt32(&s.started, 1)
}

func (s *StatsTracker) Finish() {
	atomic.AddInt32(&s.finished, 1)
}

func (s *StatsTracker) Started() int32 {
	return atomic.LoadInt32(&s.started)
}

func (s *StatsTracker) Finished() int32 {
	return atomic.LoadInt32(&s.finished)
}

func (s *StatsTracker) Active() int32 {
	return s.Started() - s.Finished()
}
