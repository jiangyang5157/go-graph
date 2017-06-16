package bfs

import "github.com/jiangyang5157/go-graph/graph"

func Bfs(g graph.Graph, id graph.Id, f func(graph.Node)) error {
	nd, err := g.GetNode(id)
	if err != nil {
		return graph.ErrNodeNotFound
	}

	f(nd)

	return nil
}
