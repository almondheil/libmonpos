package libmonpos

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

/// One monitor in the config file.
///
/// The width, height, and scale keys are required, and the position and align are optional.
type Monitor struct {
	Width uint
	Height uint
	Scale float32
	Position string `yaml:",omitempty"`
	Align string `yaml:",omitempty"`
}

/// The entire config file is made up of multiple Monitor entries, under a common monitors header.
type Config struct {
	Monitors map[string]Monitor
}

func (m Monitor) String() string {
	dimensions := fmt.Sprintf("Monitor{%dx%d@%.2fx", m.Width, m.Height, m.Scale)
	if m.Position == "" {
		return dimensions + "}"
	}

	position := fmt.Sprintf("%v, align %v", m.Position, m.Align)
	return dimensions + " " + position + "}"
}

/// Read a config file from disk. Does NOT verify the position and alignment of the monitor.
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

/// Apply defaults across a config
func apply_defaults(c Config) {
	for name, mon := range c.Monitors {
		// align defaults to center when unspecified
		if mon.Position != "" && mon.Align == "" {
			mon.Align = "center"
			c.Monitors[name] = mon
		}
	}
}

/// Read a config file from disk and ensure that it is valid.
///
/// If position has a value but align is not specified, it will be defaulted to "center"
func LoadConfig(path string) (Config, error) {
	conf, err := read_config_yaml(path)
	if err != nil {
		return Config{}, err
	}

	// Apply any defaults to that config
	apply_defaults(conf)

	err = validate_config(conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
