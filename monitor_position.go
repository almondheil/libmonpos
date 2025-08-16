package monitor_position

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

/// Valid directions
var directions = []string{"left-of", "right-of", "above", "below"}
var directions_horiz = []string{"left-of", "right-of"}

/// Valid alignments
var alignments_horiz = []string{"center", "left", "right"}
var alignments_vert = []string{"center", "top", "bottom"}


/// One monitor in the config file.
///
/// The width, height, and scale keys are required, and the position and align are optional.
type Monitor struct {
	Width int
	Height int
	Scale float32
	Position string `yaml:",omitempty"`
	Align string `yaml:",omitempty"`
}

/// The entire config file is made up of multiple Monitor entries, under a common monitors header.
type Config struct {
	Monitors map[string]Monitor
}

/// Read a config file from disk and ensure that it is valid.
///
/// If position has a value but align is not specified, it will be defaulted to "center"
func read_config(path string) (error, Config) {
	// read the file from the system
	data, err := os.ReadFile(path)
	if err != nil {
		return err, Config{}
	}

	// unmarshal from yaml into the correct structure
	conf := Config{}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		return err, Config{}
	}

	// validate the config before returning it
	err = validate_config(conf)
	if err != nil {
		return err, Config{} 
	}

	return nil, conf
}

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
		return fmt.Errorf("position must be of the form <direction> <monitor>")
	}

	// make sure the first and second word are actually a direction and name respectively
	if !slices.Contains(directions, parts[0]) {
		return fmt.Errorf("expected direction above, below, left-of, or right-of, got %v", parts[0]) 
	}
	if !slices.Contains(monitor_names, parts[1]) {
		return fmt.Errorf("expected a monitor name, got %v", parts[1])
	}

	// if m.Align is empty, now we can set it to the always-valid value "center"
	if m.Align == "" {
		m.Align = "center"
		return nil
	}

	// depending on if the direction is horizontal, decide whether the alignment is valid or not
	is_horiz := slices.Contains(directions_horiz, parts[0])
	if is_horiz && !slices.Contains(alignments_vert, m.Align) {
		return fmt.Errorf("for direction %v, only alignments top, bottom, and center are valid. got %v", parts[0], m.Align)
	} else if !is_horiz {
		return fmt.Errorf("for direction %v, only alignments left, right, and center are valid. got %v", parts[0], m.Align)
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
