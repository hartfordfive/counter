package counter

import (
	"time"
)

// Run at: https://play.golang.org/p/mRbfZD7Dt5

type RateCounter struct {
	counter   *Counter
	period    time.Duration
	terminate chan bool
}

func NewRateCounter(p time.Duration) *RateCounter {
	cnt := &RateCounter{
		counter:   NewCounter(),
		period:    p,
		terminate: make(chan bool, 1),
	}
	go func(cnt *RateCounter) {
		ticker := time.NewTicker(cnt.period)
		for {
			select {
			case <-ticker.C:
				cnt.counter.Reset()
			case <-cnt.terminate:
				ticker.Stop()
				cnt.counter.Cancel()
				close(cnt.terminate)
				return
			}
		}
	}(cnt)
	return cnt
}

func (rc *RateCounter) Incr(val int64) {
	rc.counter.Incr(val)
}

func (rc *RateCounter) CurrRate() int64 {
	return rc.counter.Value()
}

func (rc *RateCounter) Cancel() {
	rc.terminate <- true
}
