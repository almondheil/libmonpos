package libmonpos

import (
	"fmt"
	"slices"
	"strings"
)

// categorization of which positions and alignments are horizontal vs vertical
var h_positions = []string{"left-of", "right-of"}
var v_positions = []string{"above", "below"}
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

	// check the general format of position, <word> <word>
	parts := strings.Split(m.Position, " ")
	if len(parts) != 2 {
		return fmt.Errorf("position must be of the form <position> <monitor>")
	}
	position := parts[0]
	align := parts[1]

	// make sure the first and second word are actually a position and name respectively
	if !slices.Contains(h_positions, position) && !slices.Contains(v_positions, position) {
		return fmt.Errorf("expected position 'above', 'below', 'left-of', or 'right-of', got '%v'", position) 
	}
	if !slices.Contains(monitor_names, align) {
		return fmt.Errorf("expected a monitor name, got '%v'", align)
	}

	// if m.Align is empty, now we can set it to the always-valid value "center"
	if m.Align == "" {
		m.Align = "center"
		return nil
	}

	// depending on if the position is horizontal, decide whether the alignment is valid or not
	is_horiz := slices.Contains(h_positions, position)

	// h positions work with v alignments, or "center"
	if is_horiz && !slices.Contains(v_alignments, m.Align) && m.Align != "center"  {
		return fmt.Errorf("for position '%v', only alignments 'top', 'bottom', and 'center' are valid. got '%v'", position, m.Align)
	}
	// v positions work with h alignments, or "center"
	if !is_horiz && !slices.Contains(h_alignments, m.Align) && m.Align != "center" {
		return fmt.Errorf("for position '%v', only alignments 'left', 'right', and 'center' are valid. got '%v'", position, m.Align)
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
