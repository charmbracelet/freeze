package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"

	"github.com/beevik/etree"
	"github.com/charmbracelet/freeze/font"
	"github.com/kanrichan/resvg-go"
)

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

	err = os.WriteFile(output, png, 0o600)
	if err != nil {
		return err //nolint: wrapcheck
	}
	return err //nolint: wrapcheck
}
