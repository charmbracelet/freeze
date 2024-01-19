package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

func runForm(config *Config) (*Config, error) {
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Output").
				Placeholder("out.svg").
				Description("Output location for image.").
				Value(&config.Output),

			huh.NewInput().Title("Theme").
				Placeholder("charm").
				Description("Theme for syntax highlighting").
				Suggestions([]string{"charm", "dracula", "catppuccin"}).
				Value(&config.Theme),
		).Title("Settings"),

		huh.NewGroup(
			huh.NewInput().Title("Background").
				Description("Apply a background fill.").
				Placeholder("#FFF").
				Value(&config.Background).
				Validate(validateColor),

			huh.NewInput().Title("Padding").
				Description("Apply padding to the code.").
				Placeholder("20 40").
				Validate(validatePadding),

			huh.NewInput().Title("Margin").
				Description("Apply margin to the window.").
				Placeholder("20").
				Validate(validatePadding),

			huh.NewConfirm().Title("Window Controls").
				Description("Display window controls.").
				Value(&config.Window),
		).Title("Window"),

		huh.NewGroup(
			huh.NewInput().Title("Font Family").
				Description("Font family to use for code").
				Placeholder("JetBrains Mono").
				Value(&config.Font.Family),

			huh.NewInput().Title("Font Size").
				Description("Font size to use for code.").
				Placeholder("14").
				Validate(validateInteger),

			huh.NewInput().Title("Line Height").
				Description("Line height relative to font size.").
				Placeholder("1.2").
				Validate(validateFloat),
		).Title("Window"),

		huh.NewGroup(
			huh.NewInput().Title("Border Radius").
				Description("Corner radius of the window.").
				Placeholder("0").
				Validate(validateInteger),

			huh.NewInput().Title("Border Width").
				Description("Border width thickness.").
				Placeholder("1").
				Validate(validateInteger),

			huh.NewInput().Title("Border color.").
				Validate(validateColor).
				Placeholder("#515151"),
		).Title("Shadow"),

		huh.NewGroup(
			huh.NewInput().Title("Blur").
				Description("Shadow Gaussian Blur.").
				Placeholder("0").
				Validate(validateInteger),

			huh.NewInput().Title("X Offset").
				Description("Shadow offset x coordinate").
				Placeholder("0").
				Validate(validateInteger),

			huh.NewInput().Title("Y Offset").
				Description("Shadow offset y coordinate").
				Placeholder("0").
				Validate(validateInteger),
		).Title("Shadow"),
	)
	err := f.Run()
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
	_, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("must be valid integer")
	}
	return nil
}

func validateFloat(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.New("must be valid float")
	}
	return nil
}

var colorRegex = regexp.MustCompile("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$")

func validateColor(s string) error {
	if !colorRegex.MatchString(s) {
		return errors.New("must be valid color")
	}
	return nil
}
