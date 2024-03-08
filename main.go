package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma"
	formatter "github.com/alecthomas/chroma/formatters/svg"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/alecthomas/kong"
	"github.com/beevik/etree"
	in "github.com/charmbracelet/freeze/input"
	"github.com/charmbracelet/freeze/svg"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	parser "github.com/charmbracelet/x/exp/term/vtparser"
	"github.com/mattn/go-isatty"
	"github.com/rivo/uniseg"
)

const pngExportMultiplier = 2

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

	if config.Execute != "" {
		args := strings.Split(config.Execute, " ")
		cmd := exec.Command(args[0], args[1:]...)
		pty, err := runInPty(cmd)
		if err != nil {
			printErrorFatal("Something went wrong", err)
		}

		defer pty.Close()

		var out bytes.Buffer

		// Copy the pty output to buffer
		go func() {
			io.Copy(&out, pty)
		}()

		if err := cmd.Wait(); err != nil {
			printErrorFatal("Command failed", err)
		}

		input = out.String()
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
		cfg, interativeErr := runForm(&config)
		config = *cfg
		if interativeErr != nil {
			printErrorFatal("", interativeErr)
		}
		if isDefaultConfig {
			_ = saveUserConfig(*cfg)
		}
	}

	multiplier := 1
	if strings.HasSuffix(config.Output, ".png") {
		multiplier = pngExportMultiplier
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
	} else if config.Execute != "" {
		config.Language = "ansi"
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

	// adjust for 1-indexing
	for i := range config.Lines {
		config.Lines[i]--
	}

	strippedInput := StripANSI(input)
	isAnsi := strings.ToLower(config.Language) == "ansi" || strippedInput != input
	strippedInput = cut(strippedInput, config.Lines)

	if !isAnsi && lexer == nil {
		printErrorFatal("Language Unknown", errors.New("specify a language with the --language flag"))
	}

	input = cut(input, config.Lines)
	if input == "" {
		if err != nil {
			printErrorFatal("No input", err)
		} else {
			printErrorFatal("No input", errors.New("check --lines is within bounds"))
		}
	}

	s, ok := styles.Registry[strings.ToLower(config.Theme)]
	if s == nil || !ok {
		s = charmStyle
	}

	// Create a token iterator.
	var it chroma.Iterator
	if isAnsi {
		// For ANSI output, we'll inject our own SVG. For now, let's just strip the ANSI
		// codes and print the text to properly size the input.
		it = chroma.Literator(chroma.Token{Type: chroma.Text, Value: strippedInput})
	} else {
		it, err = chroma.Coalesce(lexer).Tokenise(nil, input)
		if err != nil {
			printErrorFatal("Could not lex file", err)
		}
	}

	// Format the code to an SVG.
	options, err := fontOptions(&config)
	if err != nil {
		printErrorFatal("Invalid font options", err)
	}

	f := formatter.New(options...)
	if err != nil {
		printErrorFatal("Malformed text", err)
	}

	buf := &bytes.Buffer{}
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
	if config.Width != 0 {
		w = config.Width
	}
	if config.Height != 0 {
		h = config.Height
	}

	rect := image.SelectElement("rect")

	config.Font.Size *= float64(multiplier)

	// apply multiplier
	if config.ShowLineNumbers {
		w += int(config.Font.Size * 3)
	}

	w *= multiplier
	h *= multiplier

	for i := range config.Padding {
		config.Padding[i] *= multiplier
	}
	for i := range config.Margin {
		config.Margin[i] *= multiplier
	}

	w += config.Padding[left] + config.Padding[right]
	h += config.Padding[top] + config.Padding[bottom]

	config.Shadow.Blur *= multiplier
	config.Shadow.X *= multiplier
	config.Shadow.Y *= multiplier

	config.Border.Radius *= multiplier
	config.Border.Width *= multiplier

	svg.SetDimensions(rect, w, h)
	svg.Move(rect, float64(config.Margin[left]), float64(config.Margin[top]))

	if config.Window {
		windowControls := svg.NewWindowControls(5.5*float64(multiplier), 19*multiplier, 12*multiplier)
		svg.Move(windowControls, float64(config.Margin[left]), float64(config.Margin[top]))
		image.AddChild(windowControls)
		config.Padding[top] += (15 * multiplier)
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
	svg.SetDimensions(rect, w+config.Margin[left]+config.Margin[right], h+config.Margin[top]+config.Margin[bottom])

	if config.Shadow.Blur > 0 || config.Shadow.X > 0 || config.Shadow.Y > 0 {
		id := "shadow"
		svg.AddShadow(image, id, config.Shadow.X, config.Shadow.Y, config.Shadow.Blur)
		rect.CreateAttr("filter", fmt.Sprintf("url(#%s)", id))
	}

	g := image.SelectElement("g")
	g.CreateAttr("font-size", fmt.Sprintf("%.2fpx", config.Font.Size))
	text := g.SelectElements("text")

	d := dispatcher{
		lines:  text,
		row:    0,
		svg:    g,
		config: &config,
	}

	for i, line := range text {
		if isAnsi {
			line.SetText("")
		}
		// Offset the text by padding...
		// (x, y) -> (x+p, y+p)
		if config.ShowLineNumbers {
			ln := etree.NewElement("tspan")
			ln.CreateAttr("xml:space", "preserve")
			ln.CreateAttr("fill", s.Get(chroma.LineNumbers).Colour.String())
			ln.SetText(fmt.Sprintf("%3d  ", i+1))
			line.InsertChildAt(0, ln)
		}
		x := float64(config.Padding[left] + config.Margin[left])
		y := float64(i)*(config.Font.Size*config.LineHeight) + config.Font.Size + float64(config.Padding[top]) + float64(config.Margin[top])
		svg.Move(line, x, y)
	}

	maxWidth := 0
	strippedInput = strings.ReplaceAll(strippedInput, "\t", "    ")
	for _, line := range strings.Split(strippedInput, "\n") {
		stringWidth := uniseg.StringWidth(line)
		if stringWidth > maxWidth {
			maxWidth = stringWidth
		}
	}

	textWidthPx := (float64(maxWidth+1) * (config.Font.Size / fontHeightToWidthRatio))
	if config.ShowLineNumbers {
		textWidthPx += (config.Font.Size * 3)
	}

	if config.Width != 0 {
		textWidthPx = float64(config.Width)
	}

	hPadding := float64(config.Padding[left] + config.Padding[right])
	hMargin := float64(config.Margin[left] + config.Margin[right])
	vPadding := float64(config.Padding[top] + config.Padding[bottom])
	vMargin := float64(config.Margin[top] + config.Margin[bottom])

	image.CreateAttr("width", fmt.Sprintf("%.2fpx", textWidthPx+hMargin+hPadding))
	rect.CreateAttr("width", fmt.Sprintf("%.2fpx", textWidthPx+hPadding))
	image.CreateAttr("height", fmt.Sprintf("%.2fpx", float64(h)+vMargin+vPadding))
	rect.CreateAttr("height", fmt.Sprintf("%.2fpx", float64(h)))

	if isAnsi {
		ansiParseErr := parser.New(&d).Parse(strings.NewReader(input))
		if ansiParseErr != nil {
			printErrorFatal("Could not parse ANSI", ansiParseErr)
		}
	}

	istty := isatty.IsTerminal(os.Stdout.Fd())

	switch {
	case strings.HasSuffix(config.Output, ".png"):
		// use libsvg conversion.
		svgConversionErr := libsvgConvert(doc, w, h, config.Output)
		if svgConversionErr == nil {
			break
		}

		// could not convert with libsvg, try resvg
		svgConversionErr = resvgConvert(doc, int(textWidthPx+hMargin+hPadding), h+int(vMargin), config.Output)
		if svgConversionErr != nil {
			printErrorFatal("Unable to convert SVG to PNG", svgConversionErr)
		}

	default:
		// output file specified.
		if config.Output != "" {
			err = doc.WriteToFile(config.Output)
			if err != nil {
				printErrorFatal("Unable to write output", err)
			}
			printFilenameOutput(config.Output)
			return
		}

		// reading from stdin.
		if config.Input == "" || config.Input == "-" {
			if istty {
				err = doc.WriteToFile(defaultOutputFilename)
				printFilenameOutput(defaultOutputFilename)
			} else {
				_, err = doc.WriteTo(os.Stdout)
			}
			if err != nil {
				printErrorFatal("Unable to write output", err)
			}
			return
		}

		// reading from file.
		if istty {
			config.Output = strings.TrimSuffix(filepath.Base(config.Input), filepath.Ext(config.Input)) + ".svg"
			err = doc.WriteToFile(config.Output)
			printFilenameOutput(config.Output)
		} else {
			_, err = doc.WriteTo(os.Stdout)
		}
		if err != nil {
			printErrorFatal("Unable to write output", err)
		}
	}
}

var outputHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#875fff")).Bold(true).Padding(0, 1).MarginRight(1).SetString("WROTE")

func printFilenameOutput(filename string) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Center, outputHeader.String(), filename))
}
