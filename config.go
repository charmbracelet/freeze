package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const defaultOutputFilename = "freeze.svg"

// Config is the configuration options for a screenshot.
type Config struct {
	Input string `json:",omitempty" arg:"" help:"Code to screenshot." optional:""`

	// Window
	Background string `json:"background" help:"Apply a background fill." short:"b" placeholder:"#FFF" group:"Window"`
	Margin     []int  `json:"margin" help:"Apply margin to the window." short:"m" placeholder:"0" group:"Window"`
	Padding    []int  `json:"padding" help:"Apply padding to the code." short:"p" placeholder:"0" group:"Window"`
	Window     bool   `json:"window" help:"Display window controls." group:"Window"`
	Width      int    `json:"width" help:"Window width" short:"W" group:"Window"`
	Height     int    `json:"height" help:"Window height" short:"H" group:"Window"`

	// Settings
	Config      string `json:"config,omitempty" help:"Base configuration file or template." short:"c" group:"Settings" default:"default" placeholder:"base"`
	Interactive bool   `json:",omitempty" help:"Use an interactive form for configuration options." short:"i" group:"Settings"`
	Language    string `json:"language,omitempty" help:"Language of code file." short:"l" group:"Settings" placeholder:"go"`
	Theme       string `json:"theme" help:"Theme to use for syntax highlighting." short:"t" group:"Settings" placeholder:"charm"`

	Output  string `json:"output,omitempty" help:"Output location for {{.svg}}, {{.png}}, or {{.webp}}." short:"o" group:"Settings" default:"" placeholder:"out.svg"`
	Execute string `json:"-" help:"Capture output of command" short:"x" group:"Settings" default:""`

	// Decoration
	Border Border `json:"border" embed:"" prefix:"border." group:"Border"`
	Shadow Shadow `json:"shadow" embed:"" prefix:"shadow." help:"add a shadow to the window" short:"s" group:"Shadow"`

	// Font
	Font Font `json:"font" embed:"" prefix:"font." group:"Font"`

	// Line
	LineHeight      float64 `json:"line_height" help:"Line height relative to font size." group:"Line" placeholder:"1.2"`
	Lines           []int   `json:"-" help:"Lines to capture (start,end)." group:"Line" placeholder:"0,-1" value:"0,-1"`
	ShowLineNumbers bool    `json:"line_numbers" help:"" group:"Line" placeholder:"false"`
}

// Shadow is the configuration options for a drop shadow.
type Shadow struct {
	Blur int `json:"blur" help:"Shadow Gaussian Blur." placeholder:"0"`
	X    int `json:"x" help:"Shadow offset {{x}} coordinate." placeholder:"0"`
	Y    int `json:"y" help:"Shadow offset {{y}} coordinate." placeholder:"0"`
}

// Border is the configuration options for a window border.
type Border struct {
	Radius int    `json:"radius" help:"Corner radius of window." short:"r" placeholder:"0"`
	Width  int    `json:"width" help:"Border width thickness." placeholder:"1"`
	Color  string `json:"color" help:"Border color." placeholder:"#000"`
}

// Font is the configuration options for a font.
type Font struct {
	Family    string  `json:"family" help:"Font family to use for code." placeholder:"monospace"`
	File      string  `json:"file" help:"Font file to embed." placeholder:"monospace.ttf"`
	Size      float64 `json:"size" help:"Font size to use for code." placeholder:"14"`
	Ligatures bool    `json:"ligatures" help:"Use ligatures in the font." placeholder:"true" value:"true" negatable:""`
}

//go:embed configurations/*
var configs embed.FS

func expandPadding(p []int) []int {
	switch len(p) {
	case 1:
		return []int{p[top], p[top], p[top], p[top]}
	case 2:
		return []int{p[top], p[right], p[top], p[right]}
	case 4:
		return []int{p[top], p[right], p[bottom], p[left]}
	default:
		return []int{0, 0, 0, 0}

	}
}

var expandMargin = expandPadding

type side int

const (
	top    side = 0
	right  side = 1
	bottom side = 2
	left   side = 3
)

var userConfigPath = filepath.Join(xdg.ConfigHome, "freeze", "default.json")

func loadUserConfig() (fs.File, error) {
	return os.Open(userConfigPath)
}

func saveUserConfig(config Config) error {
	config.Input = ""
	config.Output = ""
	config.Interactive = false

	err := os.MkdirAll(filepath.Dir(userConfigPath), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(userConfigPath)
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}
