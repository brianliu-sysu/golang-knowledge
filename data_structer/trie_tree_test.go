package datastructer

import "testing"

func TestTrieTreeSearch(t *testing.T) {
	trie := NewTrieNode()
	trie.Insert("hello")
	trie.Insert("world")
	trie.Insert("hell")
	trie.Insert("world")
	if !trie.Search("hello") {
		t.Fatalf("should be true")
	}
	if !trie.Search("world") {
		t.Fatalf("should be true")
	}
	if !trie.Search("hell") {
		t.Fatalf("should be true")
	}
	if !trie.Search("world") {
		t.Fatalf("should be true")
	}
}

func TestTrieTreeDelete(t *testing.T) {
	trie := NewTrieNode()
	trie.Insert("hello")
	trie.Insert("world")
	trie.Insert("hell")
	trie.Insert("world")
	if !trie.Delete("hello") {
		t.Fatalf("should be true")
	}
	if trie.Search("hello") {
		t.Fatalf("should be false")
	}
	if !trie.Search("world") {
		t.Fatalf("should be true")
	}
	if !trie.Search("hell") {
		t.Fatalf("should be true")
	}
	if !trie.Search("world") {
		t.Fatalf("should be true")
	}
	if !trie.Delete("world") {
		t.Fatalf("should be true")
	}
	if trie.Delete("world") {
		t.Fatalf("should be false")
	}
}

func TestTrieTreeIsEmpty(t *testing.T) {
	trie := NewTrieNode()
	if !trie.IsEmpty() {
		t.Fatalf("should be true")
	}
	trie.Insert("hello")
	if trie.IsEmpty() {
		t.Fatalf("should be false")
	}
}
