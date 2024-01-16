package main

type Config struct {
	Input    string `arg:"" help:"code file to screenshot" optional:""`
	Config   string `help:"screenshot configuration" short:"c" group:"Settings" default:"base"`
	Output   string `help:"output of the image" short:"o" group:"Settings" default:"out.svg"`
	Language string `help:"code language" short:"l" group:"Settings"`
	Theme    string `help:"theme" short:"t" group:"Settings"`

	Window     bool   `help:"show window controls" short:"w"`
	Border     Border `embed:"" prefix:"border." group:"Border"`
	Shadow     Shadow `embed:"" prefix:"shadow." help:"add a shadow to the window" short:"s" group:"shadow"`
	Padding    []int  `help:"terminal padding" short:"p"`
	Margin     []int  `help:"window margin" short:"m"`
	Background string `help:"background fill" short:"b"`

	Font       Font    `embed:"" prefix:"font." group:"font"`
	LineHeight float64 `group:"font"`
}

type Shadow struct {
	Blur int `help:"shadow blur"`
	X    int `help:"x offset"`
	Y    int `help:"y offset"`
}

type Border struct {
	Radius int    `help:"corner radius" short:"r"`
	Width  int    `help:"border width"`
	Color  string `help:"border color"`
}

type Font struct {
	Family string
	Size   float64
}

var configs = map[string]string{
	"base": `{
	"window": false,
	"border": {
		"radius": 0,
		"width": 0,
		"color": "#515151"
	},
	"shadow": {
		"blur": 0,
		"x": 0,
		"y": 0
	},
	"padding": [20, 40, 20, 20],
	"margin": "0",
	"background": "#FFFFFF",
	"font": {
		"family": "JetBrains Mono",
		"size": 14
	},
	"line_height": 1.2
}`,
	"full": `{
	"window": true,
	"border": {
		"radius": 8,
		"width": 1,
		"color": "#515151"
	},
	"shadow": {
		"blur": 24,
		"x": 0,
		"y": 12
	},
	"padding": [20, 40, 20, 20],
	"margin": [50, 60, 100, 60],
	"background": "#FFFFFF",
	"font": {
		"family": "JetBrains Mono",
		"size": 14
	},
	"line_height": 1.2
}`,
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
