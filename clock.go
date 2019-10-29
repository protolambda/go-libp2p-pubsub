package pubsub

import (
	"time"
)

type NanoClock interface {
	// Current unix nano second time
	Now() time.Duration
}

type SystemClock struct{}

func (sc SystemClock) Now() time.Duration {
	return time.Duration(time.Now().UnixNano())
}

