package libmonpos

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct{
		conf_file string
		expected bool
	}{
		{"test_files/not_yaml.txt", false},
		{"test_files/basic.yaml", true},
		{"test_files/basic_with_yaml_header.yaml", true},
		{"test_files/complex.yaml", true},
		{"test_files/missing_names.yaml", false},
		{"test_files/cycle.yaml", false},
		{"test_files/cycle_one.yaml", false},
		{"test_files/two_roots.yaml", false},
		{"test_files/unspecified_align.yaml", true},
		{"test_files/unspecified_scale.yaml", true},
		{"test_files/unspecified_width.yaml", false},
		{"test_files/unspecified_height.yaml", false},
	}

	for _, tt := range tests {
		t.Run(tt.conf_file, func (t *testing.T) {
			// check that the file exists first
			_, err := os.Stat(tt.conf_file)
			if err != nil {
				t.Errorf("error finding config file: %v", err)
				return
			}

			_, err = LoadConfig(tt.conf_file)
			actual := err == nil
			if tt.expected != actual {
				t.Errorf("%v: %v", tt.conf_file, err)
			}
		})
	}
}

func TestApplyDefaults(t *testing.T) {
	tests := []struct{
		conf_file string
		name string
		expected Monitor
	}{
		{"test_files/unspecified_scale.yaml", "B", Monitor{20, 20, 1.0, "below A", "right"}},
		{"test_files/unspecified_align.yaml", "B", Monitor{20, 20, 2.0, "below A", "center"}},
	}

	for _, tt := range tests {
		t.Run(tt.conf_file, func (t *testing.T) {
			conf, err := LoadConfig(tt.conf_file)
			if err != nil {
				t.Errorf("unexpected error reading config: %v", err)
				return
			}

			actual := conf.Monitors[tt.name]
			if tt.expected != actual {
				t.Errorf("definitions did not match: expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
