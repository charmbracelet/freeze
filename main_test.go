package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("no home dir available: %v", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		in   string
		want string
	}{
		{"tilde alone", "~", home},
		{"tilde slash file", "~/out.png", filepath.Join(home, "out.png")},
		{"tilde slash nested", "~/Pictures/out.png", filepath.Join(home, "Pictures", "out.png")},
		{"relative file", "out.png", filepath.Join(cwd, "out.png")},
		{"relative parent", "../out.png", filepath.Clean(filepath.Join(cwd, "..", "out.png"))},
		{"absolute file", "/tmp/out.png", "/tmp/out.png"},
		{"tilde in middle is left alone", "foo/~/bar", filepath.Join(cwd, "foo/~/bar")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := expandPath(tc.in)
			if err != nil {
				t.Fatalf("expandPath(%q) returned err: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("expandPath(%q) = %q, want %q", tc.in, got, tc.want)
			}
			if strings.Contains(got, "~") && !strings.Contains(tc.in[1:], "~") {
				t.Errorf("result still contains ~: %q", got)
			}
		})
	}
}
