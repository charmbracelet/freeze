package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strconv"

	"github.com/beevik/etree"
	"github.com/charmbracelet/freeze/font"
	"github.com/kanrichan/resvg-go"
)

func libsvgConvert(doc *etree.Document, w, _ float64, output string) error {
	_, err := exec.LookPath("rsvg-convert")
	if err != nil {
		return err
	}

	svg, err := doc.WriteToBytes()
	if err != nil {
		return err
	}

	// rsvg-convert is installed use that to convert the SVG to PNG,
	// since it is faster.
	rsvgConvert := exec.Command("rsvg-convert",
		"--width", strconv.Itoa(int(w)),
		"--keep-aspect-ratio",
		"-f", "png",
		"-o", output,
	)
	rsvgConvert.Stdin = bytes.NewReader(svg)
	err = rsvgConvert.Run()
	return err
}

func resvgConvert(doc *etree.Document, w, h float64, output string) error {
	svg, err := doc.WriteToBytes()
	if err != nil {
		return err
	}

	worker, err := resvg.NewDefaultWorker(context.Background())
	defer worker.Close()
	if err != nil {
		printErrorFatal("Unable to write output", err)
	}

	fontdb, err := worker.NewFontDBDefault()
	defer fontdb.Close()
	if err != nil {
		printErrorFatal("Unable to write output", err)
	}
	fontdb.LoadFontData(font.JetBrainsMonoTTF)
	fontdb.LoadFontData(font.JetBrainsMonoNLTTF)

	pixmap, err := worker.NewPixmap(uint32(w), uint32(h))
	defer pixmap.Close()
	if err != nil {
		printError("Unable to write output", err)
		os.Exit(1)
	}

	tree, err := worker.NewTreeFromData(svg, &resvg.Options{
		Dpi:                96,
		ShapeRenderingMode: resvg.ShapeRenderingModeGeometricPrecision,
		TextRenderingMode:  resvg.TextRenderingModeOptimizeLegibility,
		ImageRenderingMode: resvg.ImageRenderingModeOptimizeQuality,
		DefaultSizeWidth:   float32(w),
		DefaultSizeHeight:  float32(h),
	})
	defer tree.Close()
	if err != nil {
		printError("Unable to write output", err)
		os.Exit(1)
	}

	err = tree.ConvertText(fontdb)
	if err != nil {
		return err
	}
	err = tree.Render(resvg.TransformIdentity(), pixmap)
	if err != nil {
		return err
	}
	png, err := pixmap.EncodePNG()
	if err != nil {
		return err
	}

	err = os.WriteFile(output, png, 0644)
	if err != nil {
		return err
	}
	return err
}
