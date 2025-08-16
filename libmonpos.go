package libmonpos

import (
	"os"

	"gopkg.in/yaml.v3"
)

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
