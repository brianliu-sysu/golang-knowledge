package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// 缓存项
type cacheItem struct {
	value      interface{}
	expireTime time.Time
}

// 带缓存的请求合并器
type CachedCoalescer struct {
	mu       sync.RWMutex
	cache    map[string]*cacheItem
	group    singleflight.Group
	cacheTTL time.Duration
}

func NewCachedCoalescer(cacheTTL time.Duration) *CachedCoalescer {
	cc := &CachedCoalescer{
		cache:    make(map[string]*cacheItem),
		cacheTTL: cacheTTL,
	}

	// 启动缓存清理
	go cc.cleanupExpired()

	return cc
}

// 获取数据（优先从缓存）
func (cc *CachedCoalescer) Get(
	ctx context.Context,
	key string,
	fn func() (interface{}, error),
) (interface{}, error) {
	// 1. 检查缓存
	cc.mu.RLock()
	if item, ok := cc.cache[key]; ok {
		if time.Now().Before(item.expireTime) {
			cc.mu.RUnlock()
			fmt.Printf("缓存命中: %s\n", key)
			return item.value, nil
		}
	}
	cc.mu.RUnlock()

	// 2. 缓存未命中，使用 singleflight 获取
	fmt.Printf("缓存未命中: %s\n", key)

	// 创建带超时的执行
	type result struct {
		val interface{}
		err error
	}
	resultCh := make(chan result, 1)

	go func() {
		val, err, _ := cc.group.Do(key, fn)

		// 存入缓存
		if err == nil {
			cc.mu.Lock()
			cc.cache[key] = &cacheItem{
				value:      val,
				expireTime: time.Now().Add(cc.cacheTTL),
			}
			cc.mu.Unlock()
		}

		resultCh <- result{val: val, err: err}
	}()

	// 3. 等待结果或超时
	select {
	case res := <-resultCh:
		return res.val, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// 清理过期缓存
func (cc *CachedCoalescer) cleanupExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		cc.mu.Lock()
		for key, item := range cc.cache {
			if now.After(item.expireTime) {
				delete(cc.cache, key)
				fmt.Printf("清理过期缓存: %s\n", key)
			}
		}
		cc.mu.Unlock()
	}
}

// 示例使用
func main() {
	cc := NewCachedCoalescer(2 * time.Second)

	// 第一轮请求（缓存未命中）
	fmt.Println("=== 第一轮请求 ===")
	for i := 0; i < 5; i++ {
		go func(id int) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			result, err := cc.Get(ctx, "data-key", func() (interface{}, error) {
				fmt.Println("执行昂贵操作...")
				time.Sleep(1 * time.Second)
				return "expensive-result", nil
			})

			if err != nil {
				fmt.Printf("请求 %d 失败: %v\n", id, err)
			} else {
				fmt.Printf("请求 %d 成功: %v\n", id, result)
			}
		}(i)
	}

	time.Sleep(2 * time.Second)

	// 第二轮请求（缓存命中）
	fmt.Println("\n=== 第二轮请求（缓存命中）===")
	for i := 5; i < 10; i++ {
		go func(id int) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			result, err := cc.Get(ctx, "data-key", func() (interface{}, error) {
				return "expensive-result", nil
			})

			if err != nil {
				fmt.Printf("请求 %d 失败: %v\n", id, err)
			} else {
				fmt.Printf("请求 %d 成功: %v\n", id, result)
			}
		}(i)
	}

	time.Sleep(2 * time.Second)

	// 第三轮请求（缓存过期）
	fmt.Println("\n=== 第三轮请求（缓存过期）===")
	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, _ := cc.Get(ctx, "data-key", func() (interface{}, error) {
		fmt.Println("执行昂贵操作（缓存已过期）...")
		time.Sleep(1 * time.Second)
		return "new-result", nil
	})
	fmt.Printf("最终结果: %v\n", result)

	time.Sleep(5 * time.Second)
}
