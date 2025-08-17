package libmonpos

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

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

// Read a config file from disk and check that it is valid,
func LoadConfig(path string) (Config, error) {
	conf, err := read_config_yaml(path)
	if err != nil {
		return Config{}, err
	}

	// Apply any defaults to that config
	apply_defaults(conf)

	// For all monitors, check that their direction and alignment are valid and agree with each other.
	for name, monitor := range conf.Monitors {
		// If the dimensions are unspecified, scream and cry
		if monitor.Width == 0 || monitor.Height == 0 || monitor.Scale == 0{
			return Config{}, fmt.Errorf("monitor '%v' width, height, and scale must be specified and nonzero", name)
		}

		// Get the direction, either from splitting or just leaving it empty
		direction, _, err := split_position(monitor.Position)
		if err != nil {
			return Config{}, err
		}
		err = check_direction_alignment(direction, monitor.Align)
		if err != nil {
			return Config{}, err
		}
	}

	return conf, nil
}
