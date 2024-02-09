package main

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/mattn/go-runewidth"
)

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
		p.row++
		p.col = 0
		p.endBackground()
	}
}
func (p *dispatcher) OscDispatch(params [][]byte, bellTerminated bool)      {}
func (p *dispatcher) EscDispatch(intermediates []byte, r rune, ignore bool) {}
func (p *dispatcher) DcsHook(prefix string, params [][]uint16, intermediates []byte, r rune, ignore bool) {
}
func (p *dispatcher) DcsPut(code byte) {}
func (p *dispatcher) DcsUnhook()       {}

const fontHeightToWidthRatio = 1.67

func (p *dispatcher) beginBackground(fill string) {
	rect := etree.NewElement("rect")
	rect.CreateAttr("fill", fill)
	rect.CreateAttr("x", fmt.Sprintf("%.2fpx", (float64(p.col)*p.config.Font.Size/fontHeightToWidthRatio)+float64(p.config.Margin[left]+p.config.Padding[left])))
	rect.CreateAttr("y", fmt.Sprintf("%.2fpx", float64(p.row)*p.config.Font.Size*p.config.LineHeight+float64(p.config.Margin[top]+p.config.Padding[top])))
	rect.CreateAttr("height", fmt.Sprintf("%.2fpx", p.config.Font.Size*p.config.LineHeight+1))
	p.bg = rect
}

func (p *dispatcher) endBackground() {
	if p.bg == nil {
		return
	}

	p.bg.CreateAttr("width", fmt.Sprintf("%.2fpx", float64(p.bgWidth)*p.config.Font.Size/fontHeightToWidthRatio+1))
	p.svg.InsertChildAt(0, p.bg)
	p.bg = nil
	p.bgWidth = 0
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
