package datastructer

import "cmp"

type BinaryTreeNode[T cmp.Ordered] struct {
	Value T
	Left  *BinaryTreeNode[T]
	Right *BinaryTreeNode[T]
}

type BinaryTree[T cmp.Ordered] struct {
	Root *BinaryTreeNode[T]
}

func NewBinaryTree[T cmp.Ordered]() *BinaryTree[T] {
	return &BinaryTree[T]{Root: nil}
}

func (b *BinaryTree[T]) Insert(value T) {
	if b.Root == nil {
		b.Root = &BinaryTreeNode[T]{Value: value}
		return
	}

	current := b.Root
	for current != nil {
		if value < current.Value {
			if current.Left == nil {
				current.Left = &BinaryTreeNode[T]{Value: value}
				return
			}
			current = current.Left
		} else {
			if current.Right == nil {
				current.Right = &BinaryTreeNode[T]{Value: value}
				return
			}
			current = current.Right
		}
	}
}

func (b *BinaryTree[T]) Search(value T) bool {
	current := b.Root
	for current != nil {
		if value < current.Value {
			current = current.Left
		} else if value > current.Value {
			current = current.Right
		} else {
			return true
		}
	}
	return false
}
