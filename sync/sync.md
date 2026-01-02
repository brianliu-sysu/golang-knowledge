# **sync**

## **Mutex**
基本结构
```
// src/sync/mutex.go
type Mutex struct {
    state int32   // 锁的状态
    sema  uint32  // 信号量，用于阻塞/唤醒 goroutine
}
```
### **加锁流程**
Mutex的加锁分为两种情况，如果一个G等待Mutex超过1ms，就会进入饥饿模式
* 正常模式，Mutex尝试自旋来获取锁，最多尝试4次，如果获取不到就加入等待队列
* 饥饿模式，当前的G加入等待队尾，等待队列的第一个G获取

### **解锁流程**
首先Mutex会清除Lock状态，如果没有等待者，直接返回
* 正常模式，尝试唤醒一个G
* 饥饿模式，直接传递队头的G

## **RWMutex**
基本结构
```
// src/sync/rwmutex.go
type RWMutex struct {
    w           Mutex  // 写锁，互斥锁
    writerSem   uint32 // 写者的信号量
    readerSem   uint32 // 读者的信号量
    readerCount int32  // 当前读者数量（负数表示有写者等待）
    readerWait  int32  // 写者需要等待的读者数量
}
```

### **加锁流程**
* 读锁：
  - 尝试 readerCount + 1
  -  if > 0 获得锁
  -  if <0 阻塞
* 写锁：
  - w.lock()
  - readerCount - rwmutexMaxReaders
  - 设置readerWait
  - 如果有readerWait，阻塞
  
### **解锁流程**
* 读锁： 
  - 尝试 readerCount - 1
  - if > 0 解锁成功
  - if < 0 && readerWait - 1 = 0, 唤醒等待的写锁的G
* 写锁：
  - readerCount + rwmutexMaxReaders
  - 唤醒等待读锁的G
  - w.unlock（）

## **WaitGroup**
基本结构
```
// src/sync/waitgroup.go
type WaitGroup struct {
    noCopy noCopy
    state1 uint64    // 高 32 位：计数器，低 32 位：等待者数量
    state2 uint32    // 信号量
}
```
主要用于多任务同步

## **Once**
基本结构
```
// src/sync/once.go
type Once struct {
    done uint32    // 标记是否已执行（0 或 1）
    m    Mutex     // 互斥锁，保护初始化过程
}
```

### **实现**
* load done； if != 0 结束
* m.lock()
* load done; if != 0 结束
* 执行 fn
* set done = 1
* m.unlock()

## **Cond**
基本结构
```
// src/sync/cond.go
type Cond struct {
    noCopy  noCopy
    L       Locker      // 锁（通常是 *Mutex 或 *RWMutex）
    notify  notifyList  // 等待队列
    checker copyChecker // 检测是否被复制
}
```
主要的目的是跨协程的同步

## **Pool 对象池**
Pool的目的是为了复用对象，减少频繁分配内存造成的性能损失

### 基本结构
```
type Pool struct {
    noCopy noCopy
    
    local     unsafe.Pointer // 指向 [P]poolLocal 数组
    localSize uintptr        // local 数组的大小
    
    victim     unsafe.Pointer // 上一轮的 local
    victimSize uintptr        // victim 的大小
    
    New func() interface{}    // 创建新对象的函数
}

type poolLocal struct {
    poolLocalInternal
    pad [128 - unsafe.Sizeof(poolLocalInternal{})%128]byte
}

type poolLocalInternal struct {
    private interface{}   // 私有对象，只能被当前 P 访问
    shared  poolChain     // 共享对象链表，可以被其他 P 访问
}

type poolChain struct {
    head *poolChainElt
    tail *poolChainElt
}

type poolChainElt struct {
    poolDequeue
    next, prev *poolChainElt
}

```

### 工作原理
* Get: func (p *Pool) Get() interface{}
  - 查看当前P的private是否可用，可用直接返回
  - 查看当前P的shared是否可用，可用直接返回
  - 查看其他P的shared是否可用，可用直接返回
  - 查看victim cache，重复上面的流程
  - New一个新的对象
* Put：func (p *Pool) Put(x interface{})
  - 查看当前P的private是否为nil，为nil直接设置
  - 放入当前P的shared中

## **Map(old)**
并发安全的map

### 基本结构
```
type Map struct {
    mu Mutex
    
    // read 包含可以并发安全访问的部分
    read atomic.Value // readOnly
    
    // dirty 包含需要持有 mu 才能访问的部分
    dirty map[interface{}]*entry
    
    // misses 计数器，记录 read 未命中次数
    misses int
}

type readOnly struct {
    m       map[interface{}]*entry
    amended bool // 标记 dirty 中是否有 read 中没有的键
}

type entry struct {
    p unsafe.Pointer // *interface{}
}

```

### 工作原理
* 双map设计
  - read map 可以无锁访问
  - dirty map 包含最新数据，需要加锁访问

* 读取流程
  - 尝试从read map读取
  - 如果不存在且amended=true，加锁从dirtry map中读取
  - 记录miss次数，if miss > len(dirty map), 将dirty map提升为read map

* 写入流程
  - 如果key在read map中存在，尝试直接更新
  - 如果read map不存在，加锁更新dirty map
  - 如果dirty map为空，从read map中复制没有被删除的条目到dirty map
  
* 删除流程
  - 先标记entry为nil
  - 在dirty map提升时，标记为expunged
  
## **Map(new)**
并发安全的map

### 基本结构
```
type Map struct {
	_ noCopy

	m isync.HashTrieMap[any, any]
}

type HashTrieMap[K comparable, V any] struct {
	inited   atomic.Uint32
	initMu   Mutex
	root     atomic.Pointer[indirect[K, V]] // 每层16个桶
	keyHash  hashFunc
	valEqual equalFunc
	seed     uintptr
}

// indirect is an internal node in the hash-trie.
type indirect[K comparable, V any] struct {
	node[K, V]
	dead     atomic.Bool
	mu       Mutex // Protects mutation to children and any children that are entry nodes.
	parent   *indirect[K, V]
	children [nChildren]atomic.Pointer[node[K, V]]
}


// entry is a leaf node in the hash-trie.
type entry[K comparable, V any] struct {
	node[K, V]
	overflow atomic.Pointer[entry[K, V]] // Overflow for hash collisions.
	key      K
	value    V
}

// node is the header for a node. It's polymorphic and
// is actually either an entry or an indirect.
type node[K comparable, V any] struct {
	isEntry bool
}
```

### 工作原理
* 读取流程
  - 根据key计算hash
  - 根据hash的4位来决定存储在每层桶的位置（譬如，第一层 0 ～4 位的值，第一层 4～8，以此类推）
    - 如果位置空，直接返回
    - 如果是entry类型，判断key是否相等，相等直接返回，如果不等，遍历overflow的内容，相等返回，不等，返回空
  - 如果是中间节点，层级加一，继续上面的步骤

* 写入流程
  - 根据key计算hash
  -根据hash的4位确定当前层级的位置
    - 如果为空，直接新建一个entry，挂到父节点对应的槽位上
    - 如果是entry
      - 如果key相等，直接更新
      - 如果key不相等
        - 只是当前层级的hash相等，新建一个中间节点替代当前节点，把原来的entry和新的节点放在中间节点的槽中
        - 如果hash完全相等，只是key不等，插入entry节点的overflow上
