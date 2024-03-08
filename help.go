package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/lipgloss"
)

const space = 18

var highlighter = regexp.MustCompile("{{(.+?)}}")

func helpPrinter(_ kong.HelpOptions, ctx *kong.Context) error {
	codeBlockStyle := lipgloss.NewStyle().Background(lipgloss.Color("0")).Padding(1, 0)
	programStyle := codeBlockStyle.Copy().Foreground(lipgloss.Color("12")).PaddingLeft(3).MarginLeft(2)
	argumentStyle := codeBlockStyle.Copy().Foreground(lipgloss.Color("7")).Padding(1, 1)
	flagStyle := codeBlockStyle.Copy().Foreground(lipgloss.Color("244")).PaddingRight(3)
	titleStyle := lipgloss.NewStyle().Bold(true).Margin(1, 0, 0, 2).Foreground(lipgloss.Color("#875FFF"))
	dashStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginLeft(2)
	keywordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	fmt.Println()
	fmt.Println("  Generate images of code and terminal output. 📸")
	fmt.Println()
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Center, programStyle.Render("freeze"), argumentStyle.Render("main.go"), flagStyle.Render("[-o code.svg] [--flags]")))

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
			fmt.Print("  ", dashStyle.Render("-"), string(f.Short))
			fmt.Print(dashStyle.Render("--"), f.Name)
			fmt.Print(strings.Repeat(" ", space-len(f.Name)))
		} else {
			fmt.Print("  ", dashStyle.Render(" "), " ")
			fmt.Print(dashStyle.Render("--"), f.Name)
			fmt.Print(strings.Repeat(" ", space-len(f.Name)))

		}
		help := highlighter.ReplaceAllString(f.Help, keywordStyle.Render("$1"))
		fmt.Println(help)
	}
	fmt.Println()
	return nil
}
