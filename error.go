package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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
			Background(lipgloss.Color("52")).
			Foreground(lipgloss.Color("#757575")).
			Margin(0, 0, 1, 2)
)

func printError(title string, err error) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Center, errorHeader.String(), title))
	rendered := errorDetails.Render(err.Error())
	fmt.Println(ansi.Strip(rendered))
}

func printErrorFatal(title string, err error) {
	printError(title, err)
	os.Exit(1)
}
