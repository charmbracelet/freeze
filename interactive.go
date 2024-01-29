package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	green = lipgloss.Color("#03BF87")
)

func runForm(config *Config) (*Config, error) {
	var (
		padding      string = strings.Trim(fmt.Sprintf("%v", config.Padding), "[]")
		margin       string = strings.Trim(fmt.Sprintf("%v", config.Margin), "[]")
		fontSize     string = fmt.Sprintf("%d", int(config.Font.Size))
		lineHeight   string = fmt.Sprintf("%.1f", config.LineHeight)
		borderRadius string = fmt.Sprintf("%d", config.Border.Radius)
		borderWidth  string = fmt.Sprintf("%d", config.Border.Width)
		shadowBlur   string = fmt.Sprintf("%d", config.Shadow.Blur)
		shadowX      string = fmt.Sprintf("%d", config.Shadow.X)
		shadowY      string = fmt.Sprintf("%d", config.Shadow.Y)
	)

	theme := huh.ThemeCharm()
	theme.FieldSeparator = lipgloss.NewStyle()
	theme.Blurred.Title = theme.Blurred.Title.Copy().Width(14).Foreground(lipgloss.Color("7"))
	theme.Focused.Title = theme.Focused.Title.Copy().Width(14).Foreground(green).Bold(true)
	theme.Focused.Title = theme.Focused.Title.Copy().Width(14).Foreground(green).Bold(true)
	theme.Blurred.Description = theme.Blurred.Description.Copy().Foreground(lipgloss.Color("0"))
	theme.Focused.Description = theme.Focused.Description.Copy().Foreground(lipgloss.Color("7"))
	theme.Blurred.BlurredButton = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).PaddingRight(1)
	theme.Blurred.FocusedButton = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).PaddingRight(1)
	theme.Focused.BlurredButton = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).PaddingRight(1)
	theme.Focused.FocusedButton = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).PaddingRight(1)
	theme.Focused.NoteTitle = theme.Focused.NoteTitle.Copy().Margin(1, 0)
	theme.Blurred.NoteTitle = theme.Blurred.NoteTitle.Copy().Margin(1, 0)

	if config.Input == "" {
		_ = huh.NewForm(
			huh.NewGroup(
				huh.NewNote().Title("Capture file"),
				huh.NewFilePicker().
					Height(10).
					Value(&config.Input),
			),
		).WithTheme(theme).Run()

		base, ext := filepath.Base(config.Input), filepath.Ext(config.Input)
		config.Output = strings.TrimSuffix(base, ext) + ".svg"
	}

	theme.Focused.Base.
		Border(lipgloss.Border{Left: "> "}, false).
		BorderLeft(true).
		BorderForeground(green)

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Settings"),

			huh.NewInput().
				Title("Output").
				Placeholder("out.svg").
				// Description("Output location for image.").
				Inline(true).
				Prompt("").
				Value(&config.Output),

			huh.NewInput().Title("Theme ").
				Placeholder("charm").
				// Description("Theme for syntax highlighting.").
				Suggestions([]string{"charm", "dracula"}).
				Inline(true).
				Prompt("").
				Value(&config.Theme),

			huh.NewInput().Title("Background ").
				// Description("Apply a background fill.").
				Placeholder("#FFF").
				Value(&config.Background).
				Inline(true).
				Prompt("").
				Validate(validateColor),
		),

		huh.NewGroup(
			huh.NewNote().Title("Window"),

			huh.NewInput().Title("Padding ").
				// Description("Apply padding to the code.").
				Placeholder("20 40").
				Inline(true).
				Value(&padding).
				Prompt("").
				Validate(validatePadding),

			huh.NewInput().Title("Margin ").
				// Description("Apply margin to the window.").
				Placeholder("20").
				Inline(true).
				Value(&margin).
				Prompt("").
				Validate(validatePadding),

			huh.NewConfirm().Title("Controls").
				Inline(true).
				Value(&config.Window),
		),

		huh.NewGroup(
			huh.NewNote().Title("Font"),

			huh.NewInput().Title("Font Family ").
				// Description("Font family to use for code").
				Placeholder("JetBrains Mono").
				Inline(true).
				Prompt("").
				Value(&config.Font.Family),

			huh.NewInput().Title("Font Size ").
				// Description("Font size to use for code.").
				Placeholder("14").
				Inline(true).
				Prompt("").
				Value(&fontSize).
				Validate(validateInteger),

			huh.NewInput().Title("Line Height ").
				// Description("Line height relative to size.").
				Placeholder("1.2").
				Inline(true).
				Prompt("").
				Value(&lineHeight).
				Validate(validateFloat),
		),

		huh.NewGroup(
			huh.NewNote().Title("Border"),

			huh.NewInput().Title("Border Radius ").
				// Description("Corner radius of the window.").
				Placeholder("0").
				Inline(true).
				Prompt("").
				Value(&borderRadius).
				Validate(validateInteger),

			huh.NewInput().Title("Border Width ").
				// Description("Border width thickness.").
				Placeholder("1").
				Inline(true).
				Prompt("").
				Value(&borderWidth).
				Validate(validateInteger),

			huh.NewInput().Title("Border color ").
				// Description("Color of outline stroke.").
				Validate(validateColor).
				Inline(true).
				Prompt("").
				Value(&config.Border.Color).
				Placeholder("#515151"),
		),

		huh.NewGroup(
			huh.NewNote().Title("Shadow"),

			huh.NewInput().Title("Blur ").
				// Description("Shadow Gaussian Blur.").
				Placeholder("0").
				Inline(true).
				Prompt("").
				Value(&shadowBlur).
				Validate(validateInteger),

			huh.NewInput().Title("X Offset ").
				// Description("Shadow offset x coordinate").
				Placeholder("0").
				Inline(true).
				Prompt("").
				Value(&shadowX).
				Validate(validateInteger),

			huh.NewInput().Title("Y Offset ").
				// Description("Shadow offset y coordinate").
				Placeholder("0").
				Inline(true).
				Prompt("").
				Value(&shadowY).
				Validate(validateInteger),
		),
	).WithTheme(theme).WithHeight(10)

	err := f.Run()

	config.Padding = parsePadding(padding)
	config.Margin = parseMargin(margin)
	config.Font.Size, _ = strconv.ParseFloat(fontSize, 64)
	config.LineHeight, _ = strconv.ParseFloat(lineHeight, 64)
	config.Border.Radius, _ = strconv.Atoi(borderRadius)
	config.Border.Width, _ = strconv.Atoi(borderWidth)
	config.Shadow.Blur, _ = strconv.Atoi(shadowBlur)
	config.Shadow.X, _ = strconv.Atoi(shadowX)
	config.Shadow.Y, _ = strconv.Atoi(shadowY)
	return config, err
}

func validateMargin(s string) error {
	tokens := strings.Fields(s)
	if len(tokens) > 4 {
		return errors.New("maximum four values")
	}
	for _, t := range tokens {
		_, err := strconv.Atoi(t)
		if err != nil {
			return errors.New("must be valid space-separated integers")
		}
	}
	return nil
}

func validatePadding(s string) error {
	return validateMargin(s)
}

func validateInteger(s string) error {
	if len(s) <= 0 {
		return nil
	}

	_, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("must be valid integer")
	}
	return nil
}

func validateFloat(s string) error {
	if len(s) <= 0 {
		return nil
	}

	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.New("must be valid float")
	}
	return nil
}

var colorRegex = regexp.MustCompile("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$")

func validateColor(s string) error {
	if len(s) <= 0 {
		return nil
	}

	if !colorRegex.MatchString(s) {
		return errors.New("must be valid color")
	}
	return nil
}

func parsePadding(v string) []int {
	var values []int
	for _, p := range strings.Fields(v) {
		pi, _ := strconv.Atoi(p) // already validated
		values = append(values, pi)
	}
	return expandPadding(values)
}

var parseMargin = parsePadding
