package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	errorHeader  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#FF5F87")).Bold(true).Padding(0, 1).Margin(1).MarginLeft(2).SetString("ERROR")
	errorDetails = lipgloss.NewStyle().Foreground(lipgloss.Color("#757575"))
)

func printError(title string, err error) {
	// fmt.Println(lipgloss.JoinHorizontal(lipgloss.Center, errorHeader.String(), title))
	var b strings.Builder
	e := strings.TrimSpace(err.Error())
	e = strings.ReplaceAll(e, "\r", "")
	e = strings.ReplaceAll(e, "\t", "")
	str := strings.Split(e, "\t")
	for _, s := range str {
		b.WriteString(errorDetails.Render(s))
	}
	fmt.Println(b.String())
}

func printErrorFatal(title string, err error) {
	printError(title, err)
	os.Exit(1)
}
