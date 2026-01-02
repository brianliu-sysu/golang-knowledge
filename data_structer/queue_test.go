package datastructer

import "testing"

func TestQueue_Enqueue(t *testing.T) {
	queue := NewQueue[int]()
	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	item, ok := queue.Peek()
	if !ok {
		t.Errorf("Expected peek value 1, got %d", item)
	}
	if item != 1 {
		t.Errorf("Expected peek value 1, got %d", item)
	}
	item, ok = queue.Dequeue()
	if !ok {
		t.Errorf("Expected dequeue value 1, got %d", item)
	}
	if item != 1 {
		t.Errorf("Expected dequeue value 1, got %d", item)
	}
	item, ok = queue.Dequeue()
	if !ok {
		t.Errorf("Expected dequeue value 2, got %d", item)
	}
	if item != 2 {
		t.Errorf("Expected dequeue value 2, got %d", item)
	}
	item, ok = queue.Dequeue()
	if !ok {
		t.Errorf("Expected dequeue value 3, got %d", item)
	}
	if item != 3 {
		t.Errorf("Expected dequeue value 3, got %d", item)
	}
	if !queue.IsEmpty() {
		t.Errorf("Expected queue to be empty, got %v", queue)
	}
}
