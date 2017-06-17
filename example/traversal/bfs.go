package traversal

import (
	"github.com/jiangyang5157/go-graph/graph"
	"github.com/jiangyang5157/golang-start/data/queue"
)

func Bfs(g graph.Graph, id graph.Id, f func(Node) bool) (map[graph.Id]Node, error) {
	// visited holds a map of visited node
	visited := make(map[graph.Id]Node)
	nd, err := g.GetNode(id)
	if err != nil {
		return visited, graph.ErrNodeNotFound
	}

	// Visite the begin node
	visited[id] = nd.(Node)
	if f(nd.(Node)) {
		return visited, nil
	}

	targets := g.Targets()
	tmpId := id
	tmpNode := nd
	tmpQueue := queue.NewQueue()
	for {
		if _, ok := targets[tmpId]; ok {
			for id, _ := range targets[tmpId] {
				if _, ok := visited[id]; ok {
					continue
				}
				tmpId = id
				tmpNode, _ = g.GetNode(tmpId)
				visited[tmpId] = tmpNode.(Node)
				if f(tmpNode.(Node)) {
					return visited, nil
				}
				tmpQueue.Push(tmpId)
			}
		}

		if tmpQueue.IsEmpty() {
			break
		} else {
			tmpId = tmpQueue.Pop().(graph.Id)
		}
	}

	return visited, nil
}
