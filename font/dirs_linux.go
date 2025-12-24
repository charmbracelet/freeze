//go:build linux

package font

import (
	"os"
	"path/filepath"
)

func init() {
	DefaultFontsDirs = []string{
		"/usr/share/fonts",
		"/usr/local/share/fonts",
		"/usr/share/X11/fonts/Type1",
		"/usr/share/X11/fonts/TTF",
	}

	// Only add user font directories if home dir is available
	if home, err := os.UserHomeDir(); err == nil {
		DefaultFontsDirs = append(DefaultFontsDirs,
			filepath.Join(home, ".fonts"),
			filepath.Join(home, ".local/share/fonts"),
		)
	}
}
