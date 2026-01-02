package algorithm

import (
	"testing"
)

func TestBinarySearch(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	index := BinarySearch(arr, 5)
	if index != 4 {
		t.Fatalf("should be 4")
	}
	index = BinarySearch(arr, 11)
	if index != -1 {
		t.Fatalf("should be -1, %d", index)
	}
	index = BinarySearch(arr, 1)
	if index != 0 {
		t.Fatalf("should be 0, %d", index)
	}
	index = BinarySearch(arr, 10)
	if index != 9 {
		t.Fatalf("should be 9, %d", index)
	}
}
