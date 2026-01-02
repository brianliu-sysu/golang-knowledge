package cache

import (
	"fmt"
	"sync"
	"testing"
)

const (
	CacheCapacity = 10000
	KeySpace      = 50000
)

func onEvict(k, v int) {
	fmt.Printf("key:%v, value:%v is expired", k, v)
}

func TestThreadLRU(t *testing.T) {
	lru := NewThreadSafeLRU[int, int](5, onEvict)

	for i := range 10 {
		lru.Put(i, i)
	}

	for i := range 10 {
		_, ok := lru.Get(i)
		if i >= 5 && !ok {
			t.Fatalf("should be ture")
		}

		if i < 5 && ok {
			t.Fatalf("should be false, %d", i)
		}
	}
}

func TestConcurrency(t *testing.T) {
	lru := NewShardedLRU[int, int](256, 2, nil, onEvict)

	wg := sync.WaitGroup{}

	for i := range 8 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := i; j < 64; j = j + 8 {
				lru.Put(j, j)
			}
		}(i)
	}

	wg.Wait()

	for i := range 64 {
		d, ok := lru.Get(i)
		if !ok || d != i {
			t.Fatalf("data is not right,key:%v, value:%v, exist:%v", i, d, ok)
		}
	}
}

func TestShardedLRU(t *testing.T) {
	lru := NewShardedLRU[int, int](8, 2, nil, onEvict)

	for i := range 8 {
		lru.Put(i, i)
	}

	for i := range 8 {
		_, ok := lru.Get(i)
		if !ok {
			t.Fatalf("should be ture")
		}
	}
}

func fastIntHasher(k int) uint64 {
	return uint64(k) * 2654435761
}

func BenchmarkLRU_NoSharding(b *testing.B) {
	lru := NewThreadSafeLRU[int, int](CacheCapacity, nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := counter % KeySpace

			if key%2 == 0 {
				lru.Put(key, key)
			} else {
				lru.Get(key)
			}
		}
	})
}

func BenchmarkLRU_sharding_16(b *testing.B) {
	lru := NewShardedLRU[int, int](CacheCapacity, 16, fastIntHasher, nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := counter % KeySpace

			if key%2 == 0 {
				lru.Put(key, key)
			} else {
				lru.Get(key)
			}
		}
	})
}

func BenchmarkLRU_sharding_256(b *testing.B) {
	lru := NewShardedLRU[int, int](CacheCapacity, 256, fastIntHasher, nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := counter % KeySpace

			if key%2 == 0 {
				lru.Put(key, key)
			} else {
				lru.Get(key)
			}
		}
	})
}
