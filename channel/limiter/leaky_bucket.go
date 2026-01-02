package limiter

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	mu       sync.Mutex
	capacity int
	water    int
	rate     int
	lastLeak time.Time
}

func NewLeakyBucket(capacity, rate int) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		rate:     rate,
	}
}

func (l *LeakyBucket) leak() {
	now := time.Now()
	elapsed := int(now.Sub(l.lastLeak).Seconds() * float64(l.rate))

	if elapsed > 0 {
		l.water -= elapsed
		if l.water < 0 {
			l.water = 0
		}

		l.lastLeak = now
	}

}

func (l *LeakyBucket) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.leak()
	if l.water < l.capacity {
		l.water++
		return true
	}

	return false
}
