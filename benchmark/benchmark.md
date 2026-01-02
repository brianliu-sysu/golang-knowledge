# **Benchmark**
Benchmark is a tool used to measure code perference, telling you how fast your code runs and how much memory it uses.

## Basic structure
```
// xxx_test.go
package xxx

func BenchmarkFuncName(b *testing.B) {
    // optional： Initialization
    setup()

    b.ResetTimer() // reset the timer to exclude the initialization time

    for i := range b.N {
        // code under test
        funcToTest()
    }
}

// subtest
func BenchmarkConcat(b *testing.B) {
    sizes := []{10, 100, 1000}

    for _, size := range sizes {
        b.Run(fmt.Sprint("size=%d", size), func(b *testing.B) {
            for i := range b.N {
                concat(size)
            }
        })
    }
}

// parallel testing
func BenchnarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            doWork()
        }
    })
}
```
Key Points:
* the filename must end with `_test.go`
* the function name must start with `Benchmark`
* the parameter must be `*testing.B`
* Loop b.N times (the frame automatically adjust the value of N)

## run Benchamark
```
# run all benchmark
go test -bench=.

# run a special benchmark
go test -bench=BenchmarkFuncName

# regular match
go test -bench="Concat"

# specify run time
go test -bench=. -benchtime=5s

# specify the number of runs
go test -bench=. -benchtime=1000x

# take the average after multiple runs
go test -bench=. -count=5
```
## benchmem: memory analysis
```
go test -bench=. -benchmem
```

**Output interpretation**
```
BenchmarkConcat-8    5000000    300 ns/op    64 B/op    2 allocs/op
│             │      │          │           │          │
│             │      │          │           │          └─ 每次操作分配次数
│             │      │          │           └─ 每次操作分配字节数
│             │      │          └─ 每次操作耗时
│             │      └─ 运行次数
│             └─ CPU 核心数
└─ 函数名
```

## generate profile
```
# generate CPU profile
go test -bench=. -cpuprofile=cpu.pprof

# generate memory profile
go test -bench=. -memprofile=mem.pprof

# analysis
go tool pprof cpu.pprof
go tool pprof mem.pprof
```

## comparesion tool benchstat
```
# install
go install golang.org/x/perf/cmd/benchstat@latest

# run the program multiple times and save the results
go test -bench=. -count=10 > old.txt

# Run it again after optimization
go test -bench=. -count=10 > new.txt

# Comparison
benchstat old.txt new.txt
```

## Best practices
```
┌─────────────────────────────────────────────────────────────────┐
│                    Benchmark best practices                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. always use -benchmem to check memeory allocation                                 │
│  2. use b.ResetTimer() to exclude initialization time                                │
│  3. use -count=5 to run multiple time to calculate the average                                │
│  4. use benchstat to  compare optimization results                                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```