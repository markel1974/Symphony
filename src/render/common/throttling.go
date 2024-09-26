package common

import (
	"time"
)

type DynamicThrottling struct {
	frameInterval int64
	tuning        int64
	prev          int64
	counter       uint64
}

func NewDynamicThrottling(frameInterval int) *DynamicThrottling {
	return &DynamicThrottling{
		prev:          time.Now().UnixMilli(),
		frameInterval: int64(frameInterval),
		tuning:        0,
		counter:       0,
	}
}

func (s *DynamicThrottling) DynamicThrottling() {
	now := time.Now().UnixMilli()
	diff := now - s.prev
	interval := s.frameInterval - diff
	//https://codereview.stackexchange.com/questions/40473/portable-periodic-one-shot-timer-implementation?noredirect=1&lq=1
	//https://docs.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-signalobjectandwait
	if interval < s.frameInterval {
		duration := now + interval
		if sleep := interval - s.tuning; sleep > 1 {
			time.Sleep(time.Duration(sleep) * time.Millisecond)
			now = time.Now().UnixMilli()
		}
		if s.tuning = now - duration; s.tuning < 0 {
			s.tuning = 0
		}
	}
	s.prev = now
	s.counter++
}

func (s *DynamicThrottling) Counter() uint64 {
	return s.counter
}
