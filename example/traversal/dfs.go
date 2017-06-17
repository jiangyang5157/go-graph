package traversal

import (
	"github.com/jiangyang5157/go-graph/graph"
	"github.com/jiangyang5157/golang-start/data/stack"
)

func dfs(g graph.Graph, visited map[graph.Id]Node, tmpStack *stack.Stack, f func(Node) bool) {
	if tmpStack.IsEmpty() {
		return
	}

	tmpId := tmpStack.Peek().(graph.Id)
	tmpNode, _ := g.GetNode(tmpId)
	visited[tmpId] = tmpNode.(Node)
	if f(tmpNode.(Node)) {
		tmpStack.Pop()
		return
	}

	targets := g.Targets()
	if _, ok := targets[tmpId]; !ok {
		tmpStack.Pop()
		return
	}

	for id, _ := range targets[tmpId] {
		if _, ok := visited[id]; ok {
			continue
		}
		tmpStack.Push(id)
		dfs(g, visited, tmpStack, f)
	}
	tmpStack.Pop()
	return
}

func Dfs(g graph.Graph, id graph.Id, f func(Node) bool) (map[graph.Id]Node, error) {
	// visited holds a map of visited node
	visited := make(map[graph.Id]Node)
	_, err := g.GetNode(id)
	if err != nil {
		return visited, graph.ErrNodeNotFound
	}

	tmpStack := stack.NewStack()
	tmpStack.Push(id)
	dfs(g, visited, tmpStack, f)
	return visited, nil
}
