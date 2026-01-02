package limiter

import (
	"sync"
	"testing"
	"time"
)

func TestNewLeakyBucket(t *testing.T) {
	capacity := 10
	rate := 5
	lb := NewLeakyBucket(capacity, rate)

	if lb == nil {
		t.Fatal("NewLeakyBucket returned nil")
	}

	if lb.capacity != capacity {
		t.Errorf("Expected capacity %d, got %d", capacity, lb.capacity)
	}

	if lb.rate != rate {
		t.Errorf("Expected rate %d, got %d", rate, lb.rate)
	}

	if lb.water != 0 {
		t.Errorf("Expected initial water to be 0, got %d", lb.water)
	}
}

func TestLeakyBucket_Allow_WithinCapacity(t *testing.T) {
	lb := NewLeakyBucket(5, 2)

	// 测试在容量范围内的请求
	for i := 0; i < 5; i++ {
		if !lb.Allow() {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 第6个请求应该被拒绝（超过容量）
	if lb.Allow() {
		t.Error("Request beyond capacity should be denied")
	}
}

func TestLeakyBucket_Allow_ExceedCapacity(t *testing.T) {
	lb := NewLeakyBucket(3, 1)

	// 填满桶
	for i := 0; i < 3; i++ {
		lb.Allow()
	}

	// 超过容量的请求应该被拒绝
	if lb.Allow() {
		t.Error("Request exceeding capacity should be denied")
	}
}

func TestLeakyBucket_Leak(t *testing.T) {
	lb := NewLeakyBucket(10, 2) // 每秒漏出2个单位
	lb.lastLeak = time.Now()

	// 填充一些水
	for i := 0; i < 5; i++ {
		lb.Allow()
	}

	if lb.water != 5 {
		t.Errorf("Expected water level 5, got %d", lb.water)
	}

	// 等待足够的时间让水漏出
	time.Sleep(2 * time.Second)

	// 调用 Allow 会触发 leak
	lb.Allow()

	// 2秒 * 2 rate = 4个单位漏出
	// 5 (原有) - 4 (漏出) + 1 (新请求) = 2
	expectedWater := 2
	if lb.water != expectedWater {
		t.Errorf("Expected water level %d after leak, got %d", expectedWater, lb.water)
	}
}

func TestLeakyBucket_LeakToZero(t *testing.T) {
	lb := NewLeakyBucket(10, 5)
	lb.lastLeak = time.Now()

	// 添加一些水
	lb.Allow()
	lb.Allow()

	// 等待足够长的时间，让所有水都漏完
	time.Sleep(1 * time.Second)

	lb.mu.Lock()
	lb.leak()
	lb.mu.Unlock()

	if lb.water != 0 {
		t.Errorf("Expected water to leak to 0, got %d", lb.water)
	}
}

func TestLeakyBucket_Concurrent(t *testing.T) {
	lb := NewLeakyBucket(100, 10)
	var wg sync.WaitGroup
	allowedCount := 0
	deniedCount := 0
	var countMu sync.Mutex

	// 并发发送200个请求
	numRequests := 200
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			if lb.Allow() {
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

	// 至少应该允许capacity数量的请求
	if allowedCount < lb.capacity {
		t.Errorf("Expected at least %d allowed requests, got %d", lb.capacity, allowedCount)
	}

	// 应该有一些请求被拒绝
	if deniedCount == 0 {
		t.Error("Expected some requests to be denied")
	}
}

func TestLeakyBucket_RefillAfterLeak(t *testing.T) {
	lb := NewLeakyBucket(5, 5) // 每秒漏5个
	lb.lastLeak = time.Now()

	// 填满桶
	for i := 0; i < 5; i++ {
		if !lb.Allow() {
			t.Fatalf("Failed to fill bucket at request %d", i+1)
		}
	}

	// 桶已满，下一个请求应该被拒绝
	if lb.Allow() {
		t.Error("Request should be denied when bucket is full")
	}

	// 等待1秒，让水漏出
	time.Sleep(1 * time.Second)

	// 现在应该可以再次允许请求
	for i := 0; i < 5; i++ {
		if !lb.Allow() {
			t.Errorf("Request %d should be allowed after leak", i+1)
		}
	}
}

func TestLeakyBucket_ZeroRate(t *testing.T) {
	lb := NewLeakyBucket(5, 0) // 漏出速率为0
	lb.lastLeak = time.Now()

	// 填满桶
	for i := 0; i < 5; i++ {
		lb.Allow()
	}

	// 等待一段时间
	time.Sleep(1 * time.Second)

	// 由于漏出速率为0，水不会减少，请求应该被拒绝
	if lb.Allow() {
		t.Error("Request should be denied when rate is 0 and bucket is full")
	}
}

func TestLeakyBucket_HighRate(t *testing.T) {
	lb := NewLeakyBucket(10, 100) // 高漏出速率
	lb.lastLeak = time.Now()

	// 填满桶
	for i := 0; i < 10; i++ {
		lb.Allow()
	}

	// 等待很短的时间
	time.Sleep(200 * time.Millisecond)

	// 由于高漏出速率，水应该已经漏完
	lb.mu.Lock()
	lb.leak()
	lb.mu.Unlock()

	if lb.water > 0 {
		t.Logf("Water level: %d (may not be 0 due to timing)", lb.water)
	}

	// 应该可以再次允许请求
	if !lb.Allow() {
		t.Error("Request should be allowed after high-rate leak")
	}
}

func BenchmarkLeakyBucket_Allow(b *testing.B) {
	lb := NewLeakyBucket(1000, 100)
	lb.lastLeak = time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb.Allow()
	}
}

func BenchmarkLeakyBucket_AllowConcurrent(b *testing.B) {
	lb := NewLeakyBucket(10000, 1000)
	lb.lastLeak = time.Now()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lb.Allow()
		}
	})
}
