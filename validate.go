package libmonpos

import (
	"fmt"
	"slices"
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
