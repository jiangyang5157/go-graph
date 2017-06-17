package traversal

// https://www.youtube.com/watch?v=bIA8HEEUxZI&t

import "github.com/jiangyang5157/go-graph/graph"

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
