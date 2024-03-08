package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/aymanbagabas/go-udiff"
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

func TestFreezeFlags(t *testing.T) {
	output := "artichoke-test.svg"
	defer os.Remove(output)

	cmd := exec.Command(binary, "examples/artichoke.hs", "-o", output)
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

func TestFreezeOutput(t *testing.T) {
	tests := []struct {
		input  string
		config string
		output string
	}{{
		input:  "examples/artichoke.hs",
		config: "test/configurations/base.json",
		output: "artichoke.svg",
	}}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			out := bytes.Buffer{}
			cmd := exec.Command(binary, tc.input, "--config", tc.config, "--output", "test/output/"+tc.output)
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
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
				t.Log(udiff.Unified("want", "got", string(want), string(got)))
				t.Fatalf("test/golden/%s != test/output/%s", tc.output, tc.output)
			}
		})
	}
}
