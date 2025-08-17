package libmonpos

import (
	"fmt"

	"github.com/dominikbraun/graph"
)

// Get the monitor order of the monitor graph.
func get_monitor_order(g MonitorGraph) ([]string, error) {
	// FUNCTION: Given a list of vertices in topological order,
	// return whether a graph is disconnected
	var graph_disconnected = func (g MonitorGraph, ordered_vertices []string) (bool) {
		// use a BFS to count the vertices connected to vertices[0]
		count := 0
		graph.BFS(g, ordered_vertices[0], func (string) bool {
			count++
			return false
		})

		// compare to the total number of vertices to see if we missed anything
		return count != len(ordered_vertices)
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
