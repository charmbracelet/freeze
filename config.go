package main

type Config struct {
	Input string `arg:"" help:"Code to screenshot." optional:""`

	// Window
	Background string `help:"Apply a background fill." short:"b" placeholder:"#FFF" group:"Window"`
	Margin     []int  `help:"Apply margin to the window." short:"m" placeholder:"0" group:"Window"`
	Padding    []int  `help:"Apply padding to the code." short:"p" placeholder:"0" group:"Window"`
	Window     bool   `help:"Display window controls." short:"w" group:"Window"`

	// Settings
	Config   string `help:"Base configuration file or template." short:"c" group:"Settings" default:"base" placeholder:"base"`
	Language string `help:"Language of code file." short:"l" group:"Settings" placeholder:"go"`
	Output   string `help:"Output location for {{.svg}}, {{.png}}, or {{.jpeg}}." short:"o" group:"Settings" default:"out.svg" placeholder:"out.svg"`
	Theme    string `help:"Theme to use for syntax highlighting." short:"t" group:"Settings" placeholder:"charm"`

	Border Border `embed:"" prefix:"border." group:"Border"`
	Shadow Shadow `embed:"" prefix:"shadow." help:"add a shadow to the window" short:"s" group:"Shadow"`

	Font       Font    `embed:"" prefix:"font." group:"Font"`
	LineHeight float64 `help:"Line height relative to font size." group:"Font" placeholder:"1.2"`
}

type Shadow struct {
	Blur int `help:"Shadow Gaussian Blur." placeholder:"0"`
	X    int `help:"Shadow offset {{x}} coordinate" placeholder:"0"`
	Y    int `help:"Shadow offset {{y}} coordinate" placeholder:"0"`
}

type Border struct {
	Radius int    `help:"Corner radius of window." short:"r" placeholder:"0"`
	Width  int    `help:"Border width thickness." placeholder:"1"`
	Color  string `help:"Border color." placeholder:"#000"`
}

type Font struct {
	Family string  `help:"Font family to use for code." placeholder:"monospace"`
	Size   float64 `help:"Font size to use for code." placeholder:"14"`
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
