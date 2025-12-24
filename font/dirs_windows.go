//go:build windows

package font

import (
	"os"
	"path/filepath"
)

func init() {
	windir := os.Getenv("WINDIR")
	if windir == "" {
		windir = `C:\Windows`
	}
	localAppData := os.Getenv("LOCALAPPDATA")

	DefaultFontsDirs = []string{
		filepath.Join(windir, "Fonts"),
	}

	// User fonts directory (Windows 10 1809+)
	if localAppData != "" {
		DefaultFontsDirs = append(DefaultFontsDirs,
			filepath.Join(localAppData, "Microsoft", "Windows", "Fonts"))
	}
}
