package libmonpos

import (
	"fmt"

	"github.com/dominikbraun/graph"
)

func graph_disconnected(g graph.Graph[string,string], order []string) bool {
	// count the vertices attached to order[0] with a bfs
	count := 0
	graph.BFS(g, order[0], func (string) bool {
		count++
		return false
	})

	// the graph is disconnected if not all vertices were visited
	return count != len(order)
}

func find_monitor_order(conf Config) ([]string, error) {
	g := graph.New(graph.StringHash, graph.PreventCycles(), graph.Directed())

	// first pass: add a vertex for each monitor
	for name := range conf.Monitors {
		err := g.AddVertex(name)
		if err != nil {
			return nil, err
		}
	}

	// second pass: edges between nodes that are positioned next to each other
	for name, mon := range conf.Monitors {
		if mon.Position != "" {
			_, parent, err := split_position(mon.Position)
			if err != nil {
				return nil, err
			}

			err = g.AddEdge(parent, name)
			if err != nil {
				return nil, err
			}
		}
	}

	// do a topological sort of the graph.
	// this makes sure every vertex appears before its child, but DOES NOT have any guarantees about if the graph is connected.
	order, err := graph.TopologicalSort(g)
	if err != nil {
		return nil, err
	}

	// check that all the vertices are connected, if not be mad
	if graph_disconnected(g, order) {
		return nil, fmt.Errorf("all monitors must be connected to one main monitor")
	}

	return order, nil
}
