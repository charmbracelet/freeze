package main

import "os"

func systemFontDirs() []string {
	var dirs []string
	if windir := os.Getenv("WINDIR"); windir != "" {
		dirs = append(dirs, windir+`\Fonts`)
	} else {
		dirs = append(dirs, `C:\Windows\Fonts`)
	}
	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		dirs = append(dirs, localAppData+`\Microsoft\Windows\Fonts`)
	}
	return dirs
}
