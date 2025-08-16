package libmonpos

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// One monitor in the config file.
type Monitor struct {
	Width uint
	Height uint
	Scale float32  `yaml:",omitempty"`
	Position string `yaml:",omitempty"`
	Align string `yaml:",omitempty"`
}

// The config file is made up of multiple Monitor entries, under a common monitors header.
type Config struct {
	Monitors map[string]Monitor
}

// Format monitor info as a string.
func (m Monitor) String() string {
	dimensions := fmt.Sprintf("Monitor{%dx%d@%.2fx", m.Width, m.Height, m.Scale)
	if m.Position == "" {
		return dimensions + "}"
	}

	position := fmt.Sprintf("%v align %v", m.Position, m.Align)
	return dimensions + " " + position + "}"
}

// Read a config file from disk. Does NOT verify the position and alignment of the monitor.
func read_config_yaml(path string) (Config, error) {
	// read the file from the system
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	// unmarshal from yaml into the correct structure
	conf := Config{}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}

// Apply defaults to all monitors in a config
func apply_defaults(c Config) {
	for name, mon := range c.Monitors {
		// default align: center (if position is specified)
		if mon.Position != "" && mon.Align == "" {
			mon.Align = "center"
			c.Monitors[name] = mon
		}

		// default scale: 1.0
		if mon.Scale == 0.0 {
			mon.Scale = 1.0
			c.Monitors[name] = mon
		}
	}
}

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

// Read a config file from disk and check that it is valid,
// generating an order to arrange monitors in the process.
func LoadConfig(path string) (Config, []string, error) {
	conf, err := read_config_yaml(path)
	if err != nil {
		return Config{}, []string{}, err
	}

	// Apply any defaults to that config
	apply_defaults(conf)

	// For all monitors, check that their direction and alignment are valid and agree with each other.
	for name, monitor := range conf.Monitors {
		// If the dimensions are unspecified, scream and cry
		if monitor.Width == 0 || monitor.Height == 0 || monitor.Scale == 0{
			return Config{}, []string{}, fmt.Errorf("monitor '%v' width, height, and scale must be specified and nonzero", name)
		}

		// Get the direction, either from splitting or just leaving it empty
		direction, _, err := split_position(monitor.Position)
		if err != nil {
			return Config{}, []string{}, err
		}
		err = check_direction_alignment(direction, monitor.Align)
		if err != nil {
			return Config{}, []string{}, err
		}
	}

	// Topologically sort the monitors to get an order to work with them in.
	// Doing this, we also check that all neighbor names are valid and form a tree with no disconnection.
	order, err := find_monitor_order(conf)
	if err != nil {
		return Config{}, []string{}, err
	}

	// At the end, we have read the config and know the topological order of the monitors
	return conf, order, nil
}
