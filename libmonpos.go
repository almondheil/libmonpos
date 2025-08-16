package libmonpos

import (
	"os"

	"gopkg.in/yaml.v3"
)

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
func ReadConfig(path string) (error, Config) {
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
