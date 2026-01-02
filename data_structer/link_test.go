package datastructer

import "testing"

func TestLinkedList_Add(t *testing.T) {
	ll := NewLinkedList[int]()
	ll.Add(1)
	ll.Add(2)
	ll.Add(3)

	if ll.Size != 3 {
		t.Errorf("Expected size 3, got %d", ll.Size)
	}

	node := ll.Head
	for i := range 3 {
		if node.Value != i+1 {
			t.Errorf("Expected value %d, got %d", i+1, node.Value)
		}
		node = node.Next
	}

	node = ll.Last
	for i := range 3 {
		if node.Value != 3-i {
			t.Errorf("Expected value %d, got %d", 3-i, node.Value)
		}
		node = node.Prev
	}
}

func TestLinkedList_Remove(t *testing.T) {
	ll := NewLinkedList[int]()
	ll.Add(1)
	ll.Add(2)
	ll.Add(3)

	ll.Remove(2)
	if ll.Size != 2 {
		t.Errorf("Expected size 2, got %d", ll.Size)
	}
	if ll.Head.Value != 1 {
		t.Errorf("Expected head value 1, got %d", ll.Head.Value)
	}
	if ll.Last.Value != 3 {
		t.Errorf("Expected last value 3, got %d", ll.Last.Value)
	}

	ll.Remove(1)
	if ll.Size != 1 {
		t.Errorf("Expected size 1, got %d", ll.Size)
	}
	if ll.Head.Value != 3 {
		t.Errorf("Expected head value 3, got %d", ll.Head.Value)
	}
	if ll.Last.Value != 3 {
		t.Errorf("Expected last value 3, got %d", ll.Last.Value)
	}

	ll.Remove(3)
	if ll.Size != 0 {
		t.Errorf("Expected size 0, got %d", ll.Size)
	}
	if ll.Head != nil {
		t.Errorf("Expected head to be nil, got %v", ll.Head)
	}
	if ll.Last != nil {
		t.Errorf("Expected last to be nil, got %v", ll.Last)
	}
}
