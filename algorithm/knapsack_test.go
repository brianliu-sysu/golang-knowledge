package algorithm

import (
	"reflect"
	"testing"
)

func TestKnapsack01(t *testing.T) {
	weights := []int{2, 3, 4, 5}
	values := []int{3, 4, 5, 6}
	capacity := 5
	result := knapsack01(weights, values, capacity)
	if result != 7 {
		t.Fatalf("should be 7")
	}

	capacity = 10
	result = knapsack01(weights, values, capacity)
	if result != 13 {
		t.Fatalf("should be 10")
	}
}

func TestKnapsack01_1D(t *testing.T) {
	weights := []int{2, 3, 4, 5}
	values := []int{3, 4, 5, 6}
	capacity := 5
	result, maxValue := knapsack01_1D(weights, values, capacity)
	if maxValue != 7 {
		t.Fatalf("should be 7")
	}
	if !reflect.DeepEqual(result, []int{0, 0, 3, 4, 5, 7}) {
		t.Fatalf("should be [0, 0, 3, 4, 5, 7], %v", result)
	}
}
