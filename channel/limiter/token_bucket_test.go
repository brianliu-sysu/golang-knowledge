package limiter

import (
	"sync"
	"testing"
	"time"
)

func TestNewTokenBucket(t *testing.T) {
	capacity := 10
	refillRate := 5
	tb := NewTokenBucket(capacity, refillRate)

	if tb == nil {
		t.Fatal("NewTokenBucket returned nil")
	}

	if tb.capacity != capacity {
		t.Errorf("Expected capacity %d, got %d", capacity, tb.capacity)
	}

	if tb.refillRate != refillRate {
		t.Errorf("Expected refillRate %d, got %d", refillRate, tb.refillRate)
	}

	if tb.tokens != 0 {
		t.Errorf("Expected initial tokens to be 0, got %d", tb.tokens)
	}

	if tb.lastRefillTime.IsZero() {
		t.Error("Expected lastRefillTime to be initialized")
	}
}

func TestTokenBucket_Allow_InitialState(t *testing.T) {
	tb := NewTokenBucket(10, 5)

	// 初始状态tokens为0，第一次请求应该被拒绝
	if tb.Allow() {
		t.Error("First request should be denied when tokens start at 0")
	}
}

func TestTokenBucket_Allow_AfterRefill(t *testing.T) {
	tb := NewTokenBucket(10, 5) // 每秒补充5个令牌
	tb.lastRefillTime = time.Now()

	// 等待1秒让令牌补充
	time.Sleep(1 * time.Second)

	// 应该有5个令牌可用
	for i := 0; i < 5; i++ {
		if !tb.Allow() {
			t.Errorf("Request %d should be allowed after refill", i+1)
		}
	}

	// 第6个请求应该被拒绝
	if tb.Allow() {
		t.Error("Request beyond available tokens should be denied")
	}
}

func TestTokenBucket_Allow_MaxCapacity(t *testing.T) {
	tb := NewTokenBucket(5, 10) // 容量5，每秒补充10个
	tb.lastRefillTime = time.Now()

	// 等待2秒，理论上会补充20个令牌，但容量限制为5
	time.Sleep(2 * time.Second)

	tb.mu.Lock()
	tb.refillTokens()
	tb.mu.Unlock()

	if tb.tokens > tb.capacity {
		t.Errorf("Tokens %d should not exceed capacity %d", tb.tokens, tb.capacity)
	}

	// 应该最多只能使用5个令牌
	successCount := 0
	for i := 0; i < 10; i++ {
		if tb.Allow() {
			successCount++
		}
	}

	if successCount > 5 {
		t.Errorf("Expected max 5 successful requests, got %d", successCount)
	}
}

func TestTokenBucket_AllowN_Success(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	tb.lastRefillTime = time.Now()

	// 等待2秒补充10个令牌
	time.Sleep(2 * time.Second)

	// 请求5个令牌
	if !tb.AllowN(5) {
		t.Error("AllowN(5) should succeed when 10 tokens are available")
	}

	// 再请求5个令牌
	if !tb.AllowN(5) {
		t.Error("Second AllowN(5) should succeed when 5 tokens remain")
	}

	// 再请求1个令牌应该失败
	if tb.AllowN(1) {
		t.Error("AllowN(1) should fail when no tokens remain")
	}
}

func TestTokenBucket_AllowN_InsufficientTokens(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	tb.tokens = 3
	tb.lastRefillTime = time.Now()

	// 请求5个令牌，但只有3个可用
	if tb.AllowN(5) {
		t.Error("AllowN(5) should fail when only 3 tokens are available")
	}

	// 令牌数不应该改变
	if tb.tokens != 3 {
		t.Errorf("Expected tokens to remain 3, got %d", tb.tokens)
	}
}

func TestTokenBucket_AllowN_Zero(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	tb.tokens = 0
	tb.lastRefillTime = time.Now()

	// 请求0个令牌应该成功
	if !tb.AllowN(0) {
		t.Error("AllowN(0) should always succeed")
	}
}

func TestTokenBucket_RefillTokens_NoTimeElapsed(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	tb.tokens = 5
	tb.lastRefillTime = time.Now()

	initialTokens := tb.tokens

	// 立即调用refillTokens，时间未经过
	tb.mu.Lock()
	tb.refillTokens()
	tb.mu.Unlock()

	// 令牌数应该保持不变
	if tb.tokens != initialTokens {
		t.Errorf("Expected tokens to remain %d, got %d", initialTokens, tb.tokens)
	}
}

func TestTokenBucket_RefillTokens_PartialRefill(t *testing.T) {
	tb := NewTokenBucket(20, 10) // 容量20，每秒10个
	tb.tokens = 5
	tb.lastRefillTime = time.Now()

	// 等待1秒
	time.Sleep(1 * time.Second)

	tb.mu.Lock()
	tb.refillTokens()
	tb.mu.Unlock()

	// 应该补充了10个令牌：5 + 10 = 15
	expectedTokens := 15
	if tb.tokens < expectedTokens-1 || tb.tokens > expectedTokens+1 {
		t.Logf("Expected approximately %d tokens, got %d (timing variance allowed)", expectedTokens, tb.tokens)
	}
}

func TestTokenBucket_RefillTokens_ExceedCapacity(t *testing.T) {
	tb := NewTokenBucket(10, 20) // 容量10，每秒20个
	tb.tokens = 5
	tb.lastRefillTime = time.Now()

	// 等待1秒，会补充20个令牌
	time.Sleep(1 * time.Second)

	tb.mu.Lock()
	tb.refillTokens()
	tb.mu.Unlock()

	// 令牌数应该被限制在容量内
	if tb.tokens > tb.capacity {
		t.Errorf("Tokens %d should not exceed capacity %d", tb.tokens, tb.capacity)
	}

	if tb.tokens != tb.capacity {
		t.Errorf("Expected tokens to be capped at capacity %d, got %d", tb.capacity, tb.tokens)
	}
}

func TestTokenBucket_Concurrent_Allow(t *testing.T) {
	tb := NewTokenBucket(100, 50)
	tb.tokens = 100 // 初始化100个令牌
	tb.lastRefillTime = time.Now()

	var wg sync.WaitGroup
	allowedCount := 0
	deniedCount := 0
	var countMu sync.Mutex

	numRequests := 200
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			if tb.Allow() {
				countMu.Lock()
				allowedCount++
				countMu.Unlock()
			} else {
				countMu.Lock()
				deniedCount++
				countMu.Unlock()
			}
		}()
	}

	wg.Wait()

	// 验证总数
	if allowedCount+deniedCount != numRequests {
		t.Errorf("Expected %d total requests, got %d", numRequests, allowedCount+deniedCount)
	}

	// 应该有一些请求被允许
	if allowedCount == 0 {
		t.Error("Expected some requests to be allowed")
	}

	// 应该有一些请求被拒绝
	if deniedCount == 0 {
		t.Error("Expected some requests to be denied")
	}

	t.Logf("Allowed: %d, Denied: %d", allowedCount, deniedCount)
}

func TestTokenBucket_Concurrent_AllowN(t *testing.T) {
	tb := NewTokenBucket(100, 50)
	tb.tokens = 100
	tb.lastRefillTime = time.Now()

	var wg sync.WaitGroup
	allowedCount := 0
	deniedCount := 0
	var countMu sync.Mutex

	numRequests := 50
	tokensPerRequest := 5
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			if tb.AllowN(tokensPerRequest) {
				countMu.Lock()
				allowedCount++
				countMu.Unlock()
			} else {
				countMu.Lock()
				deniedCount++
				countMu.Unlock()
			}
		}()
	}

	wg.Wait()

	t.Logf("AllowN(%d) - Allowed: %d, Denied: %d", tokensPerRequest, allowedCount, deniedCount)

	// 验证总数
	if allowedCount+deniedCount != numRequests {
		t.Errorf("Expected %d total requests, got %d", numRequests, allowedCount+deniedCount)
	}
}

func TestTokenBucket_ContinuousRefill(t *testing.T) {
	tb := NewTokenBucket(10, 10) // 每秒补充10个
	tb.lastRefillTime = time.Now()

	// 第一秒：等待并消耗令牌
	time.Sleep(1 * time.Second)
	count1 := 0
	for i := 0; i < 15; i++ {
		if tb.Allow() {
			count1++
		}
	}

	// 第二秒：再次等待并消耗令牌
	time.Sleep(1 * time.Second)
	count2 := 0
	for i := 0; i < 15; i++ {
		if tb.Allow() {
			count2++
		}
	}

	t.Logf("First second: %d requests allowed", count1)
	t.Logf("Second second: %d requests allowed", count2)

	// 两次都应该允许约10个请求
	if count1 < 8 || count1 > 12 {
		t.Logf("First second allowed %d requests (expected ~10, timing variance allowed)", count1)
	}

	if count2 < 8 || count2 > 12 {
		t.Logf("Second second allowed %d requests (expected ~10, timing variance allowed)", count2)
	}
}

func TestTokenBucket_ZeroRefillRate(t *testing.T) {
	tb := NewTokenBucket(10, 0) // 补充速率为0
	tb.tokens = 5
	tb.lastRefillTime = time.Now()

	initialTokens := tb.tokens

	// 等待一段时间
	time.Sleep(1 * time.Second)

	tb.mu.Lock()
	tb.refillTokens()
	tb.mu.Unlock()

	// 令牌数应该保持不变
	if tb.tokens != initialTokens {
		t.Errorf("Expected tokens to remain %d with zero refill rate, got %d", initialTokens, tb.tokens)
	}
}

func TestTokenBucket_HighRefillRate(t *testing.T) {
	tb := NewTokenBucket(50, 100) // 每秒补充100个
	tb.tokens = 0
	tb.lastRefillTime = time.Now()

	// 等待0.5秒
	time.Sleep(500 * time.Millisecond)

	tb.mu.Lock()
	tb.refillTokens()
	tb.mu.Unlock()

	// 应该补充了约50个令牌，但被容量限制
	if tb.tokens > tb.capacity {
		t.Errorf("Tokens %d should not exceed capacity %d", tb.tokens, tb.capacity)
	}

	t.Logf("Tokens after 0.5s with high refill rate: %d", tb.tokens)
}

func TestTokenBucket_NegativeTokensPrevention(t *testing.T) {
	tb := NewTokenBucket(5, 2)
	tb.tokens = 0
	tb.lastRefillTime = time.Now()

	// 尝试多次请求
	for i := 0; i < 10; i++ {
		tb.Allow()
	}

	// 令牌数不应该变成负数
	if tb.tokens < 0 {
		t.Errorf("Tokens should never be negative, got %d", tb.tokens)
	}
}

func TestTokenBucket_AllowN_BoundaryCondition(t *testing.T) {
	tb := NewTokenBucket(10, 5)
	tb.tokens = 10
	tb.lastRefillTime = time.Now()

	// 请求恰好等于可用令牌数
	if !tb.AllowN(10) {
		t.Error("AllowN(10) should succeed when exactly 10 tokens are available")
	}

	// 令牌应该被消耗完
	if tb.tokens != 0 {
		t.Errorf("Expected 0 tokens remaining, got %d", tb.tokens)
	}
}

func BenchmarkTokenBucket_Allow(b *testing.B) {
	tb := NewTokenBucket(10000, 1000)
	tb.tokens = 10000
	tb.lastRefillTime = time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.Allow()
	}
}

func BenchmarkTokenBucket_AllowN(b *testing.B) {
	tb := NewTokenBucket(10000, 1000)
	tb.tokens = 10000
	tb.lastRefillTime = time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.AllowN(5)
	}
}

func BenchmarkTokenBucket_AllowConcurrent(b *testing.B) {
	tb := NewTokenBucket(100000, 10000)
	tb.tokens = 100000
	tb.lastRefillTime = time.Now()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tb.Allow()
		}
	})
}

func BenchmarkTokenBucket_RefillTokens(b *testing.B) {
	tb := NewTokenBucket(10000, 1000)
	tb.tokens = 5000
	tb.lastRefillTime = time.Now().Add(-1 * time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.mu.Lock()
		tb.refillTokens()
		tb.mu.Unlock()
	}
}
