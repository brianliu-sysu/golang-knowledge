package datastructer

import "testing"

func TestMaxHeap(t *testing.T) {
	heap := NewMaxHeap[int]()
	heap.Push(1)
	heap.Push(2)
	heap.Push(3)
	item, ok := heap.Pop()
	if !ok {
		t.Errorf("Expected pop value 3, got %d", item)
	}
	if item != 3 {
		t.Errorf("Expected pop value 3, got %d", item)
	}
	item, ok = heap.Pop()
	if !ok {
		t.Errorf("Expected pop value 2, got %d", item)
	}
	if item != 2 {
		t.Errorf("Expected pop value 2, got %d", item)
	}
	item, ok = heap.Pop()
	if !ok {
		t.Errorf("Expected pop value 1, got %d", item)
	}
	if item != 1 {
		t.Errorf("Expected pop value 1, got %d", item)
	}
}

func TestMinHeap(t *testing.T) {
	heap := NewMinHeap[int]()
	heap.Push(1)
	heap.Push(2)
	heap.Push(3)
	item, ok := heap.Pop()
	if !ok {
		t.Errorf("Expected pop value 1, got %d", item)
	}
	if item != 1 {
		t.Errorf("Expected pop value 1, got %d", item)
	}
	item, ok = heap.Pop()
	if !ok {
		t.Errorf("Expected pop value 2, got %d", item)
	}
	if item != 2 {
		t.Errorf("Expected pop value 2, got %d", item)
	}
	item, ok = heap.Pop()
	if !ok {
		t.Errorf("Expected pop value 3, got %d", item)
	}
	if item != 3 {
		t.Errorf("Expected pop value 3, got %d", item)
	}
}
