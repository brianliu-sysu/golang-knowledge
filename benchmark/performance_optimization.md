# methodologies and common bottlenecks
## The golden rules of optimization
```
┌───────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                      Three principles of optimiztion                                                  │
├───────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                       │
│  1. Don't optimize prematurely ── get the code working correctly first                                │
│  2. Measure first, then optimize ── let the data speak for itself, don't rely on guesswork            │
│  3. Optimize hotspots ── 80% of the time is spent on 20% of the code                                  │
│                                                                                                       │
└───────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

## Optimize process
```
┌─────────────────────────────────────────────────────────────────┐
│                      Performance optimization process           │
└─────────────────────────────────────────────────────────────────┘
                          │
     ┌────────────────────┴────────────────────┐
     ▼                                         │
  1. define goals                              │
  "latency < 10ms，QPS > 10000"                 │
     │                                         │
     ▼                                         │
  2. Measurement status                        │
  Benchmark / pprof / monitor                  │
     │                                         │
     ▼                                         │
  3. Locate bottlenecks                        │
  CPU？memory？IO？lock？                        │
     │                                         │
     ▼                                         │
  4. Analyze the reasons                       │
  why slow？                                   │
     │                                         │
     ▼                                         │
  5. Implement optimization                    │
  Modify code                                  │
     │                                         │
     ▼                                         │
  6. verify effect ────────────────────────────┘
  Compare before and after optimization
```

## Common types of bottlenecks
```
┌────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                     Performance bottleneck classification                                                  │
├─────────────┬──────────────────────────────────────────────────────────────────────────────────────────────┤
│  CPU        │  computationally intensive,serialization,regular expressions, encryption/decrytion           │
│  Memory     │  Frequent allocation,large objects, GC pressure                                              │
│  IO         │  Disk read/write, network requests, database queries                                         │
│  Concurrent │  Lock contention, channel blocking, goroutine leaks                                          │
└─────────────┴──────────────────────────────────────────────────────────────────────────────────────────────┘
```

## CPU bottleneck
Diagnosis
```
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
```

FAQ and optimization
|question|optimization plan|
|---|---|
|Frequent serialization|use faster libraries|
|Regular match|Pre-compiled regular expressions, or use string functions instead|
|reflection operation|Buffered reflection results or generated using code |
|encryption/decryption|use hardware acceleration or reduce the encryption scope|
|Repeated calculations in a loop|Extract to outside the loop|

## Memory bottleneck
Diagnosis
```
# current heap
go tool pprof http://localhost:8080/debug/pprof/heap

# assign hotspots
go tool pprof -alloc_space http://localhost:8080/debug/pprof/heap
```

FAQ and optimization
|question|optimization plan|
|---|---|
|Frequent allocation of temporary objects|use sync.Pool to reuse objects|
|slice expansion|Preallocated capacity|
|String concatenation|strings.Builder |
|lots of small objects|structure merging, array instead of slice|
|Escape to the heap|reduce pointers, replace pointers with return values|

## IO bottleneck
diagnosis
```
# trace analysis
curl -o trace.out http://localhost:8080/debug/pprof/trace?seconds=5
go tool trace trace.out

# blocking analysis
go tool pprof http://localhost:8080/debug/pprof/block
```

FAQ and optimization
|question|optimization plan|
|---|---|
|Synchronous IO|change to asynchronous/concurrent|
|Frequent small request|batch merging|
|No buffering|Increase buffer |
|N+1 query|batch query|
|Serial request|change to parallel requests|

## parallel bottleneck
Diagnosis
```
# lock contention
go tool pprof http://localhost:8080/debug/pprof/mutex

# goroutine status
curl http://localhost:8080/debug/pprof/goroutine?debug=1

# blocking analysis
go tool pprof http://localhost:8080/debug/pprof/block
```

FAQ and optimization
|question|optimization plan|
|---|---|
|Global lock contention|Segmented lock or unlockless structure|
|Lock granularity is too large|reduce critical area|
|channel blocking|use a buffered channel|
|goroutine leak|contest cancel, timeout|
|too many goroutine|use work pool|

## GC bottleneck
diagnosis
```
GODEBUG=gctrace=1 ./app
```
FAQ and optimization
|question|optimization plan|
|---|---|
|Frequent GC|increase the value of GOGC|
|the time of STW is long|reduce live objects |
|Memory grows quickly|reuse objects, reduce object allocate|

## Quick location tool
```
┌────────────────────────────────────────────────────────────────────────┐
│                    problems → tools comparison table                   │
├─────────────────────────┬──────────────────────────────────────────────┤
│  CPU high               │  pprof profile                               │
│  memory high            │  pprof heap                                  │
│  memory increacement    │  pprof allocs + diff                         │
│  Latency jitter         │  pprof trace                                 │
│  lock contention        │  pprof mutex                                 │
│  goroutines             │  pprof goroutine                             │
│  blocking               │  pprof block                                 │
│  GC problems            │  GODEBUG=gctrace=1                           │
└─────────────────────────┴───────────────────────────────────────────────┘
```

## optimization checklist
```
□ Are there clear performance goals？
□ Has the true bottleneck been measured and identified？
□ Is the optimization targeting the hotspot paths？
□ Is there a benchmark to verify the effect？
□ Is the optimized code maintainable？
□ Have boundary cased been considered？
```

## Summarize
```
┌────────────────────────────────────────────────────────────────────┐
│                     Performance optimization core                  │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  methodologies：measure → locate → analysis → optimize → verify     │
│                                                                    │
│  tool chain：Benchmark + pprof + trace + gctrace                   │
│                                                                    │
│  priority：algorithm > IO > cocurrent > memory > micro-optimization│
│                                                                    │
│  principle：use data to speak itself，optimize hotspots，keep simple│
│                                                                    │
└────────────────────────────────────────────────────────────────────┘
```