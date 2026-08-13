package main

import "os"

func systemFontDirs() []string {
	dirs := []string{
		"/System/Library/Fonts",
		"/System/Library/Fonts/Supplemental",
		"/Library/Fonts",
	}
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, home+"/Library/Fonts")
	}
	return dirs
}
