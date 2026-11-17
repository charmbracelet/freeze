package main

import "testing"

func TestIsPNGOutput(t *testing.T) {
	for _, tc := range []struct {
		path string
		want bool
	}{
		{"output.png", true},
		{"output.PNG", true},
		{"output.PnG", true},
		{"output.svg", false},
		{"output.png.svg", false},
		{"output.png.backup", false},
	} {
		t.Run(tc.path, func(t *testing.T) {
			if got := isPNGOutput(tc.path); got != tc.want {
				t.Errorf("isPNGOutput(%q) = %t, want %t", tc.path, got, tc.want)
			}
		})
	}
}
