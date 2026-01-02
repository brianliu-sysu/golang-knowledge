package algorithm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGraphDFS(t *testing.T) {
	g := NewGraph(5)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	paths := g.DFSPaths(0)
	fmt.Println(paths)
	if len(paths) != 5 {
		t.Fatalf("should be 5")
	}
	if !reflect.DeepEqual(paths[0], []int{0}) {
		t.Fatalf("should be [0]")
	}
	if !reflect.DeepEqual(paths[1], []int{0, 1}) {
		t.Fatalf("should be [0, 1]")
	}
	if !reflect.DeepEqual(paths[2], []int{0, 1, 3}) {
		t.Fatalf("should be [0, 1, 3]")
	}
	if !reflect.DeepEqual(paths[3], []int{0, 2}) {
		t.Fatalf("should be [0, 2]")
	}
	if !reflect.DeepEqual(paths[4], []int{0, 2, 4}) {
		t.Fatalf("should be [0, 2, 4]")
	}
}

func TestGraphBFS(t *testing.T) {
	g := NewGraph(5)
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	paths := g.BFS(0)
	fmt.Println(paths)
	if len(paths) != 5 {
		t.Fatalf("should be 5")
	}
	if !reflect.DeepEqual(paths[0], []int{0}) {
		t.Fatalf("should be [0]")
	}
	if !reflect.DeepEqual(paths[1], []int{0, 1}) {
		t.Fatalf("should be [0, 1]")
	}
	if !reflect.DeepEqual(paths[2], []int{0, 2}) {
		t.Fatalf("should be [0, 2]")
	}
	if !reflect.DeepEqual(paths[3], []int{0, 1, 3}) {
		t.Fatalf("should be [0, 1, 3]")
	}
	if !reflect.DeepEqual(paths[4], []int{0, 2, 4}) {
		t.Fatalf("should be [0, 2, 4]")
	}
}
