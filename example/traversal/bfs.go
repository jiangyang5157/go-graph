package traversal

import (
	"github.com/jiangyang5157/go-graph/graph"
	"github.com/jiangyang5157/golang-start/data/queue"
)

func Bfs(g graph.Graph, id graph.Id, f func(Node) bool) error {
	nd, err := g.GetNode(id)
	if err != nil {
		return graph.ErrNodeNotFound
	}

	// visited holds a map of visited node Id
	visited := make(map[graph.Id]bool)

	// Visite the begin node
	visited[id] = true
	if f(nd.(Node)) {
		return nil
	}

	targets := g.Targets()
	tmpId := id
	tmpQueue := queue.NewQueue()
	for {
		if _, ok := targets[tmpId]; ok {
			for id, _ := range targets[tmpId] {
				if _, ok := visited[id]; ok {
					continue
				}
				tmpId = id
				tmpNode, _ := g.GetNode(tmpId)
				visited[tmpId] = true
				if f(tmpNode.(Node)) {
					return nil
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

	return nil
}
