# **pprof**
pprof是golang内置的性能分析工具，用于定位程序的CPU，内存，gorountine等性能瓶颈。

## 整体流程
```
┌─────────────────────────────────────────────────────────────────┐
│                      pprof 使用流程                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. 集成 pprof ──→ 在程序中引入 pprof                           │
│         ↓                                                       │
│  2. 采集数据 ──→ 生成 profile 文件或 HTTP 获取                   │
│         ↓                                                       │
│  3. 分析数据 ──→ 使用 go tool pprof 分析                        │
│         ↓                                                       │
│  4. 定位问题 ──→ 找到热点函数/分配点                             │
│         ↓                                                       │
│  5. 优化验证 ──→ 修复后再次分析对比                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```
## 集成pprof
* HTTP服务
```
import (
    "net/http"
    _ "net/http/pprof"
)

func main() {
    go func() {
        http.ListenAndServe(":8000", nil)
    }

    // 业务代码
}
```

**注册的路由**
```
/debug/pprof/              # 索引页
/debug/pprof/heap          # 堆内存
/debug/pprof/goroutine     # goroutine
/debug/pprof/profile       # CPU（需要采样时间）
/debug/pprof/trace         # 执行追踪
/debug/pprof/block         # 阻塞
/debug/pprof/mutex         # 锁竞争
/debug/pprof/allocs        # 累计分配
/debug/pprof/threadcreate  # 线程创建
```

* 文件输出
```
import (
    "os"
    "runtime/pprof"
)

func main() {
    // cpu profile
    f, _ := os.Create("cpu.pprof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPuProfile()

    // 业务代码 

    // heap profile
    f2, _ := os.Create("heap.pprof")
    pprof.WriteHeapProfile(f2)
    f2.Close()
}
```

## 采集数据
HTTP 方式采集
```
# CPU 采样（默认 30 秒）
curl -o cpu.pprof "http://localhost:6060/debug/pprof/profile?seconds=30"

# 堆内存
curl -o heap.pprof "http://localhost:6060/debug/pprof/heap"

# goroutine
curl -o goroutine.pprof "http://localhost:6060/debug/pprof/goroutine"

# 也可以直接分析，不保存文件
go tool pprof http://localhost:6060/debug/pprof/heap
```

**采样类型说明**

| 类型    | 说明       | 采样方式                     |
| ------- | ---------- | ---------------------------- |
| profile | CPU 采样   | 需要指定时间？seconds=30     |
| heap | 当前堆的内存快照 | 即时 |
| allocs | 累计分配统计 | 即时 |
| goroutine | goroutine栈 | 即时 |
| block | 阻塞分析 | 需开启 runtime.SetBlockProfileRate |
| mutex | 当前堆的内存快照 | 需开启 runtime.SetMutexProfileFraction |
| trace | 执行追踪 | 需要指定时间 ?seconds=5 |

## 分析数据
* 交互式分析
```
go tool pprof cpu.pprof

# 或直接连接 HTTP
go tool pprof http://localhost:6060/debug/pprof/heap
```

* 常用命令
```
(pprof) top 20          # 前 20 热点
(pprof) top -cum 20     # 按累计排序
(pprof) list funcName   # 查看函数源码
(pprof) peek funcName   # 查看调用者/被调用者
(pprof) tree            # 调用树
(pprof) web             # 生成 SVG 图（需要 graphviz）
(pprof) png             # 生成 PNG
(pprof) quit            # 退出
```

* Web UI分析
```
# 启动 Web 界面
go tool pprof -http=:8080 cpu.pprof

# 或直接从服务获取
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap
```

## 不同场景分析
* CPU
```
# 采集 30 秒 CPU 数据
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

(pprof) top 10
      flat  flat%   sum%        cum   cum%
     2.50s 25.00% 25.00%      3.20s 32.00%  main.compute
     1.80s 18.00% 43.00%      1.80s 18.00%  runtime.mallocgc
     ...
```

* 内存
```
# 当前堆使用（存活对象）
go tool pprof http://localhost:6060/debug/pprof/heap

# 累计分配（找分配热点）
go tool pprof -alloc_space http://localhost:6060/debug/pprof/heap

# 分配次数
go tool pprof -alloc_objects http://localhost:6060/debug/pprof/heap

(pprof) top
      flat  flat%   sum%        cum   cum%
   50.12MB 30.15% 30.15%    80.50MB 48.42%  main.processData
   30.25MB 18.20% 48.35%    30.25MB 18.20%  bytes.makeSlice
   ...
```

* Goroutine
```
# 查看 goroutine 数量和状态
curl "http://localhost:6060/debug/pprof/goroutine?debug=1"

# 或用 pprof 分析
go tool pprof http://localhost:6060/debug/pprof/goroutine

(pprof) top
      100  50.00%  50.00%       100  50.00%  main.worker
       80  40.00%  90.00%        80  40.00%  net/http.(*conn).serve
       ...
```

* 阻塞分析
```
# 前提：runtime.SetBlockProfileRate(1)
go tool pprof http://localhost:6060/debug/pprof/block

(pprof) top
      5.5s 45.00% 45.00%       5.5s 45.00%  sync.(*Mutex).Lock
      3.2s 26.00% 71.00%       3.2s 26.00%  runtime.chanrecv
      ...
```

* 对比分析
```
# 采集优化前
curl -o before.pprof http://localhost:6060/debug/pprof/heap

# 优化后再采集
curl -o after.pprof http://localhost:6060/debug/pprof/heap

# 对比分析
go tool pprof -diff_base=before.pprof after.pprof

(pprof) top
      flat  flat%   sum%        cum   cum%
   -20.5MB -30.0% -30.0%    -25.0MB -35.0%  main.processData  # 负数=减少
   ...
```

* 火焰图
```
# 方式 1: go tool pprof 自带
go tool pprof -http=:8080 cpu.pprof
# 然后选择 "Flame Graph" 视图

# 方式 2: 命令行
go tool pprof -png -output=flame.png cpu.pprof
```

* trace
```
# 采集 5 秒执行追踪
curl -o trace.out "http://localhost:6060/debug/pprof/trace?seconds=5"

# 分析
go tool trace trace.out

# 打开浏览器查看：
# - View trace：时间线
# - Goroutine analysis：goroutine 分析
# - Network blocking profile：网络阻塞
# - Synchronization blocking profile：同步阻塞
# - Syscall blocking profile：系统调用阻塞
```

## 速查表
```
┌─────────────────────────────────────────────────────────────────┐
│                      pprof 速查表                                │
├───────────────┬─────────────────────────────────────────────────┤
│  采集 CPU     │  /debug/pprof/profile?seconds=30               │
│  采集内存     │  /debug/pprof/heap                              │
│  采集分配     │  /debug/pprof/allocs                            │
│  采集 goroutine│ /debug/pprof/goroutine                         │
├───────────────┼─────────────────────────────────────────────────┤
│  top          │  热点排序                                        │
│  list func    │  函数源码                                        │
│  web          │  SVG 调用图                                      │
│  -http=:8080  │  Web UI                                          │
│  -diff_base   │  对比分析                                        │
├───────────────┼─────────────────────────────────────────────────┤
│  -inuse_space │  当前使用的内存                                  │
│  -alloc_space │  累计分配的内存                                  │
│  -alloc_objects│ 累计分配的对象数                                │
└───────────────┴─────────────────────────────────────────────────┘
```