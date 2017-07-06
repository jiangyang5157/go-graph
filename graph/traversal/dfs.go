package traversal

import (
	"github.com/jiangyang5157/go-graph/graph"
	"github.com/jiangyang5157/golang-start/data/stack"
)

func dfs(g graph.Graph, tmpStack *stack.Stack, f func(graph.Node) bool, visited map[graph.Id]bool) {
	if tmpStack.IsEmpty() {
		return
	}

	tmpId := tmpStack.Peek().(graph.Id)
	tmpNode, _ := g.GetNode(tmpId)
	visited[tmpId] = true
	if f(tmpNode) {
		tmpStack.Pop()
		return
	}

	targets := g.TargetsMap()
	if _, ok := targets[tmpId]; !ok {
		tmpStack.Pop()
		return
	}

	for id, _ := range targets[tmpId] {
		if _, ok := visited[id]; ok {
			continue
		}
		tmpStack.Push(id)
		dfs(g, tmpStack, f, visited)
	}
	tmpStack.Pop()
	return
}

func Dfs(g graph.Graph, id graph.Id, f func(graph.Node) bool) error {
	_, err := g.GetNode(id)
	if err != nil {
		return graph.ErrNodeNotFound
	}

	tmpStack := stack.NewStack()
	tmpStack.Push(id)
	dfs(g, tmpStack, f, make(map[graph.Id]bool))
	return nil
}
