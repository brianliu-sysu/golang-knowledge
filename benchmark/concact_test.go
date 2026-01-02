package main

import (
	"strings"
	"sync/atomic"
	"testing"
)

// 方式 1：+= 拼接
func BenchmarkConcatPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < 100; j++ {
			s += "a"
		}
	}
}

// 方式 2：strings.Builder
func BenchmarkConcatBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		for j := 0; j < 100; j++ {
			sb.WriteString("a")
		}
		_ = sb.String()
	}
}

// 方式 3：预分配 Builder
func BenchmarkConcatBuilderGrow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		sb.Grow(100)
		for j := 0; j < 100; j++ {
			sb.WriteString("a")
		}
		_ = sb.String()
	}
}

var count int32

func doWork() {
	atomic.AddInt32(&count, 1)
}

func BenchmarkParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doWork()
		}
	})
}
