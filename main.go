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

	for _, child := range doc.ChildElements() {
		// Increase the width and height by double the padding.
		widthAttr := child.SelectAttr("width")
		heightAttr := child.SelectAttr("height")
		w := dimensionToInt(widthAttr.Value)
		h := dimensionToInt(heightAttr.Value)

		if window {
			padding[top] += 10
		}

		rect := child.SelectElement("rect")
		if cornerRadius > 0 {
			rect.CreateAttr("rx", fmt.Sprintf("%d", cornerRadius))
			rect.CreateAttr("ry", fmt.Sprintf("%d", cornerRadius))
		}

		if outline {
			rect.CreateAttr("stroke", "#515151")
			rect.CreateAttr("stroke-width", "1")
			rect.CreateAttr("x", "0.5")
			rect.CreateAttr("y", "0.5")

			rect.SelectAttr("height").Value = fmt.Sprintf("%d", h+(padding[top]+padding[bottom]))
			rect.SelectAttr("width").Value = fmt.Sprintf("%d", w+(padding[left]+padding[right]))
			w += 1
			h += 1
		}

		heightAttr.Value = fmt.Sprintf("%d", h+(padding[top]+padding[bottom]))
		widthAttr.Value = fmt.Sprintf("%d", w+(padding[left]+padding[right]))

		if window {
			circleGroup := etree.NewElement("g")
			for i, c := range []string{"#FF5A54", "#E6BF29", "#52C12B"} {
				circle := etree.NewElement("circle")
				circle.CreateAttr("cx", fmt.Sprintf("%d", (i+1)*15))
				circle.CreateAttr("cy", "14")
				circle.CreateAttr("r", "5")
				circle.CreateAttr("fill", c)
				circleGroup.AddChild(circle)
			}
			child.AddChild(circleGroup)
		}

		textElements := child.SelectElement("g").SelectElements("text")
		for i, text := range textElements {
			// Offset the text by padding...
			// (x, y) -> (x+p, y+p)
			text.SelectAttr("x").Value = fmt.Sprintf("%dpx", padding[left])
			text.SelectAttr("y").Value = fmt.Sprintf("%.2fpx", float64(i+1)*lineHeight+float64(padding[top]))
		}
	}

	err = doc.WriteToFile(output)
	if err != nil {
		log.Fatal(err)
	}
}

func dimensionToInt(px string) int {
	d := strings.TrimSuffix(px, "px")
	v, _ := strconv.ParseInt(d, 10, 64)
	return int(v)
}

func readFile(file string) (string, error) {
	b, err := os.ReadFile(file)
	return string(b), err
}

func readInput(in io.Reader) (string, error) {
	b, err := io.ReadAll(in)
	return string(b), err
}

func isPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
