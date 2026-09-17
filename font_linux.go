package main

import "os"

func systemFontDirs() []string {
	dirs := []string{
		"/usr/share/fonts",
		"/usr/local/share/fonts",
	}
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, home+"/.local/share/fonts", home+"/.fonts")
	}
	return dirs
}
