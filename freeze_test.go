package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/aymanbagabas/go-udiff"
)

const binary = "./test/freeze-test"

var (
	update = flag.Bool("update", false, "update golden files")
	png    = flag.Bool("png", false, "update pngs")
)

func TestMain(m *testing.M) {
	flag.Parse()
	cmd := exec.Command("go", "build", "-o", binary)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	exit := m.Run()
	err = os.Remove(binary)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(exit)
}

func TestFreeze(t *testing.T) {
	cmd := exec.Command(binary)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFreezeOutput(t *testing.T) {
	output := "artichoke-test.svg"
	defer os.Remove(output)

	cmd := exec.Command(binary, "test/input/artichoke.hs", "-o", output)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(output)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFreezeHelp(t *testing.T) {
	out := bytes.Buffer{}
	cmd := exec.Command(binary)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Fatal("unexpected error")
	}

	got := out.String()

	contains := []string{
		"Generate images of code and terminal output.",
		"freeze main.go [-o code.svg] [--flags]",
		"--theme", "Theme to use for syntax highlighting",
		"--border.color", "Border color.",
		"--shadow.blur", "Shadow Gaussian Blur.",
		"--font.family", "Font family to use for code.",
	}

	for _, c := range contains {
		if !strings.Contains(got, c) {
			t.Fatalf("expected %s to contain \"%s\"", got, c)
		}
	}
}

func TestFreezeErrorFileMissing(t *testing.T) {
	out := bytes.Buffer{}
	cmd := exec.Command(binary, "this-file-does-not-exist")
	cmd.Stdout = &out
	err := cmd.Run()

	if err == nil {
		t.Fatal("expected error")
	}

	got := out.String()

	contains := []string{"ERROR", "File not found", "open this-file-does-not-exist: no such file or directory"}

	for _, c := range contains {
		if !strings.Contains(got, c) {
			t.Fatalf("expected %s to contain \"%s\"", got, c)
		}
	}
}

func TestFreezeConfigurations(t *testing.T) {
	tests := []struct {
		input  string
		flags  []string
		output string
	}{
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--config", "test/configurations/base.json"},
			output: "artichoke-base",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--config", "test/configurations/full.json"},
			output: "artichoke-full",
		},
		{
			input:  "test/input/eza.ansi",
			flags:  []string{"--config", "full"},
			output: "eza",
		},
		{
			flags:  []string{"--execute", `echo "Hello, world!"`},
			output: "execute",
		},
		{
			input:  "test/input/bubbletea.model",
			flags:  []string{"--language", "go", "--height", "800", "--width", "750", "--config", "full", "--window=false", "--show-line-numbers"},
			output: "bubbletea",
		},
		// {
		// 	flags:  []string{"--execute", "layout", "--height", "800", "--config", "full", "--margin", "50,10"},
		// 	output: "composite-2",
		// },
		{
			input:  "test/input/layout.ansi",
			flags:  []string{},
			output: "layout",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--language", "haskell"},
			output: "haskell",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--theme", "dracula"},
			output: "dracula",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--border.radius", "8"},
			output: "border-radius",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--border.radius", "8", "--window"},
			output: "window",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--border.radius", "8", "--window", "--border.width", "1"},
			output: "border-width",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--border.radius", "8", "--window", "--border.width", "1", "--padding", "30,50,30,30"},
			output: "padding",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--border.radius", "8", "--window", "--border.width", "1", "--padding", "30,50,30,30", "--margin", "50,60,100,60"},
			output: "margin",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--config", "full"},
			output: "shadow",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--width", "1920", "--height", "1080"},
			output: "dimensions",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--margin", "50", "--width", "600", "--height", "300"},
			output: "dimensions-margin",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--margin", "50", "--width", "600", "--height", "300", "--show-line-numbers"},
			output: "dimensions-margin-line-numbers",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--padding", "50", "--width", "600", "--height", "300"},
			output: "dimensions-padding",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--config", "full", "--width", "600", "--height", "300"},
			output: "dimensions-config",
		},
		{
			input:  "test/input/goreleaser-full.yml",
			flags:  []string{"--config", "full", "--width", "600", "--height", "900"},
			output: "overflow",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--config", "full", "--lines", "4,8", "--show-line-numbers"},
			output: "lines",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--font.size", "28"},
			output: "font-size-28",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--font.size", "14"},
			output: "font-size-14",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--line-height", "2"},
			output: "line-height-2",
		},
		{
			input:  "test/input/goreleaser-full.yml",
			flags:  []string{"--config", "full", "--height", "2000", "--show-line-numbers"},
			output: "overflow-line-numbers",
		},
		{
			input:  "test/input/helix.ansi",
			flags:  []string{"--background", "#0d1116"},
			output: "helix",
		},
		{
			input:  "test/input/glow.ansi",
			flags:  []string{},
			output: "glow",
		},
		{
			input:  "test/input/tab.go",
			flags:  []string{},
			output: "tab",
		},
		{
			input:  "test/input/wrap.go",
			flags:  []string{"--wrap", "80", "--width", "600"},
			output: "wrap",
		},
	}

	err := os.RemoveAll("test/output/svg")
	if err != nil {
		t.Fatal("unable to remove output files")
	}
	err = os.MkdirAll("test/output/svg", 0o755)
	if err != nil {
		t.Fatal("unable to create output directory")
	}
	err = os.MkdirAll("test/golden/svg", 0o755)
	if err != nil {
		t.Fatal("unable to create output directory")
	}
	err = os.MkdirAll("test/output/png", 0o755)
	if err != nil {
		t.Fatal("unable to create output directory")
	}

	for _, tc := range tests {
		t.Run(tc.output, func(t *testing.T) {
			// output SVG
			out := bytes.Buffer{}
			args := []string{tc.input}
			args = append(args, tc.flags...)
			args = append(args, "--output", "test/output/svg/"+tc.output+".svg")
			cmd := exec.Command(binary, args...)
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				t.Log(err)
				t.Log(out.String())
				t.Fatal("unexpected error")
			}
			gotfile := "test/output/svg/" + tc.output + ".svg"
			got, err := os.ReadFile(gotfile)
			if err != nil {
				t.Fatal("no output file for:", gotfile)
			}
			goldenfile := "test/golden/svg/" + tc.output + ".svg"
			if *update {
				if err := os.WriteFile(goldenfile, got, 0o644); err != nil {
					t.Log(err)
					t.Fatal("unexpected error")
				}
			}
			want, err := os.ReadFile(goldenfile)
			if err != nil {
				t.Fatal("no golden file for:", goldenfile)
			}
			if string(want) != string(got) {
				t.Log(udiff.Unified("want", "got", string(want), string(got)))
				t.Fatalf("%s != %s", goldenfile, gotfile)
			}

			// output PNG
			if png != nil && *png {
				out = bytes.Buffer{}
				args = []string{tc.input}
				args = append(args, tc.flags...)
				args = append(args, "--output", "test/output/png/"+tc.output+".png")
				cmd = exec.Command(binary, args...)
				cmd.Stdout = &out
				err = cmd.Run()
				if err != nil {
					t.Log(err)
					t.Log(out.String())
					t.Fatal("unexpected error")
				}
			}
		})
	}
}
