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
	output     string  = "out.svg"
	window     bool    = true
	shadow     bool    = true
	rounded    bool    = true
	padding    []int   = []int{30, 30}
	fontSize   float64 = 14
	lineHeight float64 = fontSize * 1.2
)

const (
	vertical   = 0
	horizontal = 1
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
		heightAttr.Value = fmt.Sprintf("%dpx", h+(2*padding[vertical]))
		widthAttr.Value = fmt.Sprintf("%dpx", w+(3*padding[horizontal]))

		textElements := child.SelectElement("g").SelectElements("text")
		for i, text := range textElements {
			// Offset the text by padding...
			// (x, y) -> (x+p, y+p)
			text.SelectAttr("x").Value = fmt.Sprintf("%dpx", padding[horizontal])
			text.SelectAttr("y").Value = fmt.Sprintf("%.2fpx", float64(i+1)*lineHeight+float64(padding[vertical]))
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
