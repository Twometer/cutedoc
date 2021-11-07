package utils

import "time"

type Stopwatch struct {
	startTime time.Time
}

func (stw *Stopwatch) Reset() {
	stw.startTime = time.Now()
}

func (stw *Stopwatch) Microseconds() int64 {
	return time.Now().Sub(stw.startTime).Microseconds()
}
