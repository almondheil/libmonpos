package libmonpos

import (
	"fmt"
	"os"
	"testing"
)

func TestDirectionAlignmentPairs(t *testing.T) {
	var empty = []string{""}
	var center = []string{"center"}
	var invalid = []string{
		"next-to",
		"who-cares",
		"gleeble",
	}

	var tests = []struct {
		name string
		directions []string
		alignments []string
		expected bool
	}{
			{"dir horiz align vert", h_directions, v_alignments, true},
			{"dir horiz align horiz", h_directions, h_alignments, false},
			{"dir horiz align center", h_directions, center, true},
			{"dir vert align horiz", v_directions, h_alignments, true},
			{"dir vert align vert", v_directions, v_alignments, false},
			{"dir vert align center", v_directions, center, true},
			{"dir horiz align unspecified", h_directions, empty, true},
			{"dir vert align unspecified", v_directions, empty, true},
			{"dir unspecified align horiz", empty, h_alignments, false},
			{"dir unspecified align vert", empty, v_alignments, false},
			{"dir unspecified align unspecified", empty, empty, true},
			{"dir invalid align vert", invalid, v_alignments, false},
			{"dir invalid align horiz", invalid, h_alignments, false},
			{"dir invalid align center", invalid, center, false},
			{"dir horiz align invalid", h_directions, invalid, false},
			{"dir vert align invalid", v_directions, invalid, false},
			{"dir invalid align invalid", invalid, invalid, false},
	}

	// for each category of position and alignment
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T) {

			// for each individual combo of (position, alignment) within that category
			for _, dir := range tt.directions {
				for _, align := range tt.alignments {
					testname := fmt.Sprintf("'%v' and '%v'", dir, align)
					t.Run(testname, func (t *testing.T) {

						// validate the monitor, and check the result matches the expected value
						err := check_direction_alignment(dir, align)
						valid := err == nil
						if valid != tt.expected {
							t.Errorf("'%v' and '%v': %v", dir, align, err)
						}
					})
				}
			}
		})
	}
}

func TestValidateConfigFiles(t *testing.T) {
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

			_, _, err = LoadConfig(tt.conf_file)
			actual := err == nil
			if tt.expected != actual {
				t.Errorf("%v: %v", tt.conf_file, err)
			}
		})
	}
}

func TestDefaults(t *testing.T) {
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
			conf, _, err := LoadConfig(tt.conf_file)
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

/*
// confirm the monitor order is not contradictory ever
func TestMonitorOrder(t *testing.T) {

}
*/

