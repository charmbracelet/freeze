package main

import (
	"os"
	"regexp"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/alecthomas/kong"
	"github.com/charmbracelet/colorprofile"
)

const space = 18

var highlighter = regexp.MustCompile("{{(.+?)}}")

func helpPrinter(_ kong.HelpOptions, ctx *kong.Context) error {
	compat.HasDarkBackground = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	compat.Profile = colorprofile.Detect(os.Stdout, os.Environ())

	codeBlockStyle := lipgloss.NewStyle().Background(compat.AdaptiveColor{Light: lipgloss.Color("254"), Dark: lipgloss.Color("235")}).MarginLeft(2).Padding(1, 2)
	programStyle := lipgloss.NewStyle().Background(codeBlockStyle.GetBackground()).Foreground(lipgloss.Color("#7E65FF")).PaddingLeft(1)
	stringStyle := lipgloss.NewStyle().Background(codeBlockStyle.GetBackground()).Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#02BA84"), Dark: lipgloss.Color("#02BF87")}).PaddingLeft(1)
	argumentStyle := lipgloss.NewStyle().Background(codeBlockStyle.GetBackground()).Foreground(lipgloss.Color("248")).PaddingLeft(1)
	flagStyle := lipgloss.NewStyle().Background(codeBlockStyle.GetBackground()).Foreground(lipgloss.Color("244")).PaddingLeft(1)
	titleStyle := lipgloss.NewStyle().Bold(true).Transform(strings.ToUpper).Margin(1, 0, 0, 2).Foreground(lipgloss.Color("#6C50FF"))

	lipgloss.Println()
	lipgloss.Println("  Generate images of code and terminal output. 📸")

	lipgloss.Println(titleStyle.Render(strings.ToUpper("Usage")))
	lipgloss.Println()
	lipgloss.Println(
		codeBlockStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinHorizontal(lipgloss.Left, programStyle.Render("freeze"), argumentStyle.Render("main.go"), flagStyle.Render("[-o code.svg] [--flags]")),
				lipgloss.JoinHorizontal(lipgloss.Left, programStyle.Render("freeze"), argumentStyle.Render("--execute"), stringStyle.Render("\"ls -la\""), flagStyle.Render("[--flags]   ")),
			),
		),
	)

	flags := ctx.Flags()
	lastGroup := ""

	lipgloss.Println()
	for _, f := range flags {
		if f.Name == "interactive" {
			printFlag(f)
		}
	}

	lipgloss.Println(titleStyle.Render("Settings"))

	for _, f := range flags {
		if f.Group != nil && f.Group.Title == "Settings" {
			if f.Hidden || f.Name == "help" {
				continue
			}
			printFlag(f)
		}
	}

	lipgloss.Print(titleStyle.Render("Customization"))

	for _, f := range flags {
		if f.Hidden || f.Name == "help" || f.Group.Title == "Settings" {
			continue
		}

		if f.Group != nil && lastGroup != f.Group.Title {
			lastGroup = f.Group.Title
			lipgloss.Println()
		}

		printFlag(f)
	}
	lipgloss.Println()
	return nil
}

const helpForeground = "243"

func printFlag(f *kong.Flag) {
	dashStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginLeft(1)
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(helpForeground))
	keywordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	if f.Short > 0 {
		lipgloss.Print("    ", dashStyle.Render("-"), string(f.Short))
		lipgloss.Print(dashStyle.Render("--"), f.Name)
		lipgloss.Print(strings.Repeat(" ", space-len(f.Name)))
	} else {
		lipgloss.Print("    ", dashStyle.Render(" "), " ")
		lipgloss.Print(dashStyle.Render("--"), f.Name)
		lipgloss.Print(strings.Repeat(" ", space-len(f.Name)))
	}
	help := highlighter.ReplaceAllString(f.Help, keywordStyle.Render("$1")+"\x1b[38;5;"+helpForeground+"m")
	lipgloss.Println(helpStyle.Render(help))
}
