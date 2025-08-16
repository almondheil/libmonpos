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
			{"pos horiz align vert", h_positions, v_alignments, true},
			{"pos horiz align horiz", h_positions, h_alignments, false},
			{"pos horiz align both", h_positions, []string{"center"}, true},
			{"pos vert align horiz", v_positions, h_alignments, true},
			{"pos vert align vert", v_positions, v_alignments, false},
			{"pos vert align both", v_positions, []string{"center"}, true},
			{"pos horiz align unspecified", h_positions, empty, true},
			{"pos vert align unspecified", v_positions, empty, true},
			{"pos unspecified align horiz", empty, h_alignments, false},
			{"pos unspecified align vert", empty, v_alignments, false},
			{"pos unspecified align unspecified", empty, empty, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T) {
			for _, pos := range tt.positions {
				for _, align := range tt.aligns {
					testname := fmt.Sprintf("'%v' and '%v'", pos, align)

					t.Run(testname, func (t_sub *testing.T) {
						// generate a potential monitor
						//
						// if pos is empty let it be empty in the monitor definition, but if it is a position
						// place it in relation to the "fake" monitor
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
							t_sub.Errorf("%v: %v", mon, err)
						}
					})
				}
			}
		})
	}
}

// TODO: IS THE POSITION VALID?
// TODO: NONSENSICAL VALUES FOR POS AND ALIGN
