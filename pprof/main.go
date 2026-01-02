package main

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	compute()
}

func init() {
	runtime.SetBlockProfileRate(1)     // 开启阻塞分析
	runtime.SetMutexProfileFraction(1) // 开启锁分析
}

func compute() {
	for {
		arr := make([]byte, 1024*1024*10)
		_ = arr
		time.Sleep(time.Millisecond * 100)
	}
}
