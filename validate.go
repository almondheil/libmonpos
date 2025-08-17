package libmonpos

import (
	"fmt"
	"slices"
	"strings"
)

// Categorization of which directions and alignments are horizontal vs vertical
var h_directions = []string{"left-of", "right-of"}
var v_directions = []string{"above", "below"}
var h_alignments = []string{"left", "right"}
var v_alignments = []string{"top", "bottom"}

// Given a position in the format `<direction> <monitor>`, split it into two parts and return them.
// Returns an error if it is not two space-separated words.
//
// Special case: if position is an empty string, returns two empty string halves
func split_position(position string) (string, string, error) {
	// if the position is empty, just return two empty halves
	if position == "" {
		return "", "", nil
	}

	parts := strings.Split(position, " ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("position should be of the form '<direction> <monitor>', got '%v'", position)
	}

	return parts[0], parts[1], nil
}

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

	// h directions work with v alignments, "center", or ""
	// v directions work with h alignments
	if is_horiz && !slices.Contains(append(v_alignments, "center", ""), alignment) {
		return fmt.Errorf("for direction '%v', only alignments 'top', 'bottom', and 'center' are valid. got '%v'", direction, alignment)
	} else if !is_horiz && !slices.Contains(append(h_alignments, "center", ""), alignment){
		return fmt.Errorf("for direction '%v', only alignments 'left', 'right', and 'center' are valid. got '%v'", direction, alignment)
	}

	// getting to the end with no errors means we are officially done and the monitor definition is valid
	return nil
}

