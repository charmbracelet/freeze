package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/beevik/etree"
	"github.com/charmbracelet/freeze/font"
	"github.com/kanrichan/resvg-go"
	"golang.design/x/clipboard"
)

func copyToClipboard(path string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}
	// check if WAYLAND_DISPLAY is set
	switch os.Getenv("XDG_SESSION_TYPE") {
	case "wayland":
		var err error
		if _, pathErr := exec.LookPath("wl-copy"); err != nil {
			printError("Unable to find wl-copy in your path. Please install it to use clipboard features in wayland.", pathErr)
		} else {
			cmd := exec.Command("wl-copy", "--type", "image/png", "<", path)
			err = cmd.Run()
		}
		return err
		// TODO add tests for GH actions
	case "x11":
		// this is x11
		var err error
		if _, pathErr := exec.LookPath("xclip"); err != nil {
			printError("Unable to find xclip in your path. Please install it to use clipboard features in x11.", pathErr)
		} else {
			cmd := exec.Command("xclip", "-selection", "clipboard", "-t", "image/png", "-i", path)
			err = cmd.Run()
		}
		return err
	}
	png, err := os.ReadFile(path)
	defer os.Remove(path) // nolint: errcheck
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtImage, png)
	clipboard.Read(clipboard.FmtImage)
	return err
}

func libsvgConvert(doc *etree.Document, _, _ float64, output string) error {
	_, err := exec.LookPath("rsvg-convert")
	if err != nil {
		return err //nolint: wrapcheck
	}

	svg, err := doc.WriteToBytes()
	if err != nil {
		return err //nolint: wrapcheck
	}

	// rsvg-convert is installed use that to convert the SVG to PNG,
	// since it is faster.
	rsvgConvert := exec.Command("rsvg-convert", "-o", output)
	rsvgConvert.Stdin = bytes.NewReader(svg)
	err = rsvgConvert.Run()
	if err != nil {
		return err
	}
	if strings.HasPrefix(output, "clipboard") {
		return copyToClipboard(output)
	}
	return err //nolint: wrapcheck
}

func resvgConvert(doc *etree.Document, w, h float64, output string) error {
	svg, err := doc.WriteToBytes()
	if err != nil {
		return err //nolint: wrapcheck
	}

	worker, err := resvg.NewDefaultWorker(context.Background())
	if err != nil {
		printErrorFatal("Unable to write output", err)
	}
	defer worker.Close() //nolint: errcheck

	fontdb, err := worker.NewFontDBDefault()
	if err != nil {
		printErrorFatal("Unable to write output", err)
	}
	defer fontdb.Close() //nolint: errcheck
	err = fontdb.LoadFontData(font.JetBrainsMonoTTF)
	if err != nil {
		printErrorFatal("Unable to load font", err)
	}
	err = fontdb.LoadFontData(font.JetBrainsMonoNLTTF)
	if err != nil {
		printErrorFatal("Unable to load font", err)
	}

	pixmap, err := worker.NewPixmap(uint32(w), uint32(h))
	if err != nil {
		printError("Unable to write output", err)
		os.Exit(1)
	}
	defer pixmap.Close() //nolint: errcheck

	tree, err := worker.NewTreeFromData(svg, &resvg.Options{
		Dpi:                192,
		ShapeRenderingMode: resvg.ShapeRenderingModeGeometricPrecision,
		TextRenderingMode:  resvg.TextRenderingModeOptimizeLegibility,
		ImageRenderingMode: resvg.ImageRenderingModeOptimizeQuality,
		DefaultSizeWidth:   float32(w),
		DefaultSizeHeight:  float32(h),
	})
	if err != nil {
		printError("Unable to write output", err)
		os.Exit(1)
	}
	defer tree.Close() //nolint: errcheck

	err = tree.ConvertText(fontdb)
	if err != nil {
		return err //nolint: wrapcheck
	}
	err = tree.Render(resvg.TransformIdentity(), pixmap)
	if err != nil {
		return err //nolint: wrapcheck
	}
	png, err := pixmap.EncodePNG()
	if err != nil {
		return err //nolint: wrapcheck
	}

	if output == "clipboard.png" {
		return copyToClipboard(output)
	}
	return os.WriteFile(output, png, 0o600)
}
