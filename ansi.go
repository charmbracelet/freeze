package main

import (
	"fmt"
	"regexp"

	"github.com/beevik/etree"
	"github.com/mattn/go-runewidth"
)

const csiRegex = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
const oscRegex = "[\u001B\u009B\u009D]"

var ansiRe = regexp.MustCompile(csiRegex)
var oscRe = regexp.MustCompile(oscRegex)

func StripANSI(str string) string {
	str = ansiRe.ReplaceAllString(str, "")
	str = oscRe.ReplaceAllString(str, "")
	return str
}

type dispatcher struct {
	lines   []*etree.Element
	row     int
	col     int
	svg     *etree.Element
	bg      *etree.Element
	bgWidth int
	config  *Config
}

func (p *dispatcher) Print(r rune) {
	p.row = clamp(p.row, 0, len(p.lines)-1)
	// insert the rune in the last tspan
	children := p.lines[p.row].ChildElements()
	var lastChild *etree.Element
	isFirstChild := len(children) == 0
	if isFirstChild {
		lastChild = etree.NewElement("tspan")
		lastChild.CreateAttr("xml:space", "preserve")
		p.lines[p.row].AddChild(lastChild)
	} else {
		lastChild = children[len(children)-1]
	}

	if runewidth.RuneWidth(r) > 1 {
		newChild := lastChild.Copy()
		newChild.SetText(string(r))
		newChild.CreateAttr("dx", fmt.Sprintf("%.2fpx", p.config.Font.Size/5))
		p.lines[p.row].AddChild(newChild)
	} else {
		lastChild.SetText(lastChild.Text() + string(r))
	}

	p.col += runewidth.RuneWidth(r)
	if p.bg != nil {
		p.bgWidth += runewidth.RuneWidth(r)
	}
}

func (p *dispatcher) Execute(code byte) {
	if code == 0x0A {
		p.endBackground()
		p.row++
		p.col = 0
	}
}
func (p *dispatcher) OscDispatch(params [][]byte, bellTerminated bool)      {}
func (p *dispatcher) EscDispatch(intermediates []byte, r rune, ignore bool) {}
func (p *dispatcher) DcsHook(prefix string, params [][]uint16, intermediates []byte, r rune, ignore bool) {
}
func (p *dispatcher) DcsPut(code byte) {}
func (p *dispatcher) DcsUnhook()       {}

const fontHeightToWidthRatio = 1.68

func (p *dispatcher) beginBackground(fill string) {
	if p.col == 0 {
		return
	}

	rect := etree.NewElement("rect")
	rect.CreateAttr("fill", fill)
	y := fmt.Sprintf("%.2fpx", float64(p.row)*p.config.Font.Size*p.config.LineHeight+float64(p.config.Margin[top]+p.config.Padding[top]))
	x := float64(p.col) * (p.config.Font.Size / fontHeightToWidthRatio)
	x += float64(p.config.Margin[left] + p.config.Padding[left])
	if p.config.ShowLineNumbers {
		x += float64(p.config.Font.Size) * 3
	}
	rect.CreateAttr("x", fmt.Sprintf("%.2fpx", x))
	rect.CreateAttr("y", y)
	rect.CreateAttr("height", fmt.Sprintf("%.2fpx", p.config.Font.Size*p.config.LineHeight+1))
	p.bg = rect
}

func (p *dispatcher) endBackground() {
	if p.bg == nil {
		return
	}

	p.bg.CreateAttr("width", fmt.Sprintf("%.5fpx", (float64(p.bgWidth)+0.5)*(p.config.Font.Size/fontHeightToWidthRatio)))
	p.svg.InsertChildAt(0, p.bg)
	p.bg = nil
	p.bgWidth = 0
}

func (p *dispatcher) CsiDispatch(prefix string, params [][]uint16, intermediates []byte, r rune, ignore bool) {
	if ignore {
		return
	}

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
			// span.CreateAttr("font-weight", "bold")
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
				fill := fmt.Sprintf("#%02x%02x%02x", params[i+1][0], params[i+1][0], params[i+1][0])
				n := params[i+1][0]
				i++
				fill = palette[n]
				span.CreateAttr("fill", fill)
				p.lines[p.row].AddChild(span)
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
				n := params[i+1][0]
				i++
				fill = palette[n]
				p.beginBackground(fill)
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
	30: "#282a2e", // black
	31: "#D74E6F", // red
	32: "#31BB71", // green
	33: "#D3E561", // yellow
	34: "#8056FF", // blue
	35: "#ED61D7", // magenta
	36: "#04D7D7", // cyan
	37: "#C5C8C6", // white

	90: "#4B4B4B", // bright black
	91: "#FE5F86", // bright red
	92: "#00D787", // bright green
	93: "#EBFF71", // bright yellow
	94: "#8F69FF", // bright blue
	95: "#FF7AEA", // bright magenta
	96: "#00FEFE", // bright cyan
	97: "#FFFFFF", // bright white
}

var palette = []string{
	"#000000", "#800000", "#008000", "#808000", "#000080", "#800080", "#008080", "#c0c0c0",
	"#808080", "#ff0000", "#00ff00", "#ffff00", "#0000ff", "#ff00ff", "#00ffff", "#ffffff",
	"#000000", "#00005f", "#000087", "#0000af", "#0000d7", "#0000ff", "#005f00", "#005f5f",
	"#005f87", "#005faf", "#005fd7", "#005fff", "#008700", "#00875f", "#008787", "#0087af",
	"#0087d7", "#0087ff", "#00af00", "#00af5f", "#00af87", "#00afaf", "#00afd7", "#00afff",
	"#00d700", "#00d75f", "#00d787", "#00d7af", "#00d7d7", "#00d7ff", "#00ff00", "#00ff5f",
	"#00ff87", "#00ffaf", "#00ffd7", "#00ffff", "#5f0000", "#5f005f", "#5f0087", "#5f00af",
	"#5f00d7", "#5f00ff", "#5f5f00", "#5f5f5f", "#5f5f87", "#5f5faf", "#5f5fd7", "#5f5fff",
	"#5f8700", "#5f875f", "#5f8787", "#5f87af", "#5f87d7", "#5f87ff", "#5faf00", "#5faf5f",
	"#5faf87", "#5fafaf", "#5fafd7", "#5fafff", "#5fd700", "#5fd75f", "#5fd787", "#5fd7af",
	"#5fd7d7", "#5fd7ff", "#5fff00", "#5fff5f", "#5fff87", "#5fffaf", "#5fffd7", "#5fffff",
	"#870000", "#87005f", "#870087", "#8700af", "#8700d7", "#8700ff", "#875f00", "#875f5f",
	"#875f87", "#875faf", "#875fd7", "#875fff", "#878700", "#87875f", "#878787", "#8787af",
	"#8787d7", "#8787ff", "#87af00", "#87af5f", "#87af87", "#87afaf", "#87afd7", "#87afff",
	"#87d700", "#87d75f", "#87d787", "#87d7af", "#87d7d7", "#87d7ff", "#87ff00", "#87ff5f",
	"#87ff87", "#87ffaf", "#87ffd7", "#87ffff", "#af0000", "#af005f", "#af0087", "#af00af",
	"#af00d7", "#af00ff", "#af5f00", "#af5f5f", "#af5f87", "#af5faf", "#af5fd7", "#af5fff",
	"#af8700", "#af875f", "#af8787", "#af87af", "#af87d7", "#af87ff", "#afaf00", "#afaf5f",
	"#afaf87", "#afafaf", "#afafd7", "#afafff", "#afd700", "#afd75f", "#afd787", "#afd7af",
	"#afd7d7", "#afd7ff", "#afff00", "#afff5f", "#afff87", "#afffaf", "#afffd7", "#afffff",
	"#d70000", "#d7005f", "#d70087", "#d700af", "#d700d7", "#d700ff", "#d75f00", "#d75f5f",
	"#d75f87", "#d75faf", "#d75fd7", "#d75fff", "#d78700", "#d7875f", "#d78787", "#d787af",
	"#d787d7", "#d787ff", "#d7af00", "#d7af5f", "#d7af87", "#d7afaf", "#d7afd7", "#d7afff",
	"#d7d700", "#d7d75f", "#d7d787", "#d7d7af", "#d7d7d7", "#d7d7ff", "#d7ff00", "#d7ff5f",
	"#d7ff87", "#d7ffaf", "#d7ffd7", "#d7ffff", "#ff0000", "#ff005f", "#ff0087", "#ff00af",
	"#ff00d7", "#ff00ff", "#ff5f00", "#ff5f5f", "#ff5f87", "#ff5faf", "#ff5fd7", "#ff5fff",
	"#ff8700", "#ff875f", "#ff8787", "#ff87af", "#ff87d7", "#ff87ff", "#ffaf00", "#ffaf5f",
	"#ffaf87", "#ffafaf", "#ffafd7", "#ffafff", "#ffd700", "#ffd75f", "#ffd787", "#ffd7af",
	"#ffd7d7", "#ffd7ff", "#ffff00", "#ffff5f", "#ffff87", "#ffffaf", "#ffffd7", "#ffffff",
	"#080808", "#121212", "#1c1c1c", "#262626", "#303030", "#3a3a3a", "#444444", "#4e4e4e",
	"#585858", "#606060", "#666666", "#767676", "#808080", "#8a8a8a", "#949494", "#9e9e9e",
	"#a8a8a8", "#b2b2b2", "#bcbcbc", "#c6c6c6", "#d0d0d0", "#dadada", "#e4e4e4", "#eeeeee",
}
