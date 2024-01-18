package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var ErrorPadding = lipgloss.NewStyle().Padding(1, 2)
var ErrorHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#FF5F87")).Bold(true).Padding(0, 1).SetString("ERROR")
var ErrorDetails = lipgloss.NewStyle().Foreground(lipgloss.Color("#757575"))

func printFatalError(title string, err error) {
	fmt.Printf("%s", ErrorPadding.Render(ErrorHeader.String(), title))
	fmt.Printf("%s\n", ErrorPadding.Render(ErrorDetails.Render(err.Error())))
	os.Exit(1)
}
