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

func libsvgConvert(doc *etree.Document, w, h float64, output string) error {
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
	rsvgConvert := exec.Command("rsvg-convert", "-o", output)
	rsvgConvert.Stdin = bytes.NewReader(svg)
	err = rsvgConvert.Run()
	return err
}

func resvgConvert(doc *etree.Document, w, h float64, output, fontFile string) error {
	svg, err := doc.WriteToBytes()
	if err != nil {
		return err
	}

	worker, err := resvg.NewDefaultWorker(context.Background())
	if err != nil {
		printErrorFatal("Unable to write output", err)
	}
	defer worker.Close()

	fontdb, err := worker.NewFontDBDefault()
	if err != nil {
		printErrorFatal("Unable to write output", err)
	}
	defer fontdb.Close()

	if fontFile == "" {
		err = fontdb.LoadFontData(font.JetBrainsMonoTTF)
		if err != nil {
			printErrorFatal("Unable to load font", err)
		}
		err = fontdb.LoadFontData(font.JetBrainsMonoNLTTF)
		if err != nil {
			printErrorFatal("Unable to load font", err)
		}

		for _, d := range font.DefaultFontsDirectories {
			err = fontdb.LoadFontsDir(d)
			if err != nil {
				if os.IsNotExist(err) || os.IsPermission(err) {
					continue
				}
				printErrorFatal("Unable to load font dir", err)
			}
		}
	} else {
		err = fontdb.LoadFontFile(fontFile)
		if err != nil {
			printErrorFatal("Unable to load font", err)
		}
	}

	pixmap, err := worker.NewPixmap(uint32(w), uint32(h))
	if err != nil {
		printError("Unable to write output", err)
		os.Exit(1)
	}
	defer pixmap.Close()

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
	defer tree.Close()

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
