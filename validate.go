package libmonpos

import (
	"fmt"
	"slices"

	"github.com/dominikbraun/graph"
)

// Categorization of which directions and alignments are horizontal vs vertical
var h_directions = []string{"left-of", "right-of"}
var v_directions = []string{"above", "below"}
var h_alignments = []string{"left", "right"}
var v_alignments = []string{"top", "bottom"}

// Check the direction and alignment of a monitor
func check_direction_alignment(direction string, alignment string) error {
	// cases where direction is empty
	if direction == "" && alignment == "" {
		return nil // position and align may both be empty
	} else if direction == "" && alignment != "" {
		return fmt.Errorf("position is blank, so alignment must also be blank")
	}

	// Check the general format of position <direction> <monitor>
	// Make sure the direction is valid, we can't check the monitor here
	if !slices.Contains(append(h_directions, v_directions...), direction) {
		return fmt.Errorf("expected direction 'above', 'below', 'left-of', or 'right-of', got '%v'", direction) 
	}

	// depending on if the direction is horizontal, decide whether the alignment is valid or not
	is_horiz := slices.Contains(h_directions, direction)

	// h directions work with v alignments, or "center"
	if is_horiz && !slices.Contains(append(v_alignments, "center"), alignment) {
		return fmt.Errorf("for direction '%v', only alignments 'top', 'bottom', and 'center' are valid. got '%v'", direction, alignment)
	}
	// v directions work with h alignments, or "center"
	if !is_horiz && !slices.Contains(append(h_alignments, "center"), alignment){
		return fmt.Errorf("for direction '%v', only alignments 'left', 'right', and 'center' are valid. got '%v'", direction, alignment)
	}

	// getting to the end with no errors means we are officially done and the monitor definition is valid
	return nil
}

func find_monitor_order(conf Config) ([]string, error) {

	// Return whether a graph is disconnected, given a list of its nodes in any order.
	var graph_disconnected = func (g graph.Graph[string,string], vertices []string) bool {
		// use a BFS to count the vertices connected to vertices[0]
		count := 0
		graph.BFS(g, vertices[0], func (string) bool {
			count++
			return false
		})

		// compare to the total number of vertices
		return count != len(vertices)
	}

	// first pass: add a vertex for each monitor
	g := graph.New(graph.StringHash, graph.PreventCycles(), graph.Directed())
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
