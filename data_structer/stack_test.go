package datastructer

import "testing"

func TestStack_Push(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	item, ok := stack.Peek()
	if !ok {
		t.Errorf("Expected peek value 3, got %d", item)
	}
	if item != 3 {
		t.Errorf("Expected peek value 3, got %d", item)
	}
	item, ok = stack.Pop()
	if !ok {
		t.Errorf("Expected pop value 3, got %d", item)
	}
	if item != 3 {
		t.Errorf("Expected pop value 3, got %d", item)
	}
	item, ok = stack.Pop()
	if !ok {
		t.Errorf("Expected pop value 2, got %d", item)
	}
	if item != 2 {
		t.Errorf("Expected pop value 2, got %d", item)
	}
	item, ok = stack.Pop()
	if !ok {
		t.Errorf("Expected pop value 1, got %d", item)
	}
	if item != 1 {
		t.Errorf("Expected pop value 1, got %d", item)
	}
	if !stack.IsEmpty() {
		t.Errorf("Expected stack to be empty, got %v", stack)
	}
}
