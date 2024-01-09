package main

import (
	"bytes"
	"io"
	"os"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/svg"
	"github.com/alecthomas/chroma/lexers"
	"github.com/beevik/etree"
	"github.com/charmbracelet/log"
)

var (
	output  string = "out.svg"
	window  bool   = true
	shadow  bool   = true
	rounded bool   = true
	padding int    = 16
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
	} else {
		input, err = readInput(os.Stdin)
		lexer = lexers.Analyse(input)
	}
	if err != nil || input == "" {
		log.Fatal("no input provided.")
	}

	// Format code source.
	l := chroma.Coalesce(lexer)
	f := svg.New(svg.FontFamily("JetBrains Mono"))
	it, err := l.Tokenise(nil, input)
	if err != nil {
		log.Fatal(err)
	}
	buf := &bytes.Buffer{}
	err = f.Format(buf, style(), it)
	if err != nil {
		log.Fatal(err)
	}

	doc := etree.NewDocument()
	_, err = doc.ReadFrom(buf)
	if err != nil {
		log.Fatal(err)
	}
	doc.ChildElements()[0].Attr[0].Value = "39px"

	doc.WriteToFile(output)
}

func readFile(file string) (string, error) {
	b, err := os.ReadFile(file)
	return string(b), err
}

func readInput(in io.Reader) (string, error) {
	b, err := io.ReadAll(in)
	return string(b), err
}
