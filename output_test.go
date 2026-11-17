package main

import (
	stdpng "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

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

func TestUppercasePNGOutput(t *testing.T) {
	if runtime.GOOS == "linux" {
		if _, err := exec.LookPath("rsvg-convert"); err != nil {
			t.Skip("Linux resvg-go crash with Go 1.24+ is tracked in #203")
		}
	}

	output := filepath.Join(t.TempDir(), "capture.PNG")
	cmd := exec.Command(binary, "test/input/artichoke.hs", "--output", output)
	if outputText, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("freeze failed: %v\n%s", err, outputText)
	}

	f, err := os.Open(output)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := stdpng.DecodeConfig(f); err != nil {
		t.Fatalf("output is not a valid PNG: %v", err)
	}
}
