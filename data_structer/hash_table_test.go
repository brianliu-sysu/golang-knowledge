package datastructer

import (
	"testing"
)

func TestHashTable(t *testing.T) {
	ht := NewHashTable[int, string](10)
	if ht == nil {
		t.Errorf("NewHashTable returned nil")
	}
	if ht.size != 10 {
		t.Errorf("NewHashTable size is %d, want 10", ht.size)
	}
	if ht.buckets == nil {
		t.Errorf("NewHashTable buckets is nil")
	}
	if ht.count != 0 {
		t.Errorf("NewHashTable count is %d, want 0", ht.count)
	}
	if ht.hasher == nil {
		t.Errorf("NewHashTable hasher is nil")
	}
}

func TestHashTableInsert(t *testing.T) {
	ht := NewHashTable[int, string](10)
	ht.Insert(1, "one")
	ht.Insert(2, "two")
	ht.Insert(3, "three")
	value, ok := ht.Get(1)
	if !ok {
		t.Errorf("NewHashTable Get(1) is not found")
	}
	if value != "one" {
		t.Errorf("NewHashTable Get(1) is %s, want one", value)
	}
	value, ok = ht.Get(2)
	if !ok {
		t.Errorf("NewHashTable Get(2) is not found")
	}
	if value != "two" {
		t.Errorf("NewHashTable Get(2) is %s, want two", value)
	}
	value, ok = ht.Get(3)
	if !ok {
		t.Errorf("NewHashTable Get(3) is not found")
	}
	if value != "three" {
		t.Errorf("NewHashTable Get(3) is %s, want three", value)
	}
}

func TestHashTableDelete(t *testing.T) {
	ht := NewHashTable[int, string](10)
	ht.Insert(1, "one")
	ht.Insert(2, "two")
	ht.Insert(3, "three")
	ht.Delete(2)
	_, ok := ht.Get(2)
	if ok {
		t.Errorf("NewHashTable Get(2) is found")
	}
}
