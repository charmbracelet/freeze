package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2"
	formatter "github.com/alecthomas/chroma/v2/formatters/svg"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/alecthomas/kong"
	"github.com/beevik/etree"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-isatty"
	"github.com/muesli/reflow/wordwrap"

	in "github.com/charmbracelet/freeze/input"
	"github.com/charmbracelet/freeze/svg"
)

const (
	defaultFontSize   = 14.0
	defaultLineHeight = 1.2
)

var (
	// Version contains the application version number. It's set via ldflags
	// when building.
	Version = ""

	// CommitSHA contains the SHA of the commit that this application was built
	// against. It's set via ldflags when building.
	CommitSHA = ""
)

func main() {
	const shaLen = 7

	var (
		input  string
		err    error
		lexer  chroma.Lexer
		config Config
		scale  float64
	)

	k, err := kong.New(&config, kong.Help(helpPrinter))
	if err != nil {
		printErrorFatal("Something went wrong", err)
	}
	ctx, err := k.Parse(os.Args[1:])
	if err != nil || ctx.Error != nil {
		printErrorFatal("Invalid Usage", err)
	}

	if config.Version {
		info, ok := debug.ReadBuildInfo()
		if Version == "" && ok && info.Main.Sum != "" {
			Version = info.Main.Version
		} else {
			Version = "unknown (built from source)"
		}
		version := fmt.Sprintf("freeze version %s", Version)
		if len(CommitSHA) >= shaLen {
			version += " (" + CommitSHA[:shaLen] + ")"
		}
		fmt.Println(version)
		os.Exit(0)
	}

	// Copy the pty output to buffer
	if config.Execute != "" {
		input, err = executeCommand(config)
		if err != nil {
			if input != "" {
				err = fmt.Errorf("%w\n%s", err, input)
			}
			printErrorFatal("Something went wrong", err)
		}
		if input == "" {
			printErrorFatal("Something went wrong", errors.New("no command output"))
		}
	}

	isDefaultConfig := config.Config == "default"
	configFile, err := configs.Open("configurations/" + config.Config + ".json")
	if config.Config == "user" {
		configFile, err = loadUserConfig()
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
		cfg, interactiveErr := runForm(&config)
		config = *cfg
		if interactiveErr != nil {
			printErrorFatal("", interactiveErr)
		}
		if isDefaultConfig {
			_ = saveUserConfig(*cfg)
		}
	}

	autoHeight := config.Height == 0
	autoWidth := config.Width == 0

	if config.Output == "" {
		config.Output = defaultOutputFilename
	}

	scale = 1
	if autoHeight && autoWidth && strings.HasSuffix(config.Output, ".png") {
		scale = 4
	}

	config.Margin = expandMargin(config.Margin, scale)
	config.Padding = expandPadding(config.Padding, scale)

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
	highlightedLines := config.ComputeHighlightedLines()

	var strippedInput string = ansi.Strip(input)
	isAnsi := strings.ToLower(config.Language) == "ansi" || strippedInput != input
	strippedInput = cut(strippedInput, config.Lines)

	// wrap to character limit.
	if config.Wrap > 0 && !config.SoftWrap {
		strippedInput = wordwrap.String(strippedInput, config.Wrap)
		input = wordwrap.String(input, config.Wrap)
	}

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

	isRealLine := []bool{}
	strippedIsRealLine := []bool{}
	// wrap to character limit.
	if config.Wrap > 0 && config.SoftWrap {
		isRealLine = SoftWrap(input, config.Wrap)
		strippedIsRealLine = SoftWrap(strippedInput, config.Wrap)
		strippedInput = wordwrap.String(strippedInput, config.Wrap)
		input = wordwrap.String(input, config.Wrap)
	}

	if config.Wrap <= 0 {
		// If Wrap is disabled, but SoftWrap enabled, we force disable SoftWrap as it does not make sense
		// to keep this option enabled.
		printError("Wrap option disabled, but SoftWrap option enabled", fmt.Errorf("wrap option disabled"))
		config.SoftWrap = false
	}

	s, ok := styles.Registry[strings.ToLower(config.Theme)]
	if s == nil || !ok {
		s = charmStyle
	}
	if !s.Has(chroma.Background) {
		s, err = s.Builder().Add(chroma.Background, "bg:"+config.Background).Build()
		if err != nil {
			printErrorFatal("Could not add background", err)
		}
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

	terminal := image.SelectElement("rect")

	w, h := svg.GetDimensions(image)

	imageWidth := float64(w)
	imageHeight := float64(h)

	imageWidth *= scale
	imageHeight *= scale

	// chroma automatically calculates the height based on a font size of 14
	// and a line height of 1.2
	imageHeight *= (config.Font.Size / defaultFontSize)
	imageHeight *= (config.LineHeight / defaultLineHeight)

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
		windowControls := svg.NewWindowControls(5.5*float64(scale), 19.0*scale, 12.0*scale)
		svg.Move(windowControls, float64(config.Margin[left]), float64(config.Margin[top]))
		image.AddChild(windowControls)
		config.Padding[top] += (15 * scale)
	}

	if config.Border.Radius > 0 {
		svg.AddCornerRadius(terminal, config.Border.Radius*scale)
	}

	if config.Shadow.Blur > 0 || config.Shadow.X > 0 || config.Shadow.Y > 0 {
		id := "shadow"
		svg.AddShadow(image, id, config.Shadow.X*scale, config.Shadow.Y*scale, config.Shadow.Blur*scale)
		terminal.CreateAttr("filter", fmt.Sprintf("url(#%s)", id))
	}

	textGroup := image.SelectElement("g")
	textGroup.CreateAttr("font-size", fmt.Sprintf("%.2fpx", config.Font.Size*float64(scale)))
	textGroup.CreateAttr("clip-path", "url(#terminalMask)")
	text := textGroup.SelectElements("text")

	d := dispatcher{lines: text, svg: textGroup, config: &config, scale: scale}

	offsetLine := 0
	if len(config.Lines) > 0 {
		offsetLine = config.Lines[0]
	}

	config.LineHeight *= float64(scale)

	softWrapOffset := 0
	for i, line := range text {
		if isAnsi {
			line.SetText("")
		}
		// Offset the text by padding...
		// (x, y) -> (x+p, y+p)

		var bg *etree.Element
		doHighlightLine := highlightedLines[i+1+offsetLine-softWrapOffset]
		if config.SoftWrap {
			// If the current line is soft-wrapped, we need to find the previous real line
			if !isRealLine[i] {
				j := i
				// iterate previous lines until we find the previous real line
				for ; !isRealLine[j]; j-- {
				}
				// we apply to the current line the highlight status of the found real line
				doHighlightLine = highlightedLines[j+1+offsetLine-softWrapOffset]
			}
		}

		if config.ShowLineNumbers {
			ln := etree.NewElement("tspan")
			ln.CreateAttr("xml:space", "preserve")
			if doHighlightLine {
				ln.CreateAttr("fill", s.Get(chroma.LineHighlight).Colour.String())
			} else {
				ln.CreateAttr("fill", s.Get(chroma.LineNumbers).Colour.String())
			}
			if config.SoftWrap {
				if (isAnsi && strippedIsRealLine[i]) || (!isAnsi && isRealLine[i]) {
					ln.SetText(fmt.Sprintf("%3d  ", i+1+offsetLine-softWrapOffset))
				} else {
					ln.SetText("     ")
				}
			} else {
				ln.SetText(fmt.Sprintf("%3d  ", i+1+offsetLine))
			}
			ln.CreateAttr("height", strconv.Itoa(int(math.Round(config.Font.Size*config.LineHeight))))
			line.InsertChildAt(0, ln)
		}
		if config.SoftWrap && !((isAnsi && strippedIsRealLine[i]) || (!isAnsi && isRealLine[i])) {
			softWrapOffset++
		}
		x := float64(config.Padding[left] + config.Margin[left])
		// Rounding required to ensure that each line have the same height
		y := (float64(i+1))*math.Round(config.Font.Size*config.LineHeight) + float64(config.Padding[top]) + float64(config.Margin[top])
		if doHighlightLine {
			// Create a background element, with grey color and 50% opacity
			bg = etree.NewElement("rect")
			bg.CreateAttr("fill", "grey")
			bg.CreateAttr("fill-opacity", "0.5")
			// This lineWidth is not accurate when the width is dynamic, it will be computed later
			lineWidth := imageWidth + config.Margin[left] + config.Padding[left] + config.Margin[right] + config.Padding[right]
			bg.CreateAttr("width", strconv.Itoa(int(lineWidth)))
			// We round to ensure that two highlighted consecutive lines do not leave a one pixel line between
			bg.CreateAttr("height", strconv.Itoa(int(math.Round(config.Font.Size*config.LineHeight))))
			line.Parent().InsertChildAt(0, bg)

			yRect := float64(i)*math.Round(config.Font.Size*config.LineHeight) +
				config.Padding[top] +
				config.Margin[top] + math.Round(config.LineHeight*config.Font.Size)/4
			// We round to ensure that two highlighted consecutive lines do not leave a one pixel line between
			svg.Move(bg, 0, math.Round(yRect))
		}
		svg.Move(line, x, y)

		// We are passed visible lines, remove the rest.
		if y > float64(imageHeight-config.Margin[bottom]-config.Padding[bottom]) {
			textGroup.RemoveChild(line)
		}
	}

	if autoWidth {
		tabWidth := 4
		if isAnsi {
			tabWidth = 6
		}
		longestLine := lipgloss.Width(strings.ReplaceAll(strippedInput, "\t", strings.Repeat(" ", tabWidth)))
		terminalWidth = float64(longestLine+1) * (config.Font.Size / fontHeightToWidthRatio)
		terminalWidth *= scale
		terminalWidth += hPadding
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
			terminalWidth += config.Font.Size * 3 * scale
			imageWidth += config.Font.Size * 3 * scale
		} else {
			terminalWidth -= config.Font.Size * 3
		}
	}

	// Adjust the highlighted rect width with the accurate computed width
	if len(highlightedLines) != 0 {
		for _, elem := range textGroup.SelectElements("rect") {
			elem.CreateAttr("width", strconv.Itoa(int(imageWidth)))
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
		parser := ansi.NewParser()
		parser.SetHandler(ansi.Handler{
			Print:     d.Print,
			HandleCsi: d.CsiDispatch,
			Execute:   d.Execute,
		})
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
			printErrorFatal("Unable to convert SVG to PNG", svgConversionErr)
		}
		printFilenameOutput(config.Output)

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

var outputHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#6C50FF")).Bold(true).Padding(0, 1).MarginRight(1).SetString("WROTE")

func printFilenameOutput(filename string) {
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Center, outputHeader.String(), filename))
}
