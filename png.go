package main

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

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

func resvgConvert(doc *etree.Document, w, h float64, output string, config *Config) error {
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

	// resvg's wasm sandbox cannot open arbitrary host paths, so we must read
	// font files on the host and pass the bytes through.
	if config != nil && config.Font.File != "" {
		bts, err := os.ReadFile(config.Font.File)
		if err != nil {
			printErrorFatal("Unable to load font", err)
		}
		if err := fontdb.LoadFontData(bts); err != nil {
			printErrorFatal("Unable to load font", err)
		}
	}

	// Load system fonts so a user-specified --font.family can resolve.
	// We walk on the host and feed bytes via LoadFontData, since the wasm
	// sandbox blocks direct filesystem access. Errors are intentionally
	// ignored: unreadable files and missing directories are skipped.
	for _, dir := range systemFontDirs() {
		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil //nolint:nilerr
			}
			ext := filepath.Ext(path)
			if ext != ".ttf" && ext != ".otf" && ext != ".ttc" {
				return nil
			}
			bts, err := os.ReadFile(path) //nolint:gosec
			if err != nil {
				return nil //nolint:nilerr
			}
			_ = fontdb.LoadFontData(bts)
			return nil
		})
	}

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

	fontFamily := "JetBrains Mono"
	if config != nil && config.Font.Family != "" {
		fontFamily = config.Font.Family
	}

	tree, err := worker.NewTreeFromData(svg, &resvg.Options{
		Dpi:                192,
		FontFamily:         fontFamily,
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
