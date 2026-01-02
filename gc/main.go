package main

import (
	"fmt"
	"runtime"
	"time"

	"runtime/debug"
)

func main() {
	testGoGC(50)
	testGoGC(100)
	testGoGC(200)
	testGoGC(400)
}

func printMemStats(msg string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("--------------------------------")
	fmt.Println(msg)
	fmt.Println("Alloc:", m.Alloc/1024/1024, "MB")
	fmt.Println("TotalAlloc:", m.TotalAlloc/1024/1024, "MB")
	fmt.Println("Sys:", m.Sys/1024/1024, "MB")
	fmt.Println("Lookups:", m.Lookups)
	fmt.Println("Mallocs:", m.Mallocs/1024/1024, "MB")
	fmt.Println("Frees:", m.Frees/1024/1024, "MB")
}

func testGoGC(percent int) {
	runtime.GC()
	time.Sleep(time.Second)

	old := debug.SetGCPercent(percent)
	defer debug.SetGCPercent(old)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startGC := m.NumGC

	fmt.Println("======= GoGC=", percent, "=======")

	var data [][]byte

	for i := range 20 {
		chunk := make([]byte, 1024*1024*10)
		data = append(data, chunk)

		runtime.ReadMemStats(&m)
		fmt.Printf("Alloc: %dMB, HeapAlloc: %dMB, NextGC: %dMB, NumGC: %d\n", (i+1)*10, m.HeapAlloc/1024/1024, m.NextGC/1024/1024, m.NumGC-startGC)
	}

	data = nil
	runtime.GC()
	time.Sleep(time.Second)

	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc: %dMB, NextGC: %dMB, NumGC: %d\n", m.Alloc/1024/1024, m.NextGC/1024/1024, m.NumGC-startGC)

}
