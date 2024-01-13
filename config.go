package main

type Configuration struct {
	Input string `arg:"" help:"code file to screenshot" optional:""`

	Language string `help:"code language" short:"l"`
	Theme    string `help:"theme" short:"t"`
	Output   string `help:"output of the image" short:"o" default:"out.svg"`
	Window   bool   `help:"show window controls" short:"w" default:"false"`
	Border   struct {
		Radius int    `help:"corner radius" short:"r" default:"0"`
		Width  int    `help:"border width" default:"0"`
		Color  string `help:"border color" default:"#515151"`
	} `embed:"" prefix:"border."`
	Shadow  bool  `help:"add a shadow to the window" short:"s" default:"false"`
	Padding []int `help:"terminal padding" short:"p" default:"20,40,20,20"`
	Margin  []int `help:"window margin" short:"m" default:"0"`
	Font    struct {
		Family string  `default:"JetBrains Mono"`
		Size   float64 `default:"14"`
	} `embed:"" prefix:"font."`
	LineHeight float64 `default:"1.2"`
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
