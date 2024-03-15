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
	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/mattn/go-isatty"
	"golang.org/x/net/context"
)

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

	// Copy the pty output to buffer
	if config.Execute != "" {
		input = executeCommand(config)
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

	var strippedInput string = ansi.Strip(input)
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

	hPadding := config.Padding[left] + config.Padding[right]
	hMargin := config.Margin[left] + config.Margin[right]
	vMargin := config.Margin[top] + config.Margin[bottom]
	vPadding := config.Padding[top] + config.Padding[bottom]

	autoHeight := config.Height == 0
	autoWidth := config.Width == 0

	terminal := image.SelectElement("rect")

	imageWidth, imageHeight := svg.GetDimensions(image)
	terminalWidth := imageWidth
	terminalHeight := imageHeight

	if !autoWidth {
		imageWidth = config.Width
		terminalWidth = config.Width - hMargin
	} else {
		imageWidth += hMargin + hPadding
		terminalWidth += hPadding
	}

	if !autoHeight {
		imageHeight = config.Height
		terminalHeight = config.Height - vMargin
	} else {
		imageHeight += vMargin + vPadding
		terminalHeight += vPadding
	}

	if config.Window {
		windowControls := svg.NewWindowControls(5.5, 19, 12)
		svg.Move(windowControls, float64(config.Margin[left]), float64(config.Margin[top]))
		image.AddChild(windowControls)
		config.Padding[top] += (15)
	}

	if config.Border.Radius > 0 {
		svg.AddCornerRadius(terminal, config.Border.Radius)
	}

	if config.Shadow.Blur > 0 || config.Shadow.X > 0 || config.Shadow.Y > 0 {
		id := "shadow"
		svg.AddShadow(image, id, config.Shadow.X, config.Shadow.Y, config.Shadow.Blur)
		terminal.CreateAttr("filter", fmt.Sprintf("url(#%s)", id))
	}

	textGroup := image.SelectElement("g")
	textGroup.CreateAttr("font-size", fmt.Sprintf("%.2fpx", config.Font.Size))
	textGroup.CreateAttr("clip-path", "url(#terminalMask)")
	text := textGroup.SelectElements("text")

	d := dispatcher{lines: text, svg: textGroup, config: &config}

	offsetLine := 0
	if len(config.Lines) > 0 {
		offsetLine = config.Lines[0]
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
			ln.SetText(fmt.Sprintf("%2d  ", i+1+offsetLine))
			line.InsertChildAt(0, ln)
		}
		x := float64(config.Padding[left] + config.Margin[left])
		y := (float64(i))*(config.Font.Size*config.LineHeight) + (config.Font.Size) + float64(config.Padding[top]) + float64(config.Margin[top])

		svg.Move(line, x, y)

		// We are passed visible lines, remove the rest.
		if y > float64(imageHeight-config.Margin[bottom]-config.Padding[bottom]) {
			textGroup.RemoveChild(line)
		}
	}

	if autoWidth {
		longestLine := lipgloss.Width(strippedInput)
		terminalWidth = int(float64(longestLine+1)*(config.Font.Size/fontHeightToWidthRatio)) + hPadding
		imageWidth = terminalWidth + hMargin
	}

	if config.Border.Width > 0 {
		svg.AddOutline(terminal, config.Border.Width, config.Border.Color)

		// NOTE: necessary so that we don't clip the outline.
		terminalHeight -= (config.Border.Width * 2)
		terminalWidth -= (config.Border.Width * 2)
	}

	if config.ShowLineNumbers {
		if autoWidth {
			terminalWidth += int(config.Font.Size * 2)
			imageWidth += int(config.Font.Size * 2)
		} else {
			terminalWidth -= int(config.Font.Size * 2)
		}
	}

	if !autoHeight || !autoWidth {
		svg.AddClipPath(image, "terminalMask",
			config.Margin[left], config.Margin[top],
			terminalWidth, terminalHeight-config.Padding[bottom])
	}

	svg.Move(terminal, max(float64(config.Margin[left]), float64(config.Border.Width)/2), max(float64(config.Margin[top]), float64(config.Border.Width)/2))
	svg.SetDimensions(image, imageWidth, imageHeight)
	svg.SetDimensions(terminal, terminalWidth, terminalHeight)

	if isAnsi {
		parser := ansi.Parser{
			Print:       d.Print,
			Execute:     d.Execute,
			CsiDispatch: d.CsiDispatch,
		}
		for _, line := range strings.Split(input, "\n") {
			parser.Parse([]byte(line))
			d.Execute(ansi.LF) // simulate a newline
		}
	}

	istty := isatty.IsTerminal(os.Stdout.Fd())

	switch {
	case strings.HasSuffix(config.Output, ".png"):
		// use libsvg conversion.
		svgConversionErr := libsvgConvert(doc, imageWidth, imageHeight, config.Output)
		if svgConversionErr == nil {
			printFilenameOutput(config.Output)
			break
		}

		// could not convert with libsvg, try resvg
		svgConversionErr = resvgConvert(doc, imageWidth, imageHeight, config.Output)
		if svgConversionErr != nil {
			printFilenameOutput(config.Output)
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

func executeCommand(config Config) string {
	args := strings.Split(config.Execute, " ")
	ctx, _ := context.WithTimeout(context.Background(), config.ExecuteTimeout)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	pty, err := config.runInPty(cmd)
	if err != nil {
		printErrorFatal("Something went wrong", err)
	}
	defer pty.Close()
	var out bytes.Buffer
	go func() { io.Copy(&out, pty) }()
	err = cmd.Wait()
	if err != nil {
		printError("Command failed", err)
	}
	return out.String()
}

var outputHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#875fff")).Bold(true).Padding(0, 1).MarginRight(1).SetString("WROTE")

func printFilenameOutput(filename string) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Center, outputHeader.String(), filename))
}
