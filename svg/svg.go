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
	f.CreateAttr("x", "0")
	f.CreateAttr("y", "0")
	f.CreateAttr("filterUnits", "userSpaceOnUse")

	o := etree.NewElement("feOffset")
	o.CreateAttr("result", "offOut")
	o.CreateAttr("in", "SourceAlpha")
	o.CreateAttr("dx", fmt.Sprintf("%d", x))
	o.CreateAttr("dy", fmt.Sprintf("%d", y))

	c := etree.NewElement("feColorMatrix")
	c.CreateAttr("result", "matrixOut")
	c.CreateAttr("in", "offOut")
	c.CreateAttr("type", "matrix")
	c.CreateAttr("values", "0.2 0 0 0 0 0 0.2 0 0 0 0 0 0.2 0 0 0 0 0 1 0")

	b := etree.NewElement("feGaussianBlur")
	b.CreateAttr("result", "blurOut")
	b.CreateAttr("in", "matrixOut")
	b.CreateAttr("stdDeviation", fmt.Sprintf("%d", blur))

	blend := etree.NewElement("feBlend")
	blend.CreateAttr("in", "SourceGraphic")
	blend.CreateAttr("in2", "blurOut")
	blend.CreateAttr("mode", "normal")

	f.AddChild(o)
	f.AddChild(b)
	f.AddChild(blend)

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
