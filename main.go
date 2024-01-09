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
	"github.com/beevik/etree"
	"github.com/charmbracelet/log"
)

var (
	output       string  = "out.svg"
	window       bool    = true
	outline      bool    = true
	shadow       bool    = true
	cornerRadius int     = 8
	padding      []int   = []int{20, 40, 20, 20}
	fontSize     float64 = 14
	lineHeight   float64 = fontSize * 1.2
)

var (
	red    string = "#FF5A54"
	yellow string = "#E6BF29"
	green  string = "#52C12B"

	grey string = "#515151"
)

const (
	top = iota
	right
	bottom
	left
)

func main() {
	var (
		input string
		err   error
		lexer chroma.Lexer
	)

	// Read input from file or stdin.
	if len(os.Args) > 1 {
		input, err = readFile(os.Args[1])
		lexer = lexers.Get(os.Args[1])
	} else if isPipe() {
		input, err = readInput(os.Stdin)
		lexer = lexers.Analyse(input)
	}
	if err != nil || input == "" {
		log.Fatal("no input provided.")
	}

	input = strings.TrimSpace(input)

	// Format code source.
	l := chroma.Coalesce(lexer)
	f := svg.New(svg.FontFamily("JetBrains Mono"))
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

	child := elements[0]
	w, h := getDimensions(child)

	rect := child.SelectElement("rect")
	w += padding[left] + padding[right]
	h += padding[top] + padding[bottom]
	setDimensions(rect, w, h)

	if window {
		addWindow(child)
	}

	if cornerRadius > 0 {
		addCornerRadius(rect)
	}

	if outline {
		addOutline(rect)

		// NOTE: necessary so that we don't clip the outline.
		setDimensions(child, w+1, h+1)
	} else {
		setDimensions(child, w, h)
	}

	lines := child.SelectElement("g").SelectElements("text")
	for i, line := range lines {
		// Offset the text by padding...
		// (x, y) -> (x+p, y+p)
		x := float64(padding[left])
		y := float64(i)*lineHeight + fontSize + float64(padding[top])
		move(line, x, y)
	}

	err = doc.WriteToFile(output)
	if err != nil {
		log.Fatal(err)
	}
}

// addCornerRadius adds corner radius to an element.
func addCornerRadius(element *etree.Element) {
	element.CreateAttr("rx", fmt.Sprintf("%d", cornerRadius))
	element.CreateAttr("ry", fmt.Sprintf("%d", cornerRadius))
}

// move moves the given element to the (x, y) position
func move(element *etree.Element, x, y float64) {
	element.SelectAttr("x").Value = fmt.Sprintf("%.2fpx", x)
	element.SelectAttr("y").Value = fmt.Sprintf("%.2fpx", y)
}

// addOutline adds an outline to the given element.
func addOutline(element *etree.Element) {
	element.CreateAttr("stroke", grey)
	element.CreateAttr("stroke-width", "1")
	element.CreateAttr("x", "0.5")
	element.CreateAttr("y", "0.5")
}

// addWindow adds a colorful window bar element to the given element.
func addWindow(element *etree.Element) {
	group := etree.NewElement("g")
	for i, c := range []string{red, yellow, green} {
		circle := etree.NewElement("circle")
		circle.CreateAttr("cx", fmt.Sprintf("%d", (i+1)*15))
		circle.CreateAttr("cy", "15")
		circle.CreateAttr("r", "5")
		circle.CreateAttr("fill", c)
		group.AddChild(circle)
	}
	element.AddChild(group)
	padding[top] += 15
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
