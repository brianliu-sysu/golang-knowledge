package datastructer

import (
	"fmt"
	"hash/fnv"
)

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

type HashTable[K comparable, V any] struct {
	size    int
	buckets [][]Entry[K, V]
	count   int
	hasher  func(K) int
}

func NewHashTable[K comparable, V any](size int) *HashTable[K, V] {
	ht := &HashTable[K, V]{
		size:    size,
		buckets: make([][]Entry[K, V], size),
	}

	ht.hasher = ht.defaultHasher
	return ht
}

func (h *HashTable[K, V]) defaultHasher(key K) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(fmt.Sprintf("%v", key)))
	return int(hasher.Sum32() % uint32(h.size))
}

func (h *HashTable[K, V]) Insert(key K, value V) {
	index := h.hasher(key)

	// check if key already exists
	for _, entry := range h.buckets[index] {
		if entry.Key == key {
			entry.Value = value
			return
		}
	}

	// create new entry and add to bucket
	h.buckets[index] = append(h.buckets[index], Entry[K, V]{Key: key, Value: value})
	h.count++
}

func (h *HashTable[K, V]) Get(key K) (V, bool) {
	index := h.hasher(key)

	// check if key exists
	for _, entry := range h.buckets[index] {
		if entry.Key == key {
			return entry.Value, true
		}
	}

	return *new(V), false
}

func (h *HashTable[K, V]) Delete(key K) bool {
	index := h.hasher(key)

	// check if key exists
	for i, entry := range h.buckets[index] {
		if entry.Key == key {
			h.buckets[index] = append(h.buckets[index][:i], h.buckets[index][i+1:]...)
			h.count--
			return true
		}
	}

	return false
}

func (h *HashTable[K, V]) Size() int {
	return h.count
}

func (h *HashTable[K, V]) Clear() {
	h.buckets = make([][]Entry[K, V], h.size)
	h.count = 0
}

func (h *HashTable[K, V]) Keys() []K {
	keys := make([]K, 0, h.count)
	for _, bucket := range h.buckets {
		for _, entry := range bucket {
			keys = append(keys, entry.Key)
		}
	}
	return keys
}

func (h *HashTable[K, V]) Values() []V {
	values := make([]V, 0, h.count)
	for _, bucket := range h.buckets {
		for _, entry := range bucket {
			values = append(values, entry.Value)
		}
	}
	return values
}
