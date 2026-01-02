package datastructer

import "testing"

func TestBinaryTree_Insert(t *testing.T) {
	tree := NewBinaryTree[int]()
	tree.Insert(1)
	tree.Insert(2)
	tree.Insert(3)

	if !tree.Search(1) {
		t.Errorf("Expected 1 to be found")
	}
	if !tree.Search(2) {
		t.Errorf("Expected 2 to be found")
	}
	if !tree.Search(3) {
		t.Errorf("Expected 3 to be found")
	}
	if tree.Search(4) {
		t.Errorf("Expected 4 to not be found")
	}
}
