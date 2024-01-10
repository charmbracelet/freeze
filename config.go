package main

type Configuration struct {
	Input string `arg:"" help:"code file to screenshot" optional:""`

	Language string `help:"code language" short:"l"`
	Output   string `help:"output of the image" short:"o" default:"out.svg"`
	Window   bool   `help:"show window controls" short:"w" default:"false"`
	Border   bool   `help:"add an outline to the window" short:"b" default:"false"`
	Shadow   bool   `help:"add a shadow to the window" short:"s" default:"false"`
	Radius   int    `help:"corner radius" short:"r" default:"0"`
	Padding  []int  `help:"terminal padding" short:"p" default:"20,40,20,20"`
	Margin   []int  `help:"window margin" short:"m" default:"0"`
	Font     struct {
		Family string  `default:"JetBrains Mono"`
		Size   float64 `default:"14"`
	} `embed:"" prefix:"font."`
	LineHeight float64 `default:"16.8"`
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
