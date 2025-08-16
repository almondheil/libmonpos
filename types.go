package libmonpos

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
