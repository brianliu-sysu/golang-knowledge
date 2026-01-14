# **Memory allocator**
The memory allocator of Golang is based on the Google's TCMalloc design, optimised for multi-core and high-concurrency scenarios.

## core design philosophy

|design | purpose |
|------|-----|
|multi-level cache|reduce lock contention|
|size classification| reduce memory fragmentation|
|batch operation|amortize allocation cost|
|thread-local cache|lock-free fast path|

## 整体架构

```mermaid
flowchart TB
    subgraph mheap["mheap (全局堆) 🔒"]
        direction TB
        H[全局唯一，管理所有内存页]
        
        subgraph mcentral["mcentral (中心缓存) 🔒"]
            direction LR
            C1["8B\nclass"]
            C2["16B\nclass"]
            C3["24B\nclass"]
            C4["32B\nclass"]
            C5["...\n..."]
            C6["32KB\nclass"]
        end
    end

    mcentral <-->|"批量获取/归还 mspan"| mcaches

    subgraph mcaches["mcache 层 (无锁访问)"]
        direction LR
        M0["mcache\n(P0)"]
        M1["mcache\n(P1)"]
        M2["mcache\n(P2)"]
        M3["mcache\n(P3)"]
    end

    M0 --- G0["Goroutine"]
    M1 --- G1["Goroutine"]
    M2 --- G2["Goroutine"]
    M3 --- G3["Goroutine"]

    style mheap fill:#ffebee,stroke:#c62828
    style mcentral fill:#fff3e0,stroke:#ef6c00
    style mcaches fill:#e8f5e9,stroke:#2e7d32
```

**层级说明：**

| 层级 | 组件 | 锁机制 | 说明 |
|------|------|--------|------|
| L1 | mheap | 全局锁 | 管理所有内存页，向 OS 申请/释放内存 |
| L2 | mcentral | 每个 size class 一把锁 | 68 个 size class，管理 mspan 链表 |
| L3 | mcache | **无锁** | 每个 P 一个，Goroutine 直接访问 |

## 分配决策流程

```mermaid
flowchart TD
    A[分配请求 size] --> B{size <= 16B<br/>且 noscan?}
    
    B -->|是| TINY["🔹 微对象分配<br/>Tiny Allocator<br/>─────────<br/>多个对象共用<br/>一个 16B 槽位"]
    
    B -->|否| C{size <= 32KB?}
    
    C -->|是| SMALL["🔸 小对象分配<br/>─────────<br/>mcache → mcentral → mheap<br/>按 size class 分配"]
    
    C -->|否| LARGE["🔺 大对象分配<br/>─────────<br/>直接从 mheap 分配<br/>独占 mspan"]

    style A fill:#e3f2fd,stroke:#1565c0
    style B fill:#fff8e1,stroke:#f9a825
    style C fill:#fff8e1,stroke:#f9a825
    style TINY fill:#e8f5e9,stroke:#2e7d32
    style SMALL fill:#fff3e0,stroke:#ef6c00
    style LARGE fill:#ffebee,stroke:#c62828
```

**分配阈值：**

| 类型 | 大小范围 | 分配路径 | 特点 |
|------|----------|----------|------|
| 微对象 | ≤16B 且 noscan | Tiny Allocator | 多对象合并，减少碎片 |
| 小对象 | 16B < size ≤ 32KB | mcache → span | 按 68 个 size class 分配 |
| 大对象 | > 32KB | 直接 mheap | 独占一个或多个 mspan |

## 小对象分配详细流程

```mermaid
flowchart TD
    START["📦 小对象分配 (size ≤ 32KB)"]
    
    START --> STEP1["1️⃣ 计算 size class<br/>───────────<br/>size → sizeclass (0-67)<br/>+ noscan → spanclass"]
    
    STEP1 --> STEP2["2️⃣ 获取当前 P 的 mcache<br/>───────────<br/><code>c := getg().m.p.mcache</code>"]
    
    STEP2 --> STEP3["3️⃣ 从 mcache 获取对应 span<br/>───────────<br/><code>span := c.alloc[spanclass]</code>"]
    
    STEP3 --> CHECK{span 有空闲槽位?}
    
    CHECK -->|有| FAST["4a. ⚡ 快速分配<br/>───────────<br/>addr = nextFreeIndex<br/>🔓 无锁！O(1)"]
    
    CHECK -->|无| REFILL["4b. 🔄 mcache.refill()<br/>───────────<br/>从 mcentral 获取新 span<br/>🔒 需加锁"]
    
    FAST --> RETURN
    REFILL --> RETURN
    
    RETURN["5️⃣ 返回对象地址<br/>───────────<br/><code>addr = span.base + idx*size</code>"]

    style START fill:#e3f2fd,stroke:#1565c0
    style CHECK fill:#fff8e1,stroke:#f9a825
    style FAST fill:#e8f5e9,stroke:#2e7d32
    style REFILL fill:#fff3e0,stroke:#ef6c00
    style RETURN fill:#f3e5f5,stroke:#7b1fa2
```

**关键路径对比：**

| 路径 | 锁机制 | 时间复杂度 | 触发条件 |
|------|--------|------------|----------|
| 快速路径 (4a) | 无锁 | O(1) | span 有空闲槽位 |
| 慢速路径 (4b) | mcentral 锁 | O(1)~O(n) | span 已满，需 refill |

## mcache 结构

```mermaid
flowchart LR
    subgraph mcache["mcache (每个 P 一个，无锁访问)"]
        direction TB
        
        subgraph tiny["🔹 Tiny Allocator"]
            T1["tiny: 0xc000100000<br/>(当前块地址)"]
            T2["tinyoffset: 8<br/>(下一个可用偏移)"]
            T3["tinyAllocs: 156<br/>(已分配计数)"]
        end
        
        subgraph alloc["📦 alloc [136]*mspan"]
            direction LR
            A0["[0] 8B scan"]
            A1["[1] 8B noscan"]
            A2["[2] 16B scan"]
            A3["[3] 16B noscan"]
            A4["..."]
            A134["[134] 32KB scan"]
            A135["[135] 32KB noscan"]
        end
    end
    
    A0 & A1 & A2 & A3 --> SPAN
    A134 & A135 --> SPAN
    
    subgraph SPAN["mspan"]
        S1["包含多个"]
        S2["object 槽位"]
    end

    style mcache fill:#e8f5e9,stroke:#2e7d32
    style tiny fill:#e3f2fd,stroke:#1565c0
    style alloc fill:#fff3e0,stroke:#ef6c00
    style SPAN fill:#f3e5f5,stroke:#7b1fa2
```

**mcache 字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `tiny` | `uintptr` | 当前 tiny 块地址，用于 ≤16B 且 noscan 的对象 |
| `tinyoffset` | `uintptr` | tiny 块内下一个可用偏移 |
| `tinyAllocs` | `uintptr` | tiny 分配计数（用于统计） |
| `alloc` | `[136]*mspan` | 68 个 size class × 2（scan/noscan）= 136 个 span 指针 |

## mspan 结构

```mermaid
flowchart TB
    subgraph mspan["mspan (内存跨度)"]
        direction TB
        
        subgraph meta["📋 元数据"]
            M1["startAddr: 0xc000100000"]
            M2["npages: 1"]
            M3["elemsize: 16B"]
            M4["nelems: 512"]
            M5["freeindex: 3"]
            M6["spanclass: 3"]
        end
        
        subgraph bitmap["🗺️ 位图"]
            B1["allocBits:  1,1,1,0,0,0,0,0,..."]
            B2["gcmarkBits: 1,1,0,0,0,0,0,0,..."]
        end
        
        subgraph memory["💾 内存布局 (8KB = 512 × 16B)"]
            direction LR
            O0["obj0<br/>✅已用"]
            O1["obj1<br/>✅已用"]
            O2["obj2<br/>✅已用"]
            O3["obj3<br/>⬜空闲"]
            O4["obj4<br/>⬜空闲"]
            O5["..."]
            O511["obj511<br/>⬜空闲"]
        end
    end
    
    M5 -.->|"freeindex=3"| O3

    style mspan fill:#fafafa,stroke:#424242
    style meta fill:#e3f2fd,stroke:#1565c0
    style bitmap fill:#fff3e0,stroke:#ef6c00
    style memory fill:#f3e5f5,stroke:#7b1fa2
    style O0 fill:#c8e6c9,stroke:#2e7d32
    style O1 fill:#c8e6c9,stroke:#2e7d32
    style O2 fill:#c8e6c9,stroke:#2e7d32
    style O3 fill:#fff9c4,stroke:#f9a825
    style O4 fill:#eceff1,stroke:#607d8b
    style O511 fill:#eceff1,stroke:#607d8b
```

**mspan 字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `startAddr` | `uintptr` | span 起始内存地址 |
| `npages` | `uintptr` | 占用的页数（每页 8KB） |
| `elemsize` | `uintptr` | 每个 object 的大小 |
| `nelems` | `uintptr` | object 总数 = npages × 8KB / elemsize |
| `freeindex` | `uintptr` | 下一个可能空闲的 object 索引（快速定位） |
| `spanclass` | `spanClass` | size class × 2 + noscan（0 或 1） |
| `allocBits` | `*gcBits` | 分配位图：1=已分配，0=空闲 |
| `gcmarkBits` | `*gcBits` | GC 标记位图：1=存活，0=待回收 |

## 分配与回收完整流程

```mermaid
flowchart TB
    subgraph lifecycle["🔄 完整生命周期"]
        direction TB
        
        OS["🖥️ OS<br/>mmap / munmap"]
        
        subgraph mheap["mheap (全局堆)"]
            direction TB
            
            subgraph pages["📄 pages (页分配器)"]
                PA["基数树 + 摘要<br/>O(1) 查找连续空闲页"]
            end
            
            subgraph mcentral["🏢 mcentral[0..135]"]
                direction LR
                subgraph C0["class 0"]
                    P0["partial<br/>有空闲"]
                    F0["full<br/>无空闲"]
                end
                subgraph C1["class 1"]
                    P1["partial"]
                    F1["full"]
                end
                subgraph CN["..."]
                    PN["partial"]
                    FN["full"]
                end
            end
        end
        
        subgraph mcaches["mcache 层 (无锁)"]
            direction LR
            M0["mcache P0<br/>alloc[136]"]
            M1["mcache P1<br/>alloc[136]"]
            M2["mcache P2<br/>alloc[136]"]
        end
        
        subgraph goroutines["Goroutines"]
            direction LR
            G0["G..."]
            G1["G..."]
            G2["G..."]
        end
    end
    
    OS <-->|"分配 ↓  回收 ↑"| mheap
    mcentral <-->|"批量获取/归还 span"| mcaches
    M0 --- G0
    M1 --- G1
    M2 --- G2

    style OS fill:#ffebee,stroke:#c62828
    style mheap fill:#fff3e0,stroke:#ef6c00
    style pages fill:#e3f2fd,stroke:#1565c0
    style mcentral fill:#f3e5f5,stroke:#7b1fa2
    style mcaches fill:#e8f5e9,stroke:#2e7d32
    style goroutines fill:#fafafa,stroke:#616161
```

**内存流动方向：**

| 操作 | 方向 | 路径 |
|------|------|------|
| **分配** | ↓ 向下 | OS → mheap → mcentral → mcache → Goroutine |
| **回收** | ↑ 向上 | Goroutine → mcache → mcentral → mheap → OS |

**mcentral 双链表：**

| 链表 | 说明 | 用途 |
|------|------|------|
| `partial` | 有空闲槽位的 span | mcache 优先从此获取 |
| `full` | 无空闲槽位的 span | GC 后可能转为 partial |

## Size Class 表（部分）

| Class | 对象大小 | span 页数 | 对象个数 | 浪费率 |
|-------|---------|----------|---------|--------|
| 1 | 8B | 1 | 1024 | 12.5% |
| 2 | 16B | 1 | 512 | 6.25% |
| 3 | 24B | 1 | 341 | 4.17% |
| 4 | 32B | 1 | 256 | 3.13% |
| 5 | 48B | 1 | 170 | 2.08% |
| 6 | 64B | 1 | 128 | 1.56% |
| ... | ... | ... | ... | ... |
| 67 | 32KB | 4 | 1 | ~0% |

> × 2 (scan/noscan) = 136 种 spanClass

## 一图总结

```mermaid
flowchart TD
    REQ["📥 分配请求"] --> SIZE{对象大小?}
    
    SIZE -->|"<16B 且 noscan"| TINY["🔹 Tiny Allocator<br/>多对象合并"]
    SIZE -->|"≤32KB"| MCACHE["🟢 mcache<br/>无锁访问"]
    SIZE -->|">32KB"| LARGE["🔴 mheap<br/>直接分配（加锁）"]
    
    MCACHE --> CHECK1{span 有空闲?}
    
    CHECK1 -->|"有"| RET1["✅ 返回地址"]
    CHECK1 -->|"无"| MCENTRAL["🟠 mcentral<br/>获取新 span（加锁）"]
    
    MCENTRAL --> CHECK2{mcentral 有 span?}
    
    CHECK2 -->|"有"| RET2["✅ 返回 span"]
    CHECK2 -->|"无"| MHEAP["🔴 mheap<br/>分配新页（加锁）"]
    
    MHEAP --> RET3["✅ 返回新 span"]
    
    RET2 --> MCACHE2["mcache 缓存 span"]
    MCACHE2 --> RET1
    
    RET3 --> MCENTRAL2["mcentral 缓存"]
    MCENTRAL2 --> RET2

    style REQ fill:#e3f2fd,stroke:#1565c0
    style SIZE fill:#fff8e1,stroke:#f9a825
    style TINY fill:#e8f5e9,stroke:#2e7d32
    style MCACHE fill:#e8f5e9,stroke:#2e7d32
    style LARGE fill:#ffebee,stroke:#c62828
    style MCENTRAL fill:#fff3e0,stroke:#ef6c00
    style MHEAP fill:#ffebee,stroke:#c62828
    style CHECK1 fill:#fce4ec,stroke:#c2185b
    style CHECK2 fill:#fce4ec,stroke:#c2185b
    style RET1 fill:#c8e6c9,stroke:#388e3c
    style RET2 fill:#c8e6c9,stroke:#388e3c
    style RET3 fill:#c8e6c9,stroke:#388e3c
```

**分配路径与锁机制：**

| 路径 | 锁机制 | 触发条件 | 性能 |
|------|--------|----------|------|
| Tiny | 无锁 | <16B 且 noscan | ⚡ 最快 |
| mcache 快速路径 | 无锁 | span 有空闲槽位 | ⚡ 极快 |
| mcentral 回填 | size class 锁 | mcache span 已满 | 🔸 较快 |
| mheap 分配 | 全局锁 | mcentral 无可用 span | 🔺 最慢 |
| 大对象直接分配 | 全局锁 | size > 32KB | 🔺 最慢 |

## 页分配器（Go 1.14+）

使用**基数树 + 位图 + 摘要**实现 O(1) 页查找：

```mermaid
flowchart TD
    subgraph radixtree["🌳 基数树 (Radix Tree)"]
        direction TB
        
        L0["Level 0 (根)<br/>━━━━━━━━━━<br/>summary<br/>覆盖整个堆"]
        
        L0 --> L1A & L1B & L1C
        
        subgraph level1["Level 1"]
            L1A["summary"]
            L1B["summary"]
            L1C["summary"]
        end
        
        L1A --> L2A
        L1B --> L2B
        L1C --> L2C
        
        subgraph level2["Level 2 ... Level N"]
            L2A["..."]
            L2B["..."]
            L2C["..."]
        end
        
        L2A --> B1
        L2B --> B2
        L2C --> B3
        
        subgraph bits["底层位图 (Bitmap)"]
            B1["bits<br/>01101001..."]
            B2["bits<br/>11110000..."]
            B3["bits<br/>00001111..."]
        end
    end

    style L0 fill:#e3f2fd,stroke:#1565c0,stroke-width:2px
    style level1 fill:#fff3e0,stroke:#ef6c00
    style level2 fill:#f3e5f5,stroke:#7b1fa2
    style bits fill:#e8f5e9,stroke:#2e7d32
```

**基数树查找原理：**

| 层级 | 内容 | 作用 |
|------|------|------|
| Level 0 | 根 summary | 记录整个堆的最大连续空闲页数 |
| Level 1~N | 子 summary | 记录子树的最大连续空闲页数 |
| 底层 | Bitmap | 每个 bit 表示一个页的状态（0=空闲，1=已用） |

**O(1) 查找流程：**
1. 从根 summary 快速判断是否有足够连续空闲页
2. 沿着满足条件的子树向下搜索
3. 在底层 bitmap 中定位具体页位置

**摘要 (pallocSum)** 编码三个值到一个 uint64：
- `start`: 从左边开始的连续空闲页数
- `max`: 区域内最大连续空闲页数
- `end`: 从右边开始的连续空闲页数

通过 `max` 可以 **O(1) 判断**该区域能否满足 n 页的需求。

## scan vs noscan

| 类型 | 含义 | GC 需要扫描内部 |
|------|------|----------------|
| scan | 对象内部有指针 | ✅ 需要 |
| noscan | 对象内部无指针 | ❌ 跳过 |

分离存储的好处：GC 可以**整个跳过 noscan span**，减少扫描开销。

## 微对象分配 (Tiny Allocator)

条件：`size < 16B && noscan`

```mermaid
block-beta
    columns 16
    
    A["int8\n1B"]:1
    B["pad\n3B"]:3
    C["int32\n4B"]:4
    D["int16\n2B"]:2
    E["空闲\n6B"]:6
    
    style A fill:#c8e6c9,stroke:#2e7d32
    style B fill:#ffccbc,stroke:#e64a19
    style C fill:#bbdefb,stroke:#1976d2
    style D fill:#d1c4e9,stroke:#7b1fa2
    style E fill:#eceff1,stroke:#607d8b
```

**Tiny 块布局 (16B)：**

| 偏移 | 大小 | 内容 | 说明 |
|------|------|------|------|
| 0 | 1B | `int8` | 第 1 个对象 |
| 1 | 3B | padding | 对齐填充（int32 需 4 字节对齐） |
| 4 | 4B | `int32` | 第 2 个对象 |
| 8 | 2B | `int16` | 第 3 个对象 |
| 10 | 6B | 空闲 | `tinyoffset = 10`，剩余空间 |

> 💡 Tiny Allocator 将多个 ≤16B 且 noscan 的小对象合并到同一个 16B 槽位，减少内存碎片

- 多个微对象共用一个 16B 槽位
- 使用 `tinyoffset` 追踪下一个可用位置
- 不区分内部边界，整个块作为一个整体管理

# **Go 逃逸分析原理**
逃逸分析是编译器在编译期判断变量应该分配到堆上还是栈上的技术，核心决策方法是：变量的生命周期是否小于函数的生命周期，小于就分配栈上，否则分配到堆上

## 典型场景
* 返回局部变量的指针
* 分配一个大的数据，导致栈空间不足
* 向chan存入指针对象
* 闭包变量
* 动态分发对象
* 切片和map存储指针

## 查看逃逸分析的结果
* go build -gcflags="-m" main.go
* go build -gcflags="-m -m" main.go
* go build -gcflags="-m" ./...


