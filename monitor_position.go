package main // monitor_position

import (
	"fmt"
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

/// The entire config file is made up of multiple Monitor entries, under a common Monitor header.
type Config struct {
	Monitors map[string]Monitor
}

func read_config(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	conf := Config{}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		panic(err)
	}

	return conf
}


func main() {
	config := read_config("./example.yaml")
	fmt.Printf("%v\n", config)
}
