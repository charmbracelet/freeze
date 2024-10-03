package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
)

const defaultOutputFilename = "freeze.png"

// Config is the configuration options for a screenshot.
type Config struct {
	Input string `json:",omitempty" arg:"" help:"Code to screenshot." optional:""`

	// Window
	Background string    `json:"background" help:"Apply a background fill." short:"b" placeholder:"#171717" group:"Window"`
	Margin     []float64 `json:"margin" help:"Apply margin to the window." short:"m" placeholder:"0" group:"Window"`
	Padding    []float64 `json:"padding" help:"Apply padding to the code." short:"p" placeholder:"0" group:"Window"`
	Window     bool      `json:"window" help:"Display window controls." group:"Window"`
	Width      float64   `json:"width" help:"Width of terminal window." short:"W" group:"Window"`
	Height     float64   `json:"height" help:"Height of terminal window." short:"H" group:"Window"`

	// Settings
	Config      string `json:"config,omitempty" help:"Base configuration file or template." short:"c" group:"Settings" default:"default" placeholder:"base"`
	Interactive bool   `hidden:"" json:",omitempty" help:"Use an interactive form for configuration options." short:"i" group:"Settings"`
	Language    string `json:"language,omitempty" help:"Language of code file." short:"l" group:"Settings" placeholder:"go"`
	Theme       string `json:"theme" help:"Theme to use for syntax highlighting." short:"t" group:"Settings" placeholder:"charm"`
	Wrap        int    `json:"wrap" help:"Wrap lines at a specific width." short:"w" group:"Settings" default:"0" placeholder:"80"`

	Output         string        `json:"output,omitempty" help:"Output location for {{.svg}}, {{.png}}, or {{.webp}}." short:"o" group:"Settings" default:"" placeholder:"freeze.svg"`
	Execute        string        `json:"-" help:"Capture output of command execution." short:"x" group:"Settings" default:""`
	ExecuteTimeout time.Duration `json:"-" help:"Execution timeout." group:"Settings" default:"10s" prefix:"execute." name:"timeout" hidden:""`

	// Decoration
	Border Border `json:"border" embed:"" prefix:"border." group:"Border"`
	Shadow Shadow `json:"shadow" embed:"" prefix:"shadow." help:"add a shadow to the window" short:"s" group:"Shadow"`

	// Font
	Font Font `json:"font" embed:"" prefix:"font." group:"Font"`

	// Line
	LineHeight      float64 `json:"line_height" help:"Line height relative to font size." group:"Line" placeholder:"1.2"`
	Lines           []int   `json:"-" help:"Lines to capture (start,end)." group:"Line" placeholder:"0,-1" value:"0,-1"`
	ShowLineNumbers bool    `json:"show_line_numbers" help:"" group:"Line" placeholder:"false"`
}

// Shadow is the configuration options for a drop shadow.
type Shadow struct {
	Blur float64 `json:"blur" help:"Shadow Gaussian Blur." placeholder:"0"`
	X    float64 `json:"x" help:"Shadow offset {{x}} coordinate." placeholder:"0"`
	Y    float64 `json:"y" help:"Shadow offset {{y}} coordinate." placeholder:"0"`
}

// Border is the configuration options for a window border.
type Border struct {
	Radius float64 `json:"radius" help:"Corner radius of window." short:"r" placeholder:"0"`
	Width  float64 `json:"width" help:"Border width thickness." placeholder:"1"`
	Color  string  `json:"color" help:"Border color." placeholder:"#000"`
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

func expandPadding(p []float64, scale float64) []float64 {
	switch len(p) {
	case 1:
		return []float64{p[top] * scale, p[top] * scale, p[top] * scale, p[top] * scale}
	case 2:
		return []float64{p[top] * scale, p[right] * scale, p[top] * scale, p[right] * scale}
	case 4:
		return []float64{p[top] * scale, p[right] * scale, p[bottom] * scale, p[left] * scale}
	default:
		return []float64{0, 0, 0, 0}
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

var userConfigPath = filepath.Join(xdg.ConfigHome, "freeze", "user.json")

func loadUserConfig() (fs.File, error) {
	return os.Open(userConfigPath) //nolint: wrapcheck
}

func saveUserConfig(config Config) error {
	config.Input = ""
	config.Output = ""
	config.Interactive = false

	err := os.MkdirAll(filepath.Dir(userConfigPath), os.ModePerm)
	if err != nil {
		return err //nolint: wrapcheck
	}
	f, err := os.Create(userConfigPath)
	if err != nil {
		return err //nolint: wrapcheck
	}
	b, err := json.Marshal(config)
	if err != nil {
		return err //nolint: wrapcheck
	}
	_, err = f.Write(b)

	printFilenameOutput(userConfigPath)

	return err //nolint: wrapcheck
}
