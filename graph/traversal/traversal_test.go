package traversal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/jiangyang5157/go-graph/graph"
)

var js map[string]map[string]map[string]interface{}

func setup() error {
	file, err := os.Open("../../testdata/graph.json")
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

func loadGraph(id string) (graph.Graph, error) {
	if _, ok := js[id]; !ok {
		return nil, fmt.Errorf("%s does not exist", id)
	}
	jsGraph := js[id]

	g := graph.NewGraph()
	for id, neighbour := range jsGraph {
		nd, err := g.GetNode(graph.Id(id))
		if err != nil {
			nd = graph.NewNode(id)
			g.AddNode(nd)
		}
		for id2, weight := range neighbour {
			nd2, err := g.GetNode(graph.Id(id2))
			if err != nil {
				nd2 = graph.NewNode(id2)
				g.AddNode(nd2)
			}
			edge := graph.NewEdge(weight.(float64))
			g.AddEdge(nd.Id(), nd2.Id(), edge)
		}
	}
	return g, nil
}

func Test_Traversal(t *testing.T) {
	g, err := loadGraph("graph_13")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(g)
	nodes := g.Nodes()

	// BFS
	visited := 0
	err = Bfs(g, "A", func(nd graph.Node) bool {
		fmt.Printf("BFS visite %v\n", nd.String())
		visited++
		return false
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\nGraph has %v nodes in total.\n", len(nodes))
	fmt.Printf("BFS visited %v nodes in total.\n\n", visited)

	// DFS
	visited = 0
	err = Dfs(g, "A", func(nd graph.Node) bool {
		fmt.Printf("DFS visite %v\n", nd.String())
		visited++
		return false
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\nGraph has %v nodes in total.\n", len(nodes))
	fmt.Printf("DFS visited %v nodes in total.\n\n", visited)
}
