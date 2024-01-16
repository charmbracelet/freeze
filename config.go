package main

type Config struct {
	Input    string `arg:"" help:"code file to screenshot" optional:""`
	Config   string `help:"Base configuration file or template." short:"c" group:"Settings" default:"base"`
	Output   string `help:"Output location for SVG, PNG, or JPEG." short:"o" group:"Settings" default:"out.svg"`
	Language string `help:"Language of code file." short:"l" group:"Settings"`
	Theme    string `help:"Theme to use for syntax highlighting." short:"t" group:"Settings"`

	Window     bool   `help:"Display window controls." short:"w"`
	Border     Border `embed:"" prefix:"border." group:"Border"`
	Shadow     Shadow `embed:"" prefix:"shadow." help:"add a shadow to the window" short:"s" group:"Shadow"`
	Padding    []int  `help:"Apply padding to the code." short:"p"`
	Margin     []int  `help:"Apply margin to the window." short:"m"`
	Background string `help:"Apply a background fill." short:"b"`

	Font       Font    `embed:"" prefix:"font." group:"Font"`
	LineHeight float64 `help:"Line height relative to font size." group:"font"`
}

type Shadow struct {
	Blur int `help:"Shadow Gaussian Blur."`
	X    int `help:"Shadow offset X coordinate"`
	Y    int `help:"Shadow offset Y coordinate"`
}

type Border struct {
	Radius int    `help:"Cornder radius of window." short:"r"`
	Width  int    `help:"Border width thickness."`
	Color  string `help:"Border color."`
}

type Font struct {
	Family string  `help:"Font family to use for code."`
	Size   float64 `help:"Font size to use for code."`
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
