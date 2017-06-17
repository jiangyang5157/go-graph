package dijkstra

import (
	"errors"
	"math"

	"github.com/jiangyang5157/go-graph/graph"
)

// Dijkstra's algorithm (Greedy algorithm) is an algorithm for finding the shortest paths from one node to every other node.
// Dijkstra Returns a map of distances from one node to every other node.
func Dijkstra(g graph.Graph, id graph.Id) (map[graph.Id]float64, error) {
	_, err := g.GetNode(id)
	if err != nil {
		return nil, graph.ErrNodeNotFound
	}

	// ret holds a map of visited node
	ret := make(map[graph.Id]float64)
	// dists holds a map of unvisited node
	dists := make(map[graph.Id]float64)

	// initalize distances to maximum
	for id, _ := range g.Nodes() {
		dists[id] = math.MaxFloat64
	}
	// distance from source to source is zero
	dists[id] = 0

	targets := g.Targets()
	var tmpId graph.Id
	for {
		// result completed
		if len(dists) == 0 {
			break
		}

		// find Id with minimum distance
		min := math.MaxFloat64
		for id, dist := range dists {
			if dist < min {
				min = dist
				tmpId = id
			}
		}

		if min == math.MaxFloat64 {
			return nil, errors.New("Error: Infinite distance.")
		}

		// remove id-dist that has minimum distance, and add it into result
		delete(dists, tmpId)
		ret[tmpId] = min

		if _, ok := targets[tmpId]; !ok {
			// no targets, path ends
			break
		}

		// re-calculate distances
		for id, edge := range targets[tmpId] {
			if _, ok := ret[id]; ok {
				// ignore visited node
				continue
			}
			d := ret[tmpId] + edge.GetWeight().(float64)
			if d < dists[id] {
				dists[id] = d
			}
		}
	}

	return ret, nil
}
