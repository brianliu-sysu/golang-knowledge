package datastructer

import "cmp"

type MaxHeap[T cmp.Ordered] struct {
	items []T
}

func NewMaxHeap[T cmp.Ordered]() *MaxHeap[T] {
	return &MaxHeap[T]{items: make([]T, 0)}
}

func (h *MaxHeap[T]) Push(item T) {
	h.items = append(h.items, item)
	h.heapifyUp()
}

func (h *MaxHeap[T]) Pop() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	item := h.items[0]
	h.items = h.items[1:]
	h.heapifyDown()
	return item, true
}

func (h *MaxHeap[T]) heapifyUp() {
	index := len(h.items) - 1
	for {
		parentIndex := (index - 1) / 2
		if index == 0 || cmp.Less(h.items[index], h.items[parentIndex]) {
			break
		}
		h.items[index], h.items[parentIndex] = h.items[parentIndex], h.items[index]
		index = parentIndex
	}
}

func (h *MaxHeap[T]) heapifyDown() {
	index := 0
	n := len(h.items)
	for {
		leftChildIndex := 2*index + 1
		if leftChildIndex >= n {
			break
		}

		bestChildIndex := leftChildIndex
		rightChildIndex := 2*index + 2
		if rightChildIndex < n && cmp.Less(h.items[leftChildIndex], h.items[rightChildIndex]) {
			bestChildIndex = rightChildIndex
		}
		if cmp.Less(h.items[bestChildIndex], h.items[index]) {
			break
		}

		h.items[index], h.items[bestChildIndex] = h.items[bestChildIndex], h.items[index]
		index = bestChildIndex
	}
}

func (h *MaxHeap[T]) Peek() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	return h.items[0], true
}

func (h *MaxHeap[T]) IsEmpty() bool {
	return len(h.items) == 0
}

func (h *MaxHeap[T]) Size() int {
	return len(h.items)
}

type MinHeap[T cmp.Ordered] struct {
	items []T
}

func NewMinHeap[T cmp.Ordered]() *MinHeap[T] {
	return &MinHeap[T]{items: make([]T, 0)}
}

func (h *MinHeap[T]) Push(item T) {
	h.items = append(h.items, item)
	h.heapifyUp()
}

func (h *MinHeap[T]) Pop() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	item := h.items[0]
	h.items = h.items[1:]
	h.heapifyDown()
	return item, true
}

func (h *MinHeap[T]) heapifyUp() {
	index := len(h.items) - 1
	for parentIndex := (index - 1) / 2; index > 0 && cmp.Less(h.items[index], h.items[parentIndex]); index = parentIndex {
		h.items[index], h.items[parentIndex] = h.items[parentIndex], h.items[index]
		parentIndex = (index - 1) / 2
	}
}

func (h *MinHeap[T]) heapifyDown() {
	index := 0
	n := len(h.items)
	for {
		leftChildIndex := 2*index + 1
		if leftChildIndex >= n {
			break
		}

		bestChildIndex := leftChildIndex
		rightChildIndex := 2*index + 2
		if rightChildIndex < n && cmp.Less(h.items[rightChildIndex], h.items[leftChildIndex]) {
			bestChildIndex = rightChildIndex
		}
		if cmp.Less(h.items[index], h.items[bestChildIndex]) {
			break
		}

		h.items[index], h.items[bestChildIndex] = h.items[bestChildIndex], h.items[index]
		index = bestChildIndex
	}
}

func (h *MinHeap[T]) Peek() (T, bool) {
	if len(h.items) == 0 {
		var zero T
		return zero, false
	}
	return h.items[0], true
}

func (h *MinHeap[T]) IsEmpty() bool {
	return len(h.items) == 0
}

func (h *MinHeap[T]) Size() int {
	return len(h.items)
}
