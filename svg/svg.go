package svg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

// AddShadow adds a definition of a shadow to the <defs> with the given id.
func AddShadow(element *etree.Element, id string, x, y, blur int) {
	f := etree.NewElement("filter")
	f.CreateAttr("id", id)
	f.CreateAttr("filterUnits", "userSpaceOnUse")

	b := etree.NewElement("feGaussianBlur")
	b.CreateAttr("in", "SourceAlpha")
	b.CreateAttr("stdDeviation", fmt.Sprintf("%d", blur))

	o := etree.NewElement("feOffset")
	o.CreateAttr("result", "offsetblur")
	o.CreateAttr("dx", fmt.Sprintf("%d", x))
	o.CreateAttr("dy", fmt.Sprintf("%d", y))

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

// AddCornerRadius adds corner radius to an element.
func AddCornerRadius(e *etree.Element, radius int) {
	e.CreateAttr("rx", fmt.Sprintf("%d", radius))
	e.CreateAttr("ry", fmt.Sprintf("%d", radius))
}

// Move moves the given element to the (x, y) position
func Move(e *etree.Element, x, y float64) {
	e.CreateAttr("x", fmt.Sprintf("%.2fpx", x))
	e.CreateAttr("y", fmt.Sprintf("%.2fpx", y))
}

// AddOutline adds an outline to the given element.
func AddOutline(e *etree.Element, width int, color string) {
	e.CreateAttr("stroke", color)
	e.CreateAttr("stroke-width", fmt.Sprintf("%d", width))
}

const (
	red    string = "#FF5A54"
	yellow string = "#E6BF29"
	green  string = "#52C12B"
)

// NewWindowControls returns a colorful window bar element.
func NewWindowControls() *etree.Element {
	bar := etree.NewElement("svg")
	for i, color := range []string{red, yellow, green} {
		circle := etree.NewElement("circle")
		circle.CreateAttr("cx", fmt.Sprintf("%d", (i+1)*19+-6))
		circle.CreateAttr("cy", fmt.Sprintf("%d", 12))
		circle.CreateAttr("r", "5.5")
		circle.CreateAttr("fill", color)
		bar.AddChild(circle)
	}
	return bar
}

// SetDimensions sets the width and height of the given element.
func SetDimensions(element *etree.Element, width, height int) {
	widthAttr := element.SelectAttr("width")
	heightAttr := element.SelectAttr("height")
	heightAttr.Value = fmt.Sprintf("%d", height)
	widthAttr.Value = fmt.Sprintf("%d", width)
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
