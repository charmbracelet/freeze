package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma"
	formatter "github.com/alecthomas/chroma/formatters/svg"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/alecthomas/kong"
	"github.com/beevik/etree"
	"github.com/charmbracelet/freeze/font"
	in "github.com/charmbracelet/freeze/input"
	"github.com/charmbracelet/freeze/svg"
	"github.com/charmbracelet/log"
	"github.com/kanrichan/resvg-go"
	"github.com/mattn/go-isatty"
)

const pngExportMultiplier = 3

func main() {
	var (
		input  string
		err    error
		lexer  chroma.Lexer
		config Config
	)

	k, err := kong.New(&config, kong.Help(helpPrinter))
	if err != nil {
		printErrorFatal("Something went wrong", err)
	}
	ctx, err := k.Parse(os.Args[1:])
	if err != nil || ctx.Error != nil {
		printErrorFatal("Invalid Usage", err)
	}

	isDefaultConfig := config.Config == "default"

	configFile, err := loadUserConfig()
	if err != nil || !isDefaultConfig {
		configFile, err = configs.Open("configurations/" + config.Config + ".json")
	}
	if err != nil {
		configFile, err = os.Open(config.Config)
	}
	if err != nil {
		configFile, _ = configs.Open("configurations/base.json")
	}
	r, err := kong.JSON(configFile)
	if err != nil {
		printErrorFatal("Invalid JSON", err)
	}
	k, err = kong.New(&config, kong.Help(helpPrinter), kong.Resolvers(r))
	if err != nil {
		printErrorFatal("Something went wrong", err)
	}
	ctx, err = k.Parse(os.Args[1:])
	if err != nil {
		printErrorFatal("Invalid Usage", err)
	}

	if config.Interactive {
		cfg, err := runForm(&config)
		config = *cfg
		if err != nil {
			printErrorFatal("", err)
		}
		if isDefaultConfig {
			_ = saveUserConfig(*cfg)
		}
	}

	config.Margin = expandMargin(config.Margin)
	config.Padding = expandPadding(config.Padding)

	if config.Input == "" && !in.IsPipe(os.Stdin) && len(ctx.Args) <= 0 {
		_ = helpPrinter(kong.HelpOptions{}, ctx)
		os.Exit(0)
	}

	if config.Input == "-" || in.IsPipe(os.Stdin) {
		input, err = in.ReadInput(os.Stdin)
		lexer = lexers.Analyse(input)
	} else {
		input, err = in.ReadFile(config.Input)
		if err != nil {
			printErrorFatal("File not found", err)
		}
		lexer = lexers.Get(config.Input)
	}

	if config.Language != "" {
		lexer = lexers.Get(config.Language)
	}

	if lexer == nil {
		printErrorFatal("Language Unknown", errors.New("specify a language with the --language flag"))
	}

	input = strings.TrimSpace(input)

	if err != nil || input == "" {
		printErrorFatal("No input", err)
	}

	// Format code source.
	l := chroma.Coalesce(lexer)
	ff := formatter.EmbedFont("JetBrains Mono", font.JetBrainsMono, formatter.WOFF2)
	f := formatter.New(ff, formatter.FontFamily(config.Font.Family))
	it, err := l.Tokenise(nil, input)
	if err != nil {
		printErrorFatal("Malformed text", err)
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
		printErrorFatal("Bad SVG", err)
	}

	elements := doc.ChildElements()
	if len(elements) < 1 {
		printErrorFatal("Bad Output", nil)
	}

	image := elements[0]
	w, h := svg.GetDimensions(image)

	rect := image.SelectElement("rect")
	w += config.Padding[left] + config.Padding[right]
	h += config.Padding[top] + config.Padding[bottom]
	svg.SetDimensions(rect, w, h)
	svg.Move(rect, float64(config.Margin[left]), float64(config.Margin[top]))

	if config.Window {
		windowControls := svg.NewWindowControls()
		svg.Move(windowControls, float64(config.Margin[left]), float64(config.Margin[top]))
		image.AddChild(windowControls)
		config.Padding[top] += 15
	}

	if config.Border.Radius > 0 {
		svg.AddCornerRadius(rect, config.Border.Radius)
	}

	if config.Border.Width > 0 {
		svg.AddOutline(rect, config.Border.Width, config.Border.Color)

		if config.Margin[left] <= 0 && config.Margin[top] <= 0 {
			svg.Move(rect, float64(config.Border.Width)/2, float64(config.Border.Width)/2)
		}

		// NOTE: necessary so that we don't clip the outline.
		w += config.Border.Width
		h += config.Border.Width
	}

	svg.SetDimensions(image, w+config.Margin[left]+config.Margin[right], h+config.Margin[top]+config.Margin[bottom])

	if config.Shadow.Blur > 0 || config.Shadow.X > 0 || config.Shadow.Y > 0 {
		id := "shadow"
		svg.AddShadow(image, id, config.Shadow.X, config.Shadow.Y, config.Shadow.Blur)
		rect.CreateAttr("filter", fmt.Sprintf("url(#%s)", id))
	}

	lines := image.SelectElement("g").SelectElements("text")
	for i, line := range lines {
		// Offset the text by padding...
		// (x, y) -> (x+p, y+p)
		x := float64(config.Padding[left] + config.Margin[left])
		y := float64(i)*(config.Font.Size*config.LineHeight) + config.Font.Size + float64(config.Padding[top]) + float64(config.Margin[top])
		svg.Move(line, x, y)
	}

	istty := isatty.IsTerminal(os.Stdout.Fd())

	switch {
	case strings.HasSuffix(config.Output, ".png"):
		svg, err := doc.WriteToBytes()
		if err != nil {
			printErrorFatal("Unable to write output", err)
		}
		worker, err := resvg.NewDefaultWorker(context.Background())
		defer worker.Close()
		if err != nil {
			printErrorFatal("Unable to write output", err)
		}

		fontdb, err := worker.NewFontDBDefault()
		defer fontdb.Close()
		if err != nil {
			printErrorFatal("Unable to write output", err)
		}
		fontdb.LoadFontData(font.JetBrainsMonoTTF)

		pixmap, err := worker.NewPixmap(uint32(w+config.Margin[left]+config.Margin[right]), uint32(h+config.Margin[top]+config.Margin[bottom]))
		defer pixmap.Close()
		if err != nil {
			printErrorFatal("Unable to write output", err)
		}

		tree, err := worker.NewTreeFromData(svg, &resvg.Options{})
		defer tree.Close()
		if err != nil {
			printErrorFatal("Unable to write output", err)
		}

		err = tree.ConvertText(fontdb)
		if err != nil {
			printErrorFatal("Unable to write output", err)
		}
		tree.Render(resvg.TransformIdentity(), pixmap)
		png, err := pixmap.EncodePNG()
		if err != nil {
			printErrorFatal("Unable to write output", err)
		}

		os.WriteFile(config.Output, png, 0644)

	case strings.HasSuffix(config.Output, ".svg"):
		if istty {
			err = doc.WriteToFile(config.Output)
		} else {
			_, err = doc.WriteTo(os.Stdout)
		}
	}
	if err != nil {
		printErrorFatal("Unable to write output", err)
	}
}
