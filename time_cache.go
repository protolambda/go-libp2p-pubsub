package pubsub

import (
	"container/list"
	"time"
)

type TimeCache struct {
	q     *list.List
	m     map[string]time.Duration
	clock NanoClock
	span  time.Duration
}

func NewTimeCache(clock NanoClock, span time.Duration) *TimeCache {
	return &TimeCache{
		q:     list.New(),
		m:     make(map[string]time.Duration),
		clock: clock,
		span:  span,
	}
}

// Returns true if the message was new
func (tc *TimeCache) Add(s string) bool {
	_, ok := tc.m[s]
	if ok {
		return false
	}

	tc.sweep()

	tc.m[s] = tc.clock.Now()
	tc.q.PushFront(s)
	return true
}

func (tc *TimeCache) sweep() {
	for {
		back := tc.q.Back()
		if back == nil {
			return
		}

		v := back.Value.(string)
		t, ok := tc.m[v]
		if !ok {
			panic("inconsistent cache state")
		}

		if tc.clock.Now() - t > tc.span {
			tc.q.Remove(back)
			delete(tc.m, v)
		} else {
			return
		}
	}
}

func (tc *TimeCache) Has(s string) bool {
	_, ok := tc.m[s]
	return ok
}
