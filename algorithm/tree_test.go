package algorithm

import (
	"reflect"
	"testing"
)

func TestTreeNode(t *testing.T) {
	tree := NewTreeNode(5)
	tree.Add(3)
	tree.Add(7)
	tree.Add(2)
	tree.Add(4)
	tree.Add(6)

	if !reflect.DeepEqual(tree.InOrderTraversal(), []int{2, 3, 4, 5, 6, 7}) {
		t.Fatalf("should be equal")
	}
	if !reflect.DeepEqual(tree.PreOrderTraversal(), []int{5, 3, 2, 4, 7, 6}) {
		t.Fatalf("should be equal")
	}
	if !reflect.DeepEqual(tree.PostOrderTraversal(), []int{2, 4, 3, 6, 7, 5}) {
		t.Fatalf("should be equal")
	}
}

func TestTreeDelete(t *testing.T) {
	tree := NewTree()
	tree.Add(5)
	tree.Add(3)
	tree.Add(7)
	tree.Add(2)
	tree.Add(4)
	tree.Add(6)
	tree.Delete(3)
	if !reflect.DeepEqual(tree.InOrderTraversal(), []int{2, 4, 5, 6, 7}) {
		t.Fatalf("should be equal")
	}

	tree.Delete(5)
	if !reflect.DeepEqual(tree.InOrderTraversal(), []int{2, 4, 6, 7}) {
		t.Fatalf("should be equal")
	}

	tree.Delete(7)
	if !reflect.DeepEqual(tree.InOrderTraversal(), []int{2, 4, 6}) {
		t.Fatalf("should be equal")
	}

	tree.Delete(2)
	if !reflect.DeepEqual(tree.InOrderTraversal(), []int{4, 6}) {
		t.Fatalf("should be equal")
	}

	tree.Delete(4)
	if !reflect.DeepEqual(tree.InOrderTraversal(), []int{6}) {
		t.Fatalf("should be equal")
	}
	tree.Delete(6)
	if !reflect.DeepEqual(tree.InOrderTraversal(), []int{}) {
		t.Fatalf("should be equal")
	}
}
