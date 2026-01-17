//go:build darwin

package font

import (
	"os"
	"path/filepath"
)

func init() {
	DefaultFontsDirs = []string{
		"/System/Library/Fonts",
		"/Library/Fonts",
	}

	// Only add user font directory if home dir is available
	if home, err := os.UserHomeDir(); err == nil {
		DefaultFontsDirs = append(DefaultFontsDirs, filepath.Join(home, "Library/Fonts"))
	}
}
