package traversal

// https://www.youtube.com/watch?v=bIA8HEEUxZI&t

import (
	"github.com/jiangyang5157/go-graph/graph"
	"github.com/jiangyang5157/golang-start/data/queue"
	"github.com/jiangyang5157/golang-start/data/stack"
)

// Node describes the methods of node operations.
type Node interface {
	graph.Node
	Data() interface{}
}

type node struct {
	id   string
	data interface{}
}

func (nd *node) String() string {
	return nd.id
}

func (nd *node) Id() graph.Id {
	return graph.Id(nd.id)
}

func (nd *node) Data() interface{} {
	return nd.data
}

func NewNode(id string, data interface{}) Node {
	return &node{
		id:   id,
		data: data,
	}
}

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
