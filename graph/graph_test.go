package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

var js map[string]map[string]map[string]interface{}

func setup() error {
	file, err := os.Open("../test/graph.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	js = make(map[string]map[string]map[string]interface{})
	for {
		err := decoder.Decode(&js)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func teardown() {
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		log.Println(err)
		return
	}
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func loadGraph(id string) (Graph, error) {
	if _, ok := js[id]; !ok {
		return nil, fmt.Errorf("%s does not exist", id)
	}
	jsGraph := js[id]

	g := NewGraph()
	for id, neighbour := range jsGraph {
		nd, err := g.GetNode(Id(id))
		if err != nil {
			nd = NewNode(id)
			g.AddNode(nd)
		}
		for id2, weight := range neighbour {
			nd2, err := g.GetNode(Id(id2))
			if err != nil {
				nd2 = NewNode(id2)
				g.AddNode(nd2)
			}
			edge := NewEdge(weight.(float64))
			g.AddEdge(nd.Id(), nd2.Id(), edge)
		}
	}
	return g, nil
}

func Test_Graph(t *testing.T) {
	g, err := loadGraph("graph_00")
	if err != nil {
		t.Fatal(err)
	}

	// test modify weight of edge
	eg, err := g.GetEdge("S", "A")
	eg.SetWeight(222.22)

	eg2, err := g.GetEdge("S", "A")
	if eg2.GetWeight() != 222.22 {
		t.Fatal("Modify edge S --> A failed")
	}
	// test re-add edge
	err = g.AddEdge("S", "A", NewEdge(1.11))
	if err != ErrEdgeExisted {
		t.Fatal("Should catch ErrEdgeExisted error")
	}

	// print graph
	fmt.Println(g)
}
