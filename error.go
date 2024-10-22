package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	errorHeader = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1F1F1")).
			Background(lipgloss.Color("#FF5F87")).
			Bold(true).
			Padding(0, 1).
			Margin(1).
			MarginLeft(2).
			SetString("ERROR")
	errorDetails = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#757575")).
			MarginLeft(2)
)

func printError(title string, err error) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Center, errorHeader.String(), title))
	fmt.Println(errorDetails.Render(err.Error()))
}

func printErrorFatal(title string, err error) {
	printError(title, err)
	os.Exit(1)
}
