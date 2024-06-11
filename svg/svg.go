package svg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/beevik/etree"
)

// AddShadow adds a definition of a shadow to the <defs> with the given id.
func AddShadow(element *etree.Element, id string, x, y, blur float64) {
	f := etree.NewElement("filter")
	f.CreateAttr("id", id)
	f.CreateAttr("filterUnits", "userSpaceOnUse")

	b := etree.NewElement("feGaussianBlur")
	b.CreateAttr("in", "SourceAlpha")
	b.CreateAttr("stdDeviation", fmt.Sprintf("%.2f", blur))

	o := etree.NewElement("feOffset")
	o.CreateAttr("result", "offsetblur")
	o.CreateAttr("dx", fmt.Sprintf("%.2f", x))
	o.CreateAttr("dy", fmt.Sprintf("%.2f", y))

	m := etree.NewElement("feMerge")
	mn1 := etree.NewElement("feMergeNode")
	mn2 := etree.NewElement("feMergeNode")
	mn2.CreateAttr("in", "SourceGraphic")
	m.AddChild(mn1)
	m.AddChild(mn2)

	f.AddChild(b)
	f.AddChild(o)
	f.AddChild(m)

	defs := etree.NewElement("defs")
	defs.AddChild(f)
	element.AddChild(defs)
}

// AddClipPath adds a definition of a clip path to the <defs> with the given id.
func AddClipPath(element *etree.Element, id string, x, y, w, h float64) {
	p := etree.NewElement("clipPath")
	p.CreateAttr("id", id)

	rect := etree.NewElement("rect")
	rect.CreateAttr("x", fmt.Sprintf("%.2f", x))
	rect.CreateAttr("y", fmt.Sprintf("%.2f", y))
	rect.CreateAttr("width", fmt.Sprintf("%.2f", w))
	rect.CreateAttr("height", fmt.Sprintf("%.2f", h))

	p.AddChild(rect)

	defs := etree.NewElement("defs")
	defs.AddChild(p)
	element.AddChild(defs)
}

// AddCornerRadius adds corner radius to an element.
func AddCornerRadius(e *etree.Element, radius float64) {
	e.CreateAttr("rx", fmt.Sprintf("%.2f", radius))
	e.CreateAttr("ry", fmt.Sprintf("%.2f", radius))
}

// Move moves the given element to the (x, y) position
func Move(e *etree.Element, x, y float64) {
	e.CreateAttr("x", fmt.Sprintf("%.2fpx", x))
	e.CreateAttr("y", fmt.Sprintf("%.2fpx", y))
}

// AddOutline adds an outline to the given element.
func AddOutline(e *etree.Element, width float64, color string) {
	e.CreateAttr("stroke", color)
	e.CreateAttr("stroke-width", fmt.Sprintf("%.2f", width))
}

const (
	red    string = "#FF5A54"
	yellow string = "#E6BF29"
	green  string = "#52C12B"
)

// NewWindowControls returns a colorful window bar element.
func NewWindowControls(r float64, x, y float64) (*etree.Element, float64) {
	bar := etree.NewElement("svg")
	for i, color := range []string{red, yellow, green} {
		circle := etree.NewElement("circle")
		circle.CreateAttr("cx", fmt.Sprintf("%.2f", float64(i+1)*float64(x)-float64(r)))
		circle.CreateAttr("cy", fmt.Sprintf("%.2f", y))
		circle.CreateAttr("r", fmt.Sprintf("%.2f", r))
		circle.CreateAttr("fill", color)
		bar.AddChild(circle)
	}
	controlsWidth := float64(len(bar.ChildElements()))*float64(x) + float64(r)*2
	return bar, controlsWidth
}

// NewWindowTitle returns a title element with the given text.
func NewWindowTitle(positions []float64, position string, fs float64, ff, text string, s *chroma.Style) (*etree.Element, error) {
	if text == "" || text == "-" {
		return nil, errors.New("Invalid title provided")
	}
	if position != "left" && position != "center" && position != "right" {
		return nil, errors.New("Invalid title position. Must be one of \"left\", \"center\", or \"right\"")
	}
	// positions[0] left, [1] center, [2] right, [3] top
	x := 0.0
	y := positions[3]
	var anchor string
	switch position {
	case "left":
		x = positions[0]
		anchor = "start"
		break
	case "center":
		x = positions[1]
		anchor = "middle"
		break
	case "right":
		x = positions[2]
		anchor = "end"
		break
	}
	input := etree.NewElement("text")
	input.CreateAttr("font-size", fmt.Sprintf("%.2fpx", fs))
	input.CreateAttr("fill", s.Get(chroma.Text).Colour.String())
	input.CreateAttr("font-family", ff)
	input.CreateAttr("text-anchor", anchor)
	input.CreateAttr("alignment-baseline", "middle")
	input.SetText(text)
	Move(input, float64(x), float64(y))
	return input, nil
}

// SetDimensions sets the width and height of the given element.
func SetDimensions(element *etree.Element, width, height float64) {
	widthAttr := element.SelectAttr("width")
	heightAttr := element.SelectAttr("height")
	heightAttr.Value = fmt.Sprintf("%.2f", height)
	widthAttr.Value = fmt.Sprintf("%.2f", width)
}

// GetDimensions returns the width and height of the element.
func GetDimensions(element *etree.Element) (int, int) {
	widthValue := element.SelectAttrValue("width", "0px")
	heightValue := element.SelectAttrValue("height", "0px")
	width := dimensionToInt(widthValue)
	height := dimensionToInt(heightValue)
	return width, height
}

// dimensionToInt takes a string and returns the integer value.
// e.g. "500px" -> 500
func dimensionToInt(px string) int {
	d := strings.TrimSuffix(px, "px")
	v, _ := strconv.ParseInt(d, 10, 64)
	return int(v)
}
