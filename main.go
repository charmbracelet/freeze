package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/svg"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/alecthomas/kong"
	"github.com/beevik/etree"
	"github.com/charmbracelet/log"
	"github.com/mattn/go-isatty"
)

var (
	red    string = "#FF5A54"
	yellow string = "#E6BF29"
	green  string = "#52C12B"
)

func main() {
	var (
		input  string
		err    error
		lexer  chroma.Lexer
		config Config
	)

	_ = kong.Parse(&config,
		kong.Help(helpPrinter))

	c, ok := configs[config.Config]
	if !ok {
		b, err := os.ReadFile(config.Config)
		c = string(b)
		if err != nil {
			c = configs["base"]
		}
	}

	r, err := kong.JSON(strings.NewReader(c))
	if err != nil {
		log.Fatal("invalid json configuration", "error", err)
	}
	_ = kong.Parse(&config,
		kong.Resolvers(r),
		kong.Help(helpPrinter))

	config.Margin = expandMargin(config.Margin)
	config.Padding = expandPadding(config.Padding)

	if config.Input == "" && !isPipe(os.Stdin) {
		log.Fatal("no input provided.")
	}

	if config.Input == "-" || isPipe(os.Stdin) {
		input, err = readInput(os.Stdin)
		lexer = lexers.Analyse(input)
	} else {
		input, err = readFile(config.Input)
		lexer = lexers.Get(config.Input)
	}

	if config.Language != "" {
		lexer = lexers.Get(config.Language)
	}

	if lexer == nil {
		log.Fatal("unable to detect language, specify `--language`")

	}

	input = strings.TrimSpace(input)

	if err != nil || input == "" {
		log.Fatal("no input provided.")
	}

	// Format code source.
	l := chroma.Coalesce(lexer)
	ff := svg.EmbedFont("JetBrains Mono", FontJetBrainsMono, svg.WOFF2)
	if err != nil {
		log.Fatal(err)
	}
	f := svg.New(ff, svg.FontFamily(config.Font.Family))
	it, err := l.Tokenise(nil, input)
	if err != nil {
		log.Fatal(err)
	}
	buf := &bytes.Buffer{}

	s, ok := styles.Registry[strings.ToLower(config.Theme)]
	if s == nil || !ok {
		s = charmStyle
	}

	err = f.Format(buf, s, it)
	if err != nil {
		log.Fatal(err)
	}

	// Parse SVG (XML document)
	doc := etree.NewDocument()
	_, err = doc.ReadFrom(buf)
	if err != nil {
		log.Fatal(err)
	}

	elements := doc.ChildElements()
	if len(elements) < 1 {
		log.Fatal("no svg output")
	}

	svg := elements[0]
	w, h := getDimensions(svg)

	rect := svg.SelectElement("rect")
	w += config.Padding[left] + config.Padding[right]
	h += config.Padding[top] + config.Padding[bottom]
	setDimensions(rect, w, h)
	move(rect, float64(config.Margin[left]), float64(config.Margin[top]))

	if config.Window {
		config.addWindow(svg)
	}

	if config.Border.Radius > 0 {
		config.addCornerRadius(rect)
	}

	if config.Border.Width > 0 {
		config.addOutline(rect)

		if config.Margin[left] <= 0 && config.Margin[top] <= 0 {
			move(rect, float64(config.Border.Width)/2, float64(config.Border.Width)/2)
		}

		// NOTE: necessary so that we don't clip the outline.
		w += config.Border.Width
		h += config.Border.Width
	}

	setDimensions(svg, w+config.Margin[left]+config.Margin[right], h+config.Margin[top]+config.Margin[bottom])

	if config.Shadow.Blur > 0 || config.Shadow.X > 0 || config.Shadow.Y > 0 {
		id := "shadow"
		config.addShadow(svg, id)
		svg.CreateAttr("filter", fmt.Sprintf("url(#%s)", id))
	}

	lines := svg.SelectElement("g").SelectElements("text")
	for i, line := range lines {
		// Offset the text by padding...
		// (x, y) -> (x+p, y+p)
		x := float64(config.Padding[left] + config.Margin[left])
		y := float64(i)*(config.Font.Size*config.LineHeight) + config.Font.Size + float64(config.Padding[top]) + float64(config.Margin[top])
		move(line, x, y)
	}

	istty := isatty.IsTerminal(os.Stdout.Fd())

	switch {
	case strings.HasSuffix(config.Output, ".png"):
	case strings.HasSuffix(config.Output, ".svg"):
		if istty {
			err = doc.WriteToFile(config.Output)
		} else {
			_, err = doc.WriteTo(os.Stdout)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
}

// addShadow adds a definition of a shadow to the <defs> with the given id.
func (c *Config) addShadow(element *etree.Element, id string) {
	filter := etree.NewElement("filter")
	filter.CreateAttr("id", id)
	filter.CreateAttr("x", "0")
	filter.CreateAttr("y", "0")

	offset := etree.NewElement("feOffset")
	offset.CreateAttr("result", "offOut")
	offset.CreateAttr("in", "SourceAlpha")
	offset.CreateAttr("dx", fmt.Sprintf("%d", c.Shadow.X))
	offset.CreateAttr("dy", fmt.Sprintf("%d", c.Shadow.Y))

	color := etree.NewElement("feColorMatrix")
	color.CreateAttr("result", "matrixOut")
	color.CreateAttr("in", "offOut")
	color.CreateAttr("type", "matrix")
	color.CreateAttr("values", "0.2 0 0 0 0 0 0.2 0 0 0 0 0 0.2 0 0 0 0 0 1 0")

	blur := etree.NewElement("feGaussianBlur")
	blur.CreateAttr("result", "blurOut")
	blur.CreateAttr("in", "matrixOut")
	blur.CreateAttr("stdDeviation", fmt.Sprintf("%d", c.Shadow.Blur))

	blend := etree.NewElement("feBlend")
	blend.CreateAttr("in", "SourceGraphic")
	blend.CreateAttr("in2", "blurOut")
	blend.CreateAttr("mode", "normal")

	filter.AddChild(offset)
	filter.AddChild(blur)
	filter.AddChild(blend)

	defs := etree.NewElement("defs")
	defs.AddChild(filter)
	element.AddChild(defs)
}

// addCornerRadius adds corner radius to an element.
func (c *Config) addCornerRadius(element *etree.Element) {
	element.CreateAttr("rx", fmt.Sprintf("%d", c.Border.Radius))
	element.CreateAttr("ry", fmt.Sprintf("%d", c.Border.Radius))
}

// move moves the given element to the (x, y) position
func move(element *etree.Element, x, y float64) {
	element.CreateAttr("x", fmt.Sprintf("%.2fpx", x))
	element.CreateAttr("y", fmt.Sprintf("%.2fpx", y))
}

// addOutline adds an outline to the given element.
func (c *Config) addOutline(element *etree.Element) {
	element.CreateAttr("stroke", c.Border.Color)
	element.CreateAttr("stroke-width", fmt.Sprintf("%d", c.Border.Width))
}

// addWindow adds a colorful window bar element to the given element.
func (c *Config) addWindow(element *etree.Element) {
	group := etree.NewElement("g")
	for i, color := range []string{red, yellow, green} {
		circle := etree.NewElement("circle")
		circle.CreateAttr("cx", fmt.Sprintf("%d", (i+1)*19+c.Margin[left]-6))
		circle.CreateAttr("cy", fmt.Sprintf("%d", 12+c.Margin[top]))
		circle.CreateAttr("r", "5.5")
		circle.CreateAttr("fill", color)
		group.AddChild(circle)
	}
	element.AddChild(group)
	c.Padding[top] += 15
}

// setDimensions sets the width and height of the given element.
func setDimensions(element *etree.Element, width, height int) {
	widthAttr := element.SelectAttr("width")
	heightAttr := element.SelectAttr("height")
	heightAttr.Value = fmt.Sprintf("%d", height)
	widthAttr.Value = fmt.Sprintf("%d", width)
}

// getDimensions returns the width and height of the element.
func getDimensions(element *etree.Element) (int, int) {
	widthValue := element.SelectAttrValue("width", "0px")
	heightValue := element.SelectAttrValue("height", "0px")
	width := dimensionToInt(widthValue)
	height := dimensionToInt(heightValue)
	return width, height
}

// dimensionToInt takes a string and returns the integer value.
// e.g. "500px" -> 500
func dimensionToInt(px string) int {
	d := strings.TrimSuffix(px, "px")
	v, _ := strconv.ParseInt(d, 10, 64)
	return int(v)
}

// readFile returns the files content.
func readFile(file string) (string, error) {
	b, err := os.ReadFile(file)
	return string(b), err
}

// readInput reads some input.
func readInput(in io.Reader) (string, error) {
	b, err := io.ReadAll(in)
	return string(b), err
}

func isPipe(in *os.File) bool {
	stat, err := in.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}
