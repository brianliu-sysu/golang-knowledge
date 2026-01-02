package datastructer

import "unicode/utf8"

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode), isEnd: false}
}

func (n *TrieNode) Insert(word string) {
	current := n
	for len(word) > 0 {
		r, size := utf8.DecodeRuneInString(word)
		word = word[size:]
		if _, ok := current.children[r]; !ok {
			current.children[r] = NewTrieNode()
		}
		current = current.children[r]
	}
	current.isEnd = true
}

func (n *TrieNode) Search(word string) bool {
	current := n
	for len(word) > 0 {
		r, size := utf8.DecodeRuneInString(word)
		word = word[size:]
		if _, ok := current.children[r]; !ok {
			return false
		}
		current = current.children[r]
	}
	return current.isEnd
}

type TrieNodeRouter struct {
	node *TrieNode
	rune rune
}

func (n *TrieNode) Delete(word string) bool {
	current := n
	path := make([]TrieNodeRouter, 0)
	for len(word) > 0 {
		r, size := utf8.DecodeRuneInString(word)
		word = word[size:]
		if _, ok := current.children[r]; !ok {
			return false
		}
		path = append(path, TrieNodeRouter{node: current, rune: r})
		current = current.children[r]
	}

	if current.isEnd {
		current.isEnd = false
		for i := len(path) - 1; i >= 0; i-- {
			child := path[i].node.children[path[i].rune]
			if !child.isEnd && len(child.children) == 0 {
				delete(path[i].node.children, path[i].rune)
				continue
			}
			break
		}
		return true
	}

	return false
}

func (n *TrieNode) IsEmpty() bool {
	return len(n.children) == 0
}
