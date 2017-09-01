package metrics

import (
	"time"
)

var startTime time.Time

func StartUp() {
	startTime = time.Now()
}

func TimeSinceStart() time.Duration {
	return time.Since(startTime)
}
