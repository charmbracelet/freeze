package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/acarl005/stripansi"
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
	"github.com/charmbracelet/term/ansi/parser"
	"github.com/kanrichan/resvg-go"
	"github.com/mattn/go-isatty"
	"github.com/mattn/go-runewidth"
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

	strippedInput := stripansi.Strip(input)
	isAnsi := strings.ToLower(config.Language) == "ansi" || strippedInput != input

	if !isAnsi && lexer == nil {
		printErrorFatal("Language Unknown", errors.New("specify a language with the --language flag"))
	}

	// adjust for 1-indexing
	for i := range config.Lines {
		config.Lines[i]--
	}

	input = cut(input, config.Lines)
	if err != nil || input == "" {
		printErrorFatal("No input", err)
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
	}

	// Format the code to an SVG.
	ff := formatter.EmbedFont("JetBrains Mono", font.JetBrainsMono, formatter.WOFF2)
	f := formatter.New(ff, formatter.FontFamily(config.Font.Family))
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

	rect := image.SelectElement("rect")

	// apply multiplier
	w *= multiplier
	h *= multiplier
	config.Font.Size *= float64(multiplier)

	if config.ShowLineNumbers {
		w += int(config.Font.Size * 3)
	}

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
			ln.CreateAttr("fill", s.Get(chroma.LineNumbers).Colour.String())
			ln.SetText(fmt.Sprintf("%3d  ", i+1))
			line.InsertChildAt(0, ln)
		}
		x := float64(config.Padding[left] + config.Margin[left])
		y := float64(i)*(config.Font.Size*config.LineHeight) + config.Font.Size + float64(config.Padding[top]) + float64(config.Margin[top])
		svg.Move(line, x, y)
	}

	if isAnsi {
		parser.New(&d).Parse(strings.NewReader(input))
	}

	istty := isatty.IsTerminal(os.Stdout.Fd())

	switch {
	case strings.HasSuffix(config.Output, ".png"):
		svg, err := doc.WriteToBytes()
		if err != nil {
			printErrorFatal("Unable to write output", err)
		}

		if _, err := exec.LookPath("rsvg-convert"); err == nil {
			// rsvg-convert is installed use that to convert the SVG to PNG,
			// since it is faster.
			rsvgConvert := exec.Command("rsvg-convert",
				"--width", strconv.Itoa(w),
				"--keep-aspect-ratio",
				"-f", "png",
				"-o", config.Output,
			)
			rsvgConvert.Stdin = bytes.NewReader(svg)
			err = rsvgConvert.Run()
			if err == nil {
				break
			}

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

		pixmap, err := worker.NewPixmap(uint32((w + config.Margin[left] + config.Margin[right])), uint32(h+config.Margin[top]+config.Margin[bottom]))
		defer pixmap.Close()
		if err != nil {
			printError("Unable to write output", err)
			os.Exit(1)
		}

		tree, err := worker.NewTreeFromData(svg, &resvg.Options{
			Dpi:                288.0,
			ShapeRenderingMode: resvg.ShapeRenderingModeGeometricPrecision,
			TextRenderingMode:  resvg.TextRenderingModeGeometricPrecision,
			ImageRenderingMode: resvg.ImageRenderingModeOptimizeQuality,
			FontSize:           float32(config.Font.Size),
		})
		defer tree.Close()
		if err != nil {
			printError("Unable to write output", err)
			os.Exit(1)
		}

		err = tree.ConvertText(fontdb)
		if err != nil {
			printError("Unable to render text", err)
			os.Exit(1)
		}
		err = tree.Render(resvg.TransformIdentity(), pixmap)
		if err != nil {
			printError("Unable to render PNG", err)
			os.Exit(1)
		}
		png, err := pixmap.EncodePNG()
		if err != nil {
			printError("Unable to encode PNG", err)
			os.Exit(1)
		}

		err = os.WriteFile(config.Output, png, 0644)
		if err != nil {
			printError("Unable to write output", err)
			os.Exit(1)
		}

	case strings.HasSuffix(config.Output, ".svg"):
		if istty {
			err = doc.WriteToFile(config.Output)
		} else {
			_, err = doc.WriteTo(os.Stdout)
		}
		if err != nil {
			printErrorFatal("Unable to write output", err)
		}
	}
}

type dispatcher struct {
	lines           []*etree.Element
	row             int
	col             int
	svg             *etree.Element
	background      *etree.Element
	backgroundWidth int
	config          *Config
}

func (p *dispatcher) Print(r rune) {
	// insert the rune in the last tspan
	children := p.lines[p.row].ChildElements()
	var lastChild *etree.Element
	if len(children) == 0 {
		lastChild = etree.NewElement("tspan")
		lastChild.CreateAttr("xml:space", "preserve")
		p.lines[p.row].AddChild(lastChild)
	} else {
		lastChild = children[len(children)-1]
	}
	lastChild.SetText(lastChild.Text() + string(r))
	p.col += runewidth.RuneWidth(r)
	if p.background != nil {
		p.backgroundWidth += runewidth.RuneWidth(r)
	}
}

func (p *dispatcher) Execute(code byte) {
	if code == 0x0A {
		p.row++
		p.col = 0
		p.endBackground()
	}
}
func (p *dispatcher) DcsPut(code byte) {}
func (p *dispatcher) DcsUnhook()       {}

func (p *dispatcher) OscDispatch(params [][]byte, bellTerminated bool)      {}
func (p *dispatcher) EscDispatch(intermediates []byte, r rune, ignore bool) {}
func (p *dispatcher) DcsHook(prefix string, params [][]uint16, intermediates []byte, r rune, ignore bool) {
}

func (p *dispatcher) beginBackground(fill string) {
	rect := etree.NewElement("rect")
	rect.CreateAttr("fill", fill)
	rect.CreateAttr("x", fmt.Sprintf("%.2fpx", (float64(p.col)*p.config.Font.Size/1.658)+float64(p.config.Margin[left]+p.config.Padding[left])))
	rect.CreateAttr("y", fmt.Sprintf("%.2fpx", float64(p.row)*p.config.Font.Size*p.config.LineHeight+float64(p.config.Margin[top]+p.config.Padding[top])))
	rect.CreateAttr("height", fmt.Sprintf("%.2fpx", p.config.Font.Size*p.config.LineHeight+1))
	p.background = rect
}

func (p *dispatcher) endBackground() {
	if p.background == nil {
		return
	}

	p.background.CreateAttr("width", fmt.Sprintf("%.2fpx", float64(p.backgroundWidth)*p.config.Font.Size/1.658+1))
	p.svg.InsertChildAt(0, p.background)
	p.background = nil
	p.backgroundWidth = 0
}

func (p *dispatcher) CsiDispatch(prefix string, params [][]uint16, intermediates []byte, r rune, ignore bool) {
	span := etree.NewElement("tspan")
	span.CreateAttr("xml:space", "preserve")

	var i int
	for i < len(params) {
		v := params[i][0]
		switch v {
		case 0:
			// reset ANSI, this is done by creating a new empty tspan,
			// which would reset all the styles such that when text is appended to the last
			// child of this line there is no styling applied.
			p.lines[p.row].AddChild(span)
			p.endBackground()
		case 1:
			span.CreateAttr("font-weight", "bold")
			p.lines[p.row].AddChild(span)
		case 9:
			span.CreateAttr("text-decoration", "line-through")
			p.lines[p.row].AddChild(span)
		case 3:
			span.CreateAttr("font-style", "italic")
			p.lines[p.row].AddChild(span)
		case 4:
			span.CreateAttr("text-decoration", "underline")
			p.lines[p.row].AddChild(span)
		case 30, 31, 32, 33, 34, 35, 36, 37, 90, 91, 92, 93, 94, 95, 96, 97:
			span.CreateAttr("fill", ansi[v])
			p.lines[p.row].AddChild(span)
		case 38:
			i++
			switch params[i][0] {
			case 5:
				span.CreateAttr("fill", fmt.Sprintf("#%02x%02x%02x", params[i+1][0], params[i+1][0], params[i+1][0]))
				p.lines[p.row].AddChild(span)
				i++
			case 2:
				span.CreateAttr("fill", fmt.Sprintf("#%02x%02x%02x", params[i+1][0], params[i+2][0], params[i+3][0]))
				p.lines[p.row].AddChild(span)
				i += 3
			}
		case 48:
			i++
			switch params[i][0] {
			case 5:
				fill := fmt.Sprintf("#%02x%02x%02x", params[i+1][0], params[i+1][0], params[i+1][0])
				p.beginBackground(fill)
				i++
			case 2:
				fill := fmt.Sprintf("#%02x%02x%02x", params[i+1][0], params[i+2][0], params[i+3][0])
				p.beginBackground(fill)
				i += 3
			}
		case 100, 101, 102, 103, 104, 105, 106, 107:
			p.beginBackground(ansi[v])
		}
		i++
	}
}

var ansi = map[uint16]string{
	30: "#676767", // black
	31: "#ff5f87", // red
	32: "#00dc7f", // green
	33: "#e8e8a8", // yellow
	34: "#00aaff", // blue
	35: "#635adf", // magenta
	36: "#ff7cdb", // cyan
	37: "#C4C4C4", // white

	90: "#676767", // bright black
	91: "#ff5f87", // bright red
	92: "#00dc7f", // bright green
	93: "#e8e8a8", // bright yellow
	94: "#00aaff", // bright blue
	95: "#635adf", // bright magenta
	96: "#ff7cdb", // bright cyan
	97: "#C4C4C4", // bright white
}
