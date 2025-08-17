package libmonpos

import (
	"fmt"
	"math"
)

func round_div(a int, b int) int {
	decimal := float64(a) / float64(b)
	return int(math.Round(decimal))
}

func GeneratePositions(conf Config, g MonitorGraph) (map[string]Rect, error) {
	// get the topological order of the monitor graph
	order, err := get_monitor_order(g)
	if err != nil {
		return nil, err
	}

	// Put the first monitor at (0, 0)
	positions := make(map[string]Rect)

	// Build out from there, placing things in relation
	for _, monitor := range order {
		// divide and round the monitor's dimensions to get a scaled (AKA logical) size
		mon_conf := conf.Monitors[monitor]
		width_scaled := int(math.Round(float64(mon_conf.Width) / mon_conf.Scale))
		height_scaled := int(math.Round(float64(mon_conf.Height) / mon_conf.Scale))

		// special case out order[0], the root monitor.
		if monitor == order[0] {
			positions[monitor] = make_rect(Pair{0, 0}, width_scaled, height_scaled)
			continue
		}

		// split the position of this monitor
		direction, parent, err := split_position(mon_conf.Position)
		if err != nil {
			return nil, err
		}

		// based on the direction, figure out either the x position or y position of this rectangle
		parent_rect := positions[parent]
		var x_pos int = 0
		var y_pos int = 0
		var unset = ""
		switch direction {
		case "left-of":
			// left-of parent: your R.X = their L.X
			//                 your L.X + your Size.X = their L.X
			//                 your L.X = their L.X - your Size.X
			x_pos = parent_rect.L.X - width_scaled
			unset = "y"
		case "right-of":
			// right-of parent: your L.X = their R.X
			x_pos = parent_rect.R.X
			unset = "y"
		case "above":
			// above parent: your R.Y = their L.Y
			//               your L.Y + your Size.Y = their L.Y
			//               your L.Y = their L.Y - your Size.Y
			y_pos = parent_rect.L.Y - height_scaled
			unset = "x"
		case "below":
			// below parent: your L.Y = their R.Y
			y_pos = parent_rect.R.Y
			unset = "x"
		}

		// this will leave one of them undefined, so find it and define it based on alignment
		if unset == "y" {
			// direction was left-of or right-of
			switch mon_conf.Align {
			case "top":
				y_pos = parent_rect.L.Y
			case "bottom":
				y_pos = parent_rect.R.Y - height_scaled
			case "center":
				y_pos = parent_rect.L.Y - round_div(height_scaled - parent_rect.Size.Y, 2)
			}
		} else {
			// direction was above or below
			switch mon_conf.Align {
			case "left":
				x_pos = parent_rect.L.X
			case "right":
				x_pos = parent_rect.R.X - height_scaled
			case "center":
				x_pos = parent_rect.L.X - round_div(width_scaled - parent_rect.Size.X, 2)
			}
		}

		// add that monitor to the map
		positions[monitor] = make_rect(Pair{x_pos, y_pos}, width_scaled, height_scaled)
	}

	// Check if there are any overlaps, scream and cry if so
	var error_message = ""
	for name1, rect1 := range positions {
		fmt.Printf("%v: %v\n", name1, rect1)

		for name2, rect2 := range positions {
			if name1 == name2 {
				continue
			}

			// check for overlap
			if rect1.Overlaps(rect2) {
				message := fmt.Sprintf("OVERLAP between %v (%v) and %v (%v)\n", name1, rect1, name2, rect2)
				error_message += message
			}
		}
	}

	if error_message != "" {
		return positions, fmt.Errorf(error_message)
	} else {
		return positions, nil
	}
}

