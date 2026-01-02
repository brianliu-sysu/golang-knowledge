package algorithm

import (
	"reflect"
	"testing"
)

func TestQuickSort(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	QuickSort(arr)

	if !reflect.DeepEqual(arr, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Fatalf("should be equal")
	}

	arr = []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	QuickSort(arr)

	if !reflect.DeepEqual(arr, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Fatalf("should be equal")
	}

	arr = []int{3, 5, 4, 2, 1, 6, 7, 8, 9, 10}
	QuickSort(arr)

	if !reflect.DeepEqual(arr, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Fatalf("should be equal")
	}
}
