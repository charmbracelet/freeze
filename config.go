package main

type Configuration struct {
	Input string `arg:"" help:"the file to read" optional:""`

	Language     string  `help:"code language"`
	Output       string  `help:"output of the image" default:"out.svg"`
	Window       bool    `help:"whether to show window controls" default:"false"`
	Outline      bool    `help:"whether to add an outline to the window" default:"false"`
	Shadow       bool    `help:"whether to add a shadow to the window" default:"false"`
	CornerRadius int     `help:"amount to round the corners" default:"0"`
	Padding      []int   `help:"padding of the window" default:"20,40,20,20"`
	Margin       []int   `help:"margin of the window" default:"0"`
	FontFamily   string  `help:"font family" default:"JetBrains Mono"`
	FontSize     float64 `help:"font size" default:"14"`
	LineHeight   float64 `help:"line height" default:"16.8"`
}

func expandPadding(p []int) []int {
	switch len(p) {
	case 1:
		return []int{p[top], p[top], p[top], p[top]}
	case 2:
		return []int{p[top], p[right], p[top], p[right]}
	case 4:
		return []int{p[top], p[right], p[bottom], p[left]}
	default:
		return []int{0, 0, 0, 0}

	}
}

var expandMargin = expandPadding

type side int

const (
	top    side = 0
	right  side = 1
	bottom side = 2
	left   side = 3
)
