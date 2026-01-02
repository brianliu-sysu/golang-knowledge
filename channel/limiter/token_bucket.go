package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu             sync.Mutex
	capacity       int
	tokens         int
	refillRate     int
	lastRefillTime time.Time
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:       capacity,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

func (t *TokenBucket) refillTokens() {
	now := time.Now()
	elapsed := int(now.Sub(t.lastRefillTime).Seconds() * float64(t.refillRate))
	if elapsed <= 0 {
		return
	}

	tokens := t.tokens + elapsed
	if tokens > t.capacity {
		tokens = t.capacity
	}

	t.tokens = tokens
	t.lastRefillTime = time.Now()
}

func (t *TokenBucket) AllowN(n int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.refillTokens()

	if t.tokens >= n {
		t.tokens -= n
		return true
	}

	return false
}

func (t *TokenBucket) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.refillTokens()

	if t.tokens > 0 {
		t.tokens--
		return true
	}

	return false
}
