## 引言

Go 语言社区有一句名言：“不要通过共享内存来通信，而要通过通信来共享内存。” 这让很多新手误以为 `Channel` 是解决并发的唯一银弹。

然而，阅读 Go 标准库源码你会发现，`sync` 包的使用无处不在。在追求极致性能、细粒度状态控制或实现底层数据结构（如缓存）时，`sync` 包提供的原语往往比 Channel 更高效、更直接。本文将带你重新审视 `sync` 包，掌握它的正确打开方式。

---

## 第一部分：锁的艺术 (Mutex & RWMutex)

### 1. Mutex：保护数据，而不是保护代码
很多新手习惯在函数一开头加锁，函数结束解锁。但这可能导致锁的粒度过大。
*   **最佳实践**：只锁临界区（Critical Section）。
*   **代码示例**：使用 `defer` 确保解锁，防止 panic 导致死锁。

```go
// ❌ 错误：锁住了整个耗时的 I/O 操作
func heavyOperation() {
    mu.Lock()
    defer mu.Unlock()
    http.Get("https://google.com") // 耗时操作
    count++
}

// ✅ 正确：只锁状态变更
func heavyOperation() {
    resp, _ := http.Get("https://google.com")
    
    mu.Lock()
    count++
    mu.Unlock() // 尽快释放
}
```

### 2. RWMutex：读多写少的性能救星
*   **场景**：配置热加载、缓存读取。
*   **注意**：如果写操作非常频繁，`RWMutex` 的性能可能不如普通 `Mutex`，因为读锁的维护也有开销。

---

## 第二部分：流程控制 (WaitGroup & Once)

### 1. WaitGroup：并发编排
*   **核心坑点**：**WaitGroup 绝对不能被复制！**
*   如果在函数间传递 `WaitGroup`，必须传**指针**。传值会导致死锁，因为子函数里的 `Done()` 操作的是副本。

```go
// ❌ 错误：wg sync.WaitGroup 传值
func worker(wg sync.WaitGroup) { ... }

// ✅ 正确：wg *sync.WaitGroup 传指针
func worker(wg *sync.WaitGroup) { ... }
```

### 2. Once：单例模式的最佳实现
*   相比于 `init()` 函数，`sync.Once` 可以在运行时懒加载（Lazy Loading），且线程安全。
*   **场景**：数据库连接池初始化、加载大型配置文件。

---

## 第三部分：性能优化神器 (Pool & Map)

### 1. sync.Pool：减轻 GC 压力的核武器
*   **原理**：对象复用。用完不扔，擦擦还能用。
*   **场景**：高频创建的临时对象，如 `bytes.Buffer`、Request Context。
*   **避坑指南**：
    1.  **一定要 Reset**：取出来的对象可能包含脏数据。
    2.  **不要做连接池**：`sync.Pool` 中的对象随时会被 GC 回收，不适合存数据库连接。

### 2. sync.Map：特定场景的特种兵
*   很多同学问：为什么不用 `map + Mutex`？
*   `sync.Map` 针对 **“读多写少”** 或 **“Key 集合稳定”** 的场景做了极致优化（空间换时间，读写分离）。
*   **结论**：通用场景请继续使用 `map + Mutex`，只有在性能分析（Profile）证明锁竞争是瓶颈时，才考虑切换到 `sync.Map`。

---

## 第四部分：致命陷阱 —— 禁止复制 (No Copy)

这是 `sync` 包最容易被忽视的规则：**sync 包里的结构体（Mutex, WaitGroup, Cond）在首次使用后，都不应该被复制。**

*   **原因**：锁的内部状态依赖于内存地址。复制结构体等于复制了锁的状态，但新锁和旧锁已经分离，导致并发控制失效。
*   **检测手段**：使用 `go vet`。

```go
type SafeCounter struct {
    mu sync.Mutex
    v  map[string]int
}

// ❌ 接收者是值类型，导致调用时 mu 被复制
func (c SafeCounter) Inc(key string) {
    c.mu.Lock() // 锁的是副本，根本没用！
    defer c.mu.Unlock()
    c.v[key]++
}

// ✅ 接收者必须是指针
func (c *SafeCounter) Inc(key string) { ... }
```

---

## 总结：Channel 还是 Sync？

*   **使用 Channel**：当你需要传递数据所有权、分发任务、或者实现异步消息流时。
*   **使用 Sync**：当你需要保护某个结构体的内部状态、实现高性能缓存、或者对延迟极其敏感时。

掌握 `sync` 包，是 Go 工程师从“会写业务”进阶到“能写中间件”的必经之路。
