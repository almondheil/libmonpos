package libmonpos

import (
	"fmt"

	"github.com/dominikbraun/graph"
)

// One monitor in the config file.
type Monitor struct {
	Width uint
	Height uint
	Scale float64  `yaml:",omitempty"`
	Position string `yaml:",omitempty"`
	Align string `yaml:",omitempty"`
}

// The config file is made up of multiple Monitor entries, under a common monitors header.
type Config struct {
	Monitors map[string]Monitor
}

// A graph representing connections between monitors
type MonitorGraph graph.Graph[string,string]

// A rectangle defined by the top-left and bottom-right corner.
type Rect struct {
	L Pair
	R Pair
	Size Pair
}

// An (x, y) pair
type Pair struct {
	X int
	Y int
}

// Construct a rect from the upper-left corner and dimensions.
func make_rect(pos Pair, width, height int) Rect {
	return Rect{pos, Pair{pos.X + width, pos.Y + height}, Pair{width, height}}
}

// Decide whether two Rects are overlapping
func (r1 Rect) Overlaps(r2 Rect) bool {
	// top-left corner of one is to the right of the bottom-right corner of the other
	if r1.L.X >= r2.R.X || r2.L.X >= r1.R.X {
		return false
	}

	// top-left corner of one is below the bottom-right corner of the other
	if r1.L.Y >= r2.R.Y || r2.L.Y >= r1.R.Y {
		return false
	}

	// if neither of those, there's overlap
	return true
}

// Format monitor info as a string.
func (m Monitor) String() string {
	result := fmt.Sprintf("%vx%v", m.Width, m.Height)

	// add on scale if non-standard
	if m.Scale != 1.0 {
		result += fmt.Sprintf("@%2.3vx", m.Scale)
	}

	// tack on position + align if position is defined
	if m.Position == "" {
		return result
	} else {
		result += fmt.Sprintf(" %v align %v", m.Position, m.Align)
		return result
	}
}
