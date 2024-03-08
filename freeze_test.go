package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const binary = "./freeze-test"

func TestMain(m *testing.M) {
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
		"WINDOW",
		"--background", "Apply a background fill.",
		"SETTINGS",
		"--theme", "Theme to use for syntax highlighting",
		"BORDER",
		"--border.color", "Border color.",
		"SHADOW",
		"--shadow.blur", "Shadow Gaussian Blur.",
		"FONT",
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
			output: "artichoke-base.svg",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--config", "test/configurations/full.json"},
			output: "artichoke-full.svg",
		},
		// {
		// 	flags:  []string{"--execute", "eza --no-time --no-filesize --no-user -l"},
		// 	output: "eza.svg",
		// },
		// {
		// 	flags:  []string{"--execute", "layout", "--width", "820", "--config", "full"},
		// 	output: "ansi.svg",
		// },
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--language", "haskell"},
			output: "haskell.svg",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--theme", "dracula"},
			output: "dracula.svg",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--border.radius", "8"},
			output: "border-radius.svg",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--border.radius", "8", "--window"},
			output: "window.svg",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--border.radius", "8", "--window", "--border.width", "1"},
			output: "border-width.svg",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--border.radius", "8", "--window", "--border.width", "1", "--padding", "30,50,30,30"},
			output: "padding.svg",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--border.radius", "8", "--window", "--border.width", "1", "--padding", "30,50,30,30", "--margin", "50,60,100,60"},
			output: "margin.svg",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--config", "full"},
			output: "shadow.svg",
		},
		{
			input:  "test/input/artichoke.hs",
			flags:  []string{"--lines", "4,8", "--show-line-numbers"},
			output: "lines.svg",
		},
	}

	err := os.RemoveAll("test/output")
	if err != nil {
		t.Fatal("unable to remove output files")
	}
	err = os.Mkdir("test/output", 0755)
	if err != nil {
		t.Fatal("unable to create output directory")
	}

	for _, tc := range tests {
		t.Run(strings.Join(tc.flags, " "), func(t *testing.T) {
			out := bytes.Buffer{}
			args := append([]string{tc.input}, tc.flags...)
			args = append(args, "--output", "test/output/"+tc.output)
			cmd := exec.Command(binary, args...)
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				t.Log(err)
				t.Log(out.String())
				t.Fatal("unexpected error")
			}
			want, err := os.ReadFile("test/golden/" + tc.output)
			if err != nil {
				t.Fatal("no golden file for:", "test/golden/"+tc.output)
			}
			got, err := os.ReadFile("test/output/" + tc.output)
			if err != nil {
				t.Fatal("no output file for:", "test/output/"+tc.output)
			}
			if string(want) != string(got) {
				// t.Log(udiff.Unified("want", "got", string(want), string(got)))
				t.Fatalf("test/golden/%s != test/output/%s", tc.output, tc.output)
			}
		})
	}
}
