package libmonpos

import (
	"fmt"
	"testing"
)

func TestPositionAlignmentPairs(t *testing.T) {
	var empty = []string{""}

	var tests = []struct {
		name string
		positions []string
		aligns []string
		expected bool
	}{
			{"pos horiz align vert", h_directions, v_alignments, true},
			{"pos horiz align horiz", h_directions, h_alignments, false},
			{"pos horiz align both", h_directions, []string{"center"}, true},
			{"pos vert align horiz", v_directions, h_alignments, true},
			{"pos vert align vert", v_directions, v_alignments, false},
			{"pos vert align both", v_directions, []string{"center"}, true},
			{"pos horiz align unspecified", h_directions, empty, true},
			{"pos vert align unspecified", v_directions, empty, true},
			{"pos unspecified align horiz", empty, h_alignments, false},
			{"pos unspecified align vert", empty, v_alignments, false},
			{"pos unspecified align unspecified", empty, empty, true},
	}

	// for each category of position and alignment
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T) {

			// for each individual combo of (position, alignment) within that category
			for _, pos := range tt.positions {
				for _, align := range tt.aligns {
					testname := fmt.Sprintf("'%v' and '%v'", pos, align)
					t.Run(testname, func (t *testing.T) {
						// make config that could be a monitor.
						// (tack on a fake monitor to make the string valid if the position is not empty)
						var mon Monitor
						if pos == "" {
							mon = Monitor{1920, 1080, 1.0, pos, align}
						} else {
							mon = Monitor{1920, 1080, 1.0, pos + " fake", align}
						}

						// validate the monitor, and check the result matches the expected value
						err := validate_monitor(mon, []string{"fake"})
						valid := err == nil
						if valid != tt.expected {
							t.Errorf("%v: %v", mon, err)
						}
					})
				}
			}
		})
	}
}

func TestRelativePositions(t *testing.T) {
	// this list of monitor names will be used for all these tests
	monitor_list := []string{"mon1", "mon2", "mon3", "mon4"}

	tests := []struct{
		target_name string
		expected bool
	}{
		{"mon1", true},
		{"mon2", true},
		{"mon3", true},
		{"mon4", true},
		{"mon5", false},
		{"fake", false},
	}

	// We expect it to return no error when our monitor is in the list, and an error when it's not
	for _, tt := range tests {
		t.Run(tt.target_name, func (t *testing.T) {
			for _, direction := range append(h_directions, v_directions...) {
				mon := Monitor{1920, 1080, 1.0, direction + " " + tt.target_name, "center"}
				err := validate_monitor(mon, monitor_list)

				value := err == nil
				if value != tt.expected {
					t.Errorf("%v: %v", mon, err)
				}
			}
		})
	}
}

func TestNonsensicalPositions(t *testing.T) {
	// for these tests, "mon1" is a valid monitor
	tests := []string{
		"diagonal-to mon1",
		"glerb",
		"left-of",
		"mon1 above",
	}

	for _, position := range tests {
		for _, alignment := range append(h_alignments, append(v_alignments, "center")...) {
			testname := fmt.Sprintf("%v %v", position, alignment)
			t.Run(testname, func (t *testing.T) {
				mon := Monitor{1920, 1080, 1.0, position, alignment}
				err := validate_monitor(mon, []string{"mon1"})

				if err == nil {
					t.Errorf("%v: %v", mon, err)
				}
			})
		}
	}
}

func TestNonsensicalAlignments(t *testing.T) {
	tests := []string{
		"diagonal",
		"nowhere",
		"same",
	}

	for _, alignment := range tests {
		for _, direction := range append(h_directions, v_directions...) {
			position := direction + " mon1"
			testname := fmt.Sprintf("%v %v", position, alignment)

			t.Run(testname, func (t *testing.T) {
				mon := Monitor{1920, 1080, 1.0, position, alignment}
				err := validate_monitor(mon, []string{"mon1"})

				if err == nil {
					t.Errorf("%v: %v", mon, err)
				}
			})
		}
	}
}

// TODO: CONFIG VALIDATION AS A WHOLE
