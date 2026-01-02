package algorithm

type Graph struct {
	Vertices int
	Edges    int
	AdjList  [][]int
}

func NewGraph(vertices int) *Graph {
	return &Graph{Vertices: vertices, Edges: 0, AdjList: make([][]int, vertices)}
}

func (g *Graph) AddEdge(u, v int) {
	g.AdjList[u] = append(g.AdjList[u], v)
	g.Edges++
}

// Depth First Search, return the path of the graph
func (g *Graph) DFSPaths(start int) [][]int {
	paths := [][]int{}
	visited := make([]bool, g.Vertices)

	var dfs func(v int, path []int)
	dfs = func(v int, path []int) {
		visited[v] = true
		paths = append(paths, append([]int{}, path...))
		for _, u := range g.AdjList[v] {
			if !visited[u] {
				dfs(u, append(path, u))
			}
		}

		visited[v] = false
	}

	dfs(start, []int{start})
	return paths
}

func (g *Graph) BFS(start int) [][]int {
	paths := [][]int{}
	visited := make([]bool, g.Vertices)
	visited[start] = true
	queues := [][]int{{start}}

	for len(queues) > 0 {
		path := queues[0]
		queues = queues[1:]

		v := path[len(path)-1]
		paths = append(paths, path)
		for _, u := range g.AdjList[v] {
			if !visited[u] {
				visited[u] = true
				newPath := append([]int{}, path...)
				newPath = append(newPath, u)
				queues = append(queues, newPath)
			}
		}
	}
	return paths
}
