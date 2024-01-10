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
	"github.com/alecthomas/kong"
	"github.com/beevik/etree"
	"github.com/charmbracelet/log"
)

var (
	red    string = "#FF5A54"
	yellow string = "#E6BF29"
	green  string = "#52C12B"
	grey   string = "#515151"
)

func main() {
	var (
		input  string
		err    error
		lexer  chroma.Lexer
		config Configuration
	)

	ctx := kong.Parse(&config)

	config.Margin = expandMargin(config.Margin)
	config.Padding = expandPadding(config.Padding)

	if config.Input == "" || config.Input == "-" {
		input, err = readInput(os.Stdin)
		lexer = lexers.Analyse(input)
	} else {
		input, err = readFile(ctx.Args[0])
		lexer = lexers.Get(ctx.Args[0])
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
	f := svg.New(svg.FontFamily(config.FontFamily), ff)
	it, err := l.Tokenise(nil, input)
	if err != nil {
		log.Fatal(err)
	}
	buf := &bytes.Buffer{}
	err = f.Format(buf, style, it)
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

	if config.CornerRadius > 0 {
		config.addCornerRadius(rect)
	}

	if config.Outline {
		addOutline(rect)

		if !config.Shadow {
			move(rect, 0.5, 0.5)
		}

		// NOTE: necessary so that we don't clip the outline.
		w += 1
		h += 1
	}

	setDimensions(svg, w, h)

	if config.Shadow {
		id := "shadow"
		addShadow(svg, id)
		svg.CreateAttr("filter", fmt.Sprintf("url(#%s)", id))
	}

	setDimensions(svg, w+config.Margin[left]+config.Margin[right], h+config.Margin[top]+config.Margin[bottom])

	lines := svg.SelectElement("g").SelectElements("text")
	for i, line := range lines {
		// Offset the text by padding...
		// (x, y) -> (x+p, y+p)
		x := float64(config.Padding[left] + config.Margin[left])
		y := float64(i)*config.LineHeight + config.FontSize + float64(config.Padding[top]) + float64(config.Margin[top])
		move(line, x, y)
	}

	err = doc.WriteToFile(config.Output)
	if err != nil {
		log.Fatal(err)
	}
}

// addShadow adds a definition of a shadow to the <defs> with the given id.
func addShadow(element *etree.Element, id string) {
	filter := etree.NewElement("filter")
	filter.CreateAttr("id", id)
	filter.CreateAttr("x", "0")
	filter.CreateAttr("y", "0")

	offset := etree.NewElement("feOffset")
	offset.CreateAttr("result", "offOut")
	offset.CreateAttr("in", "SourceAlpha")
	offset.CreateAttr("dx", "0")
	offset.CreateAttr("dy", "5")

	color := etree.NewElement("feColorMatrix")
	color.CreateAttr("result", "matrixOut")
	color.CreateAttr("in", "offOut")
	color.CreateAttr("type", "matrix")
	color.CreateAttr("values", "0.2 0 0 0 0 0 0.2 0 0 0 0 0 0.2 0 0 0 0 0 1 0")

	blur := etree.NewElement("feGaussianBlur")
	blur.CreateAttr("result", "blurOut")
	blur.CreateAttr("in", "matrixOut")
	blur.CreateAttr("stdDeviation", "12")

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
func (c *Configuration) addCornerRadius(element *etree.Element) {
	element.CreateAttr("rx", fmt.Sprintf("%d", c.CornerRadius))
	element.CreateAttr("ry", fmt.Sprintf("%d", c.CornerRadius))
}

// move moves the given element to the (x, y) position
func move(element *etree.Element, x, y float64) {
	element.CreateAttr("x", fmt.Sprintf("%.2fpx", x))
	element.CreateAttr("y", fmt.Sprintf("%.2fpx", y))
}

// addOutline adds an outline to the given element.
func addOutline(element *etree.Element) {
	element.CreateAttr("stroke", grey)
	element.CreateAttr("stroke-width", "1")
}

// addWindow adds a colorful window bar element to the given element.
func (c *Configuration) addWindow(element *etree.Element) {
	group := etree.NewElement("g")
	for i, color := range []string{red, yellow, green} {
		circle := etree.NewElement("circle")
		circle.CreateAttr("cx", fmt.Sprintf("%d", (i+1)*15+c.Margin[left]))
		circle.CreateAttr("cy", fmt.Sprintf("%d", 12+c.Margin[top]))
		circle.CreateAttr("r", "4.5")
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

// isPipe returns whether the stdin is piped.
func isPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
