package libmonpos

import (
	"fmt"
	"slices"
	"strings"
)

// categorization of which directions and alignments are horizontal vs vertical
var h_directions = []string{"left-of", "right-of"}
var v_directions = []string{"above", "below"}
var h_alignments = []string{"left", "right"}
var v_alignments = []string{"top", "bottom"}

/// Validate a monitor config, ensuring that the position and align are valid if they are defined.
///
/// If position has a value but align is not specified, it will be defaulted to "center"
func validate_monitor(m Monitor, monitor_names []string) error {
	if m.Position == "" && m.Align == "" {
		return nil // position and align may both be empty
	} else if m.Position == "" && m.Align != "" {
		return fmt.Errorf("position is blank, so alignment must also be blank")
	}

	// check the general format of position <direction> <monitor>
	parts := strings.Split(m.Position, " ")
	if len(parts) != 2 {
		return fmt.Errorf("position must be of the form <direction> <monitor>")
	}
	direction := parts[0]
	monitor := parts[1]

	// make sure the first and second word are actually a direction and name respectively
	if !slices.Contains(append(h_directions, v_directions...), direction) {
		return fmt.Errorf("expected direction 'above', 'below', 'left-of', or 'right-of', got '%v'", direction) 
	}
	if !slices.Contains(monitor_names, monitor) {
		return fmt.Errorf("expected a monitor name, got '%v'", monitor)
	}

	// if m.Align is empty, now we can set it to the always-valid value "center"
	if m.Align == "" {
		m.Align = "center"
		return nil
	}

	// depending on if the direction is horizontal, decide whether the alignment is valid or not
	is_horiz := slices.Contains(h_directions, direction)

	// h directions work with v alignments, or "center"
	if is_horiz && !slices.Contains(append(v_alignments, "center"), m.Align) {
		return fmt.Errorf("for direction '%v', only alignments 'top', 'bottom', and 'center' are valid. got '%v'", direction, m.Align)
	}
	// v directions work with h alignments, or "center"
	if !is_horiz && !slices.Contains(append(h_alignments, "center"), m.Align) {
		return fmt.Errorf("for direction '%v', only alignments 'left', 'right', and 'center' are valid. got '%v'", direction, m.Align)
	}

	// getting to the end with no errors means we are officially done and the monitor definition is valid
	return nil
}

/// Validate all monitors in a config, returning an error as soon as it is detected.
func validate_config(c Config) error {
	monitors := c.Monitors

	// put all monitor names in a list
	monitor_names := make([]string, len(monitors))
	for name := range monitors {
		monitor_names = append(monitor_names, name)
	}

	// check the parameters of each monitor
	for name, monitor := range monitors {
		err := validate_monitor(monitor, monitor_names)
		if err != nil {
			return fmt.Errorf("monitor %v is invalid: %v", name, err)
		}
	}

	return nil
}
