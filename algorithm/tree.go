package algorithm

type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

func NewTreeNode(value int) *TreeNode {
	return &TreeNode{Value: value}
}

func (t *TreeNode) Add(value int) {
	if value < t.Value {
		if t.Left == nil {
			t.Left = NewTreeNode(value)
		} else {
			t.Left.Add(value)
		}
	} else {
		if t.Right == nil {
			t.Right = NewTreeNode(value)
		} else {
			t.Right.Add(value)
		}
	}

}

func (t *TreeNode) InOrderTraversal() []int {
	if t == nil {
		return []int{}
	}
	result := []int{}
	result = append(result, t.Left.InOrderTraversal()...)
	result = append(result, t.Value)
	result = append(result, t.Right.InOrderTraversal()...)
	return result
}

func (t *TreeNode) PreOrderTraversal() []int {
	if t == nil {
		return []int{}
	}
	result := []int{}
	result = append(result, t.Value)
	result = append(result, t.Left.PreOrderTraversal()...)
	result = append(result, t.Right.PreOrderTraversal()...)
	return result
}

func (t *TreeNode) PostOrderTraversal() []int {
	if t == nil {
		return []int{}
	}
	result := []int{}
	result = append(result, t.Left.PostOrderTraversal()...)
	result = append(result, t.Right.PostOrderTraversal()...)
	result = append(result, t.Value)
	return result
}

func (t *TreeNode) LevelOrderTraversal() [][]int {
	if t == nil {
		return [][]int{}
	}
	result := [][]int{}
	result = append(result, []int{t.Value})
	result = append(result, t.Left.LevelOrderTraversal()...)
	result = append(result, t.Right.LevelOrderTraversal()...)
	return result
}

func (t *TreeNode) Height() int {
	if t == nil {
		return 0
	}
	return 1 + max(t.Left.Height(), t.Right.Height())
}

func (t *TreeNode) Max() int {
	if t == nil {
		return 0
	}
	return max(t.Value, max(t.Left.Max(), t.Right.Max()))
}

func (t *TreeNode) Min() int {
	if t == nil {
		return 0
	}
	return min(t.Value, min(t.Left.Min(), t.Right.Min()))
}

func (t *TreeNode) Search(value int) *TreeNode {
	if t == nil {
		return nil
	}
	if t.Value == value {
		return t
	}
	if value < t.Value {
		return t.Left.Search(value)
	}
	return t.Right.Search(value)
}

func (t *TreeNode) LeftMin() *TreeNode {
	if t == nil {
		return nil
	}
	if t.Left == nil {
		return t
	}
	return t.Left.LeftMin()
}

type Tree struct {
	Root *TreeNode
}

func NewTree() *Tree {
	return &Tree{Root: nil}
}

func (t *Tree) Add(value int) {
	if t.Root == nil {
		t.Root = NewTreeNode(value)
		return
	}
	t.Root.Add(value)
}

func (t *Tree) InOrderTraversal() []int {
	return t.Root.InOrderTraversal()
}

func (t *Tree) PreOrderTraversal() []int {
	return t.Root.PreOrderTraversal()
}

func (t *Tree) PostOrderTraversal() []int {
	return t.Root.PostOrderTraversal()
}

func (t *Tree) LevelOrderTraversal() [][]int {
	return t.Root.LevelOrderTraversal()
}

func (t *Tree) Height() int {
	return t.Root.Height()
}

func (t *Tree) Max() int {
	return t.Root.Max()
}

func (t *Tree) Min() int {
	return t.Root.Min()
}

func (t *Tree) Search(value int) *TreeNode {
	return t.Root.Search(value)
}

func (t *Tree) Delete(value int) bool {
	if t.Root == nil {
		return false
	}

	// find the node to delete and record the parent
	var parent *TreeNode = nil
	current := t.Root
	for current != nil && current.Value != value {
		parent = current
		if value < current.Value {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	if current == nil {
		return false
	}

	// case 1: the node is a leaf node
	if current.Left == nil && current.Right == nil {
		if parent == nil {
			t.Root = nil
		} else if parent.Left == current {
			parent.Left = nil
		} else {
			parent.Right = nil
		}
		return true
	}

	// case 2: the node has only right child
	if current.Left == nil {
		if parent == nil {
			t.Root = current.Right
		} else if parent.Left == current {
			parent.Left = current.Right
		} else {
			parent.Right = current.Right
		}
		return true
	}

	// case 3: the node has only left child
	if current.Right == nil {
		if parent == nil {
			t.Root = current.Left
		} else if parent.Left == current {
			parent.Left = current.Left
		} else {
			parent.Right = current.Left
		}
		return true
	}

	// case 4: the node has both left and right children
	successorParent := current
	successor := current.Right
	for successor.Left != nil {
		successorParent = successor
		successor = successor.Left
	}
	// replace the current node with the successor
	current.Value = successor.Value

	// delete the successor
	if successorParent == current {
		successorParent.Right = successor.Right
	} else {
		successorParent.Left = successor.Right
	}

	return true
}
