package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/lipgloss"
)

func helpPrinter(options kong.HelpOptions, ctx *kong.Context) error {
	var titleStyle = lipgloss.NewStyle().Bold(true).Margin(1, 0, 0, 2).Foreground(lipgloss.Color("#875FFF"))
	var codeBlockStyle = lipgloss.NewStyle().Background(lipgloss.Color("0")).Padding(1, 3).Margin(0, 2)

	fmt.Println()
	fmt.Println("  Screenshot code on the command line.")
	fmt.Println()
	fmt.Println(codeBlockStyle.Render(foreground("freeze", 13) + " main.go " + foreground("[-o code.svg] [--flags]", 244)))

	flags := ctx.Flags()
	lastGroup := ""

	for _, f := range flags {
		if f.Name == "help" {
			continue
		}

		if f.Group != nil && lastGroup != f.Group.Title {
			lastGroup = f.Group.Title
			fmt.Println(titleStyle.Render(strings.ToUpper(f.Group.Title)))
		}

		if f.Short > 0 {
			fmt.Printf("    %s%c", foreground("-", 8), f.Short)
			fmt.Printf("  %s%s", foreground("--", 8), f.Name)
			fmt.Print(strings.Repeat(" ", 16-len(f.Name)))
		} else {
			fmt.Printf("    %s%c", " ", ' ')
			fmt.Printf("  %s%s", foreground("--", 8), f.Name)
			fmt.Print(strings.Repeat(" ", 16-len(f.Name)))

		}
		help := strings.ReplaceAll(f.Help, "{{", color(1))
		help = strings.ReplaceAll(help, "}}", color(7))
		fmt.Println(help)
	}
	fmt.Println()
	return nil
}

func color(c int) string {
	return fmt.Sprintf("\x1b[38;5;%dm", c)
}

func foreground(s string, c int) string {
	return color(c) + s + color(7)
}
