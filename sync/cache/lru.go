package cache

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type Node[K comparable, V any] struct {
	key  K
	data V

	prev *Node[K, V]
	next *Node[K, V]
}

func NewNode[K comparable, V any](key K, data V) *Node[K, V] {
	return &Node[K, V]{
		key:  key,
		data: data,
	}
}

type ThreadSafeLRU[K comparable, V any] struct {
	mu       sync.Mutex
	cache    map[K]*Node[K, V]
	capacity int

	head *Node[K, V]
	tail *Node[K, V]

	onEvict func(K, V)
}

func NewThreadSafeLRU[K comparable, V any](capacity int, onEvict func(K, V)) *ThreadSafeLRU[K, V] {
	if capacity <= 0 {
		panic("capacity must be > 0")
	}

	return &ThreadSafeLRU[K, V]{
		capacity: capacity,
		cache:    make(map[K]*Node[K, V]),
		onEvict:  onEvict,
	}
}

func (l *ThreadSafeLRU[K, V]) Put(key K, value V) {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, ok := l.cache[key]
	if !ok {
		newNode := NewNode(key, value)
		l.addFront(newNode)
		l.cache[key] = newNode
	} else {
		node.data = value
		l.removeNode(node)
		l.addFront(node)
	}

	if len(l.cache) > l.capacity {
		node := l.removeTail()
		if node != nil {
			delete(l.cache, node.key)

			if l.onEvict != nil {
				l.onEvict(node.key, node.data)
			}
		}
	}
}

func (l *ThreadSafeLRU[K, V]) Get(key K) (V, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, ok := l.cache[key]
	if ok {
		l.removeNode(node)
		l.addFront(node)
		return node.data, true
	}

	var zero V
	return zero, false
}

func (l *ThreadSafeLRU[K, V]) removeNode(node *Node[K, V]) {
	if node == nil {
		return
	}

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}

	node.prev = nil
	node.next = nil
}

func (l *ThreadSafeLRU[K, V]) addFront(node *Node[K, V]) {
	node.prev = nil

	if l.head == nil {
		l.head = node
		l.tail = node
		node.next = nil
		return
	}

	node.next = l.head
	l.head.prev = node
	l.head = node
}

func (l *ThreadSafeLRU[K, V]) removeTail() *Node[K, V] {
	node := l.tail
	if node == nil {
		return nil
	}

	l.removeNode(node)
	return node
}

type Hasher[K comparable] func(K) uint64

func NewDefaultHasher[K comparable]() Hasher[K] {
	return func(k K) uint64 {
		h := fnv.New64()
		_, _ = h.Write([]byte(fmt.Sprint(k)))
		return h.Sum64()
	}
}

// sharedLRU
type ShardedLRU[K comparable, V any] struct {
	shards [](*ThreadSafeLRU[K, V])
	masks  uint64
	hasher Hasher[K]
}

func NewShardedLRU[K comparable, V any](capacity,
	shards int,
	hasher Hasher[K],
	onEvict func(K, V),
) *ShardedLRU[K, V] {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}

	power := 1
	for power < shards {
		power <<= 1
	}

	shards = power

	shardedLRU := &ShardedLRU[K, V]{
		shards: make([](*ThreadSafeLRU[K, V]), shards),
		masks:  uint64(shards - 1),
		hasher: hasher,
	}

	if hasher == nil {
		shardedLRU.hasher = NewDefaultHasher[K]()
	}

	basCap := capacity / shards
	if capacity%shards > 0 {
		basCap++
	}

	if basCap <= 0 {
		basCap = 1
	}

	for i := range shards {
		shardedLRU.shards[i] = NewThreadSafeLRU[K, V](basCap, onEvict)
	}

	return shardedLRU
}

func (s *ShardedLRU[K, V]) shardForKey(k K) *ThreadSafeLRU[K, V] {
	hash := s.hasher(k)
	shard := hash & s.masks
	return s.shards[shard]
}

func (s *ShardedLRU[K, V]) Put(k K, v V) {
	cache := s.shardForKey(k)
	cache.Put(k, v)
}

func (s *ShardedLRU[K, V]) Get(k K) (V, bool) {
	cache := s.shardForKey(k)
	return cache.Get(k)
}
