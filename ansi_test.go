package main

import (
	"strings"
	"testing"

	"github.com/beevik/etree"
)

func TestAnsiBackgroundAndDimmed(t *testing.T) {
	cfg := DefaultConfig()
	doc := etree.NewDocument()
	svg := doc.CreateElement("svg")
	lines := []*etree.Element{
		svg.CreateElement("text"),
	}

	p := &dispatcher{
		scale:  1,
		svg:    svg,
		config: &cfg,
		lines:  lines,
		row:    0,
		col:    0,
	}

	// Test dimmed modifier (\x1b[2m) and standard bg (\x1b[41m)
	// \x1b[41m sets background red, \x1b[2m sets opacity 0.5
	p.CsiDispatch('m', []int{41})
	p.Print('A')
	p.CsiDispatch('m', []int{2})
	p.Print('B')
	p.endBackground()

	rendered := doc.WriteToString()
	if !strings.Contains(rendered, "fill=\"#D74E6F\"") {
		t.Fatalf("expected background color fill #D74E6F, got: %s", rendered)
	}
	if !strings.Contains(rendered, "opacity=\"0.5\"") {
		t.Fatalf("expected opacity=\"0.5\" for dimmed modifier, got: %s", rendered)
	}
}
