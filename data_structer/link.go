package datastructer

type Node[T comparable] struct {
	Value T
	Next  *Node[T]
	Prev  *Node[T]
}

type LinkedList[T comparable] struct {
	Head *Node[T]
	Last *Node[T]
	Size int
}

func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{Head: nil, Last: nil, Size: 0}
}

func (l *LinkedList[T]) Add(value T) {
	node := &Node[T]{Value: value, Next: l.Head, Prev: l.Last}
	if l.Head == nil {
		node.Next = node
		node.Prev = node

		l.Head = node
		l.Last = node

		l.Size++
		return
	}

	l.Last.Next = node
	l.Last = node

	l.Size++
}

func (l *LinkedList[T]) Remove(value T) {
	if l.Head == nil {
		return
	}
	if l.Head.Value == value {
		l.RemoveNode(l.Head)
		return
	}
	current := l.Head
	for current.Next != nil {
		if current.Next.Value == value {
			l.RemoveNode(current.Next)
			return
		}
		current = current.Next
	}
}

func (l *LinkedList[T]) RemoveNode(node *Node[T]) {
	if node == l.Head && node == l.Last {
		l.Head = nil
		l.Last = nil
		l.Size = 0
		return
	}

	if node == l.Head {
		l.Head = node.Next
		l.Head.Prev = l.Last
		l.Last.Next = l.Head
		l.Size--
		return
	}
	if node == l.Last {
		l.Last = node.Prev
		l.Last.Next = l.Head
		l.Head.Prev = l.Last
		l.Size--
		return
	}

	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	l.Size--
}
func (l *LinkedList[T]) Contains(value T) bool {
	return l.Head != nil && l.Head.Value == value
}
