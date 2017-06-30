package graph

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

// Id is an unique identifier
type Id string

// Node describes the methods of node operations.
type Node interface {

	// String describes the Graph.
	String() string

	// Id returns the node identifier.
	Id() Id
}

type node struct {
	id string
}

func NewNode(id string) Node {
	return &node{
		id: id,
	}
}

func (nd *node) String() string {
	return nd.id
}

func (nd *node) Id() Id {
	return Id(nd.id)
}

// Edge describes the methods of edge operations.
type Edge interface {

	// String describes the Graph.
	String() string

	// GetWeight returns the weight of the edge.
	GetWeight() interface{}

	// SetWeight sets the weight of the edge.
	SetWeight(wgt interface{})
}

type edge struct {
	weight float64
}

func NewEdge(wgt float64) Edge {
	return &edge{
		weight: wgt,
	}
}

func (eg *edge) String() string {
	return fmt.Sprintf("%.2f", eg.weight)
}

func (eg *edge) GetWeight() interface{} {
	return eg.weight
}
func (eg *edge) SetWeight(wgt interface{}) {
	eg.weight = wgt.(float64)
}

// Graph describes the methods of graph operations.
type Graph interface {

	// String describes the Graph.
	String() string

	// Nodes returns a map of nodes.
	Nodes() map[Id]Node

	// Sources returns a map of sources.
	SourcesMap() map[Id]map[Id]Edge

	// Targets returns a map of targets.
	TargetsMap() map[Id]map[Id]Edge

	// GetSources returns a map of parent nodes.
	GetSources(id Id) (map[Id]Node, error)

	// GetTargets returns a map of child nodes.
	GetTargets(id Id) (map[Id]Node, error)

	// GetNode finds and returns the node.
	GetNode(id Id) (Node, error)

	// AddNode adds a node to a graph.
	AddNode(nd Node) error

	// DeleteNode deletes a node from a graph.
	DeleteNode(id Id) error

	// GetEdge finds the edge.
	GetEdge(src, tgt Id) (Edge, error)

	// AddEdge adds an edge from src to tgt
	AddEdge(src, tgt Id, eg Edge) error

	// DeleteEdge deletes an edge from src to tgt.
	DeleteEdge(src, tgt Id) error
}

var (
	ErrNodeNotFound = errors.New("Node not found.")
	ErrNodeExisted  = errors.New("Node existed.")
	ErrEdgeNotFound = errors.New("Edge not found.")
	ErrEdgeExisted  = errors.New("Edge existed.")
)

type graph struct {
	sync.RWMutex

	// nodes stores all nodes.
	nodes map[Id]Node

	// sourcesMap maps a node to parents with edge.
	sourcesMap map[Id]map[Id]Edge

	// targetsMap maps a node to children with edge.
	targetsMap map[Id]map[Id]Edge
}

// newGraph returns a new graph.
func newGraph() *graph {
	return &graph{
		nodes:      make(map[Id]Node),
		sourcesMap: make(map[Id]map[Id]Edge),
		targetsMap: make(map[Id]map[Id]Edge),
	}
}

// NewGraph returns a new graph.
func NewGraph() Graph {
	return newGraph()
}

func (g *graph) String() string {
	g.RWMutex.RLock()
	defer g.RWMutex.RUnlock()

	buf := new(bytes.Buffer)
	for id, nd := range g.nodes {
		tgts, _ := g.GetTargets(id)
		for id2, nd2 := range tgts {
			wgt, _ := g.GetEdge(id, id2)
			fmt.Fprintf(buf, "%s -- %s --> %s\n", nd, wgt, nd2)
		}
	}
	return buf.String()
}

func (g *graph) Nodes() map[Id]Node {
	g.RWMutex.RLock()
	defer g.RWMutex.RUnlock()

	return g.nodes
}

func (g *graph) SourcesMap() map[Id]map[Id]Edge {
	g.RWMutex.RLock()
	defer g.RWMutex.RUnlock()

	return g.sourcesMap
}

func (g *graph) TargetsMap() map[Id]map[Id]Edge {
	g.RWMutex.RLock()
	defer g.RWMutex.RUnlock()

	return g.targetsMap
}

func (g *graph) GetSources(id Id) (map[Id]Node, error) {
	g.RWMutex.RLock()
	defer g.RWMutex.RUnlock()

	if _, ok := g.nodes[id]; !ok {
		return nil, ErrNodeNotFound
	}

	srcs := make(map[Id]Node)
	if _, ok := g.sourcesMap[id]; ok {
		for id2, _ := range g.sourcesMap[id] {
			srcs[id2] = g.nodes[id2]
		}
	}
	return srcs, nil
}

func (g *graph) GetTargets(id Id) (map[Id]Node, error) {
	g.RWMutex.RLock()
	defer g.RWMutex.RUnlock()

	if _, ok := g.nodes[id]; !ok {
		return nil, ErrNodeNotFound
	}

	tgts := make(map[Id]Node)
	if _, ok := g.targetsMap[id]; ok {
		for id2 := range g.targetsMap[id] {
			tgts[id2] = g.nodes[id2]
		}
	}
	return tgts, nil
}

func (g *graph) GetNode(id Id) (Node, error) {
	g.RWMutex.RLock()
	defer g.RWMutex.RUnlock()

	if _, ok := g.nodes[id]; !ok {
		return nil, ErrNodeNotFound
	}

	return g.nodes[id], nil
}

func (g *graph) AddNode(node Node) error {
	g.RWMutex.Lock()
	defer g.RWMutex.Unlock()

	id := node.Id()
	if _, ok := g.nodes[id]; ok {
		return ErrNodeExisted
	}

	g.nodes[id] = node
	return nil
}

func (g *graph) DeleteNode(id Id) error {
	g.RWMutex.Lock()
	defer g.RWMutex.Unlock()

	if _, ok := g.nodes[id]; !ok {
		return ErrNodeNotFound
	}

	delete(g.nodes, id)
	delete(g.sourcesMap, id)
	delete(g.targetsMap, id)

	// remove edges which source is the node
	for _, srcs := range g.sourcesMap {
		delete(srcs, id)
	}
	// remove edges which target is the node
	for _, tgts := range g.targetsMap {
		delete(tgts, id)
	}

	return nil
}

func (g *graph) GetEdge(src, tgt Id) (Edge, error) {
	g.RWMutex.RLock()
	defer g.RWMutex.RUnlock()

	if _, ok := g.nodes[src]; !ok {
		return nil, ErrNodeNotFound
	}
	if _, ok := g.nodes[tgt]; !ok {
		return nil, ErrNodeNotFound
	}

	if _, ok := g.targetsMap[src]; ok {
		// find the edge from src to tgt
		if edge, ok2 := g.targetsMap[src][tgt]; ok2 {
			return edge, nil
		}
	}

	return nil, ErrEdgeNotFound
}

func (g *graph) AddEdge(src, tgt Id, eg Edge) error {
	g.RWMutex.Lock()
	defer g.RWMutex.Unlock()

	if _, ok := g.nodes[src]; !ok {
		return ErrNodeNotFound
	}
	if _, ok := g.nodes[tgt]; !ok {
		return ErrNodeNotFound
	}

	if _, ok := g.sourcesMap[tgt]; !ok {
		g.sourcesMap[tgt] = make(map[Id]Edge)
	}
	if _, ok := g.targetsMap[src]; !ok {
		g.targetsMap[src] = make(map[Id]Edge)
	}

	if _, ok := g.sourcesMap[tgt][src]; ok {
		return ErrEdgeExisted
	}
	if _, ok := g.targetsMap[src][tgt]; ok {
		return ErrEdgeExisted
	}

	g.sourcesMap[tgt][src] = eg
	g.targetsMap[src][tgt] = eg

	return nil
}

func (g *graph) DeleteEdge(src, tgt Id) error {
	g.RWMutex.Lock()
	defer g.RWMutex.Unlock()

	if _, ok := g.nodes[src]; !ok {
		return ErrNodeNotFound
	}
	if _, ok := g.nodes[tgt]; !ok {
		return ErrNodeNotFound
	}

	if _, ok := g.sourcesMap[tgt]; ok {
		if _, ok2 := g.sourcesMap[tgt][src]; ok2 {
			delete(g.sourcesMap[tgt], src)
		}
	}

	if _, ok := g.targetsMap[src]; ok {
		if _, ok2 := g.targetsMap[src][tgt]; ok2 {
			delete(g.targetsMap[src], tgt)
		}
	}

	return nil
}
