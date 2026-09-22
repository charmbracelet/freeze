package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/x/term"
	"github.com/charmbracelet/x/xpty"
)

// shellCommand returns the shell and the flag used to run a command string
// through it, so --execute supports pipes, redirects, &&, and other shell
// operators instead of only a single literal argv.
func shellCommand() (string, string) {
	if runtime.GOOS == "windows" {
		return "cmd", "/c"
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	return shell, "-c"
}

func executeCommand(config Config) (string, error) {
	shell, shellFlag := shellCommand()

	ctx, cancel := context.WithTimeout(context.Background(), config.ExecuteTimeout)
	defer cancel()

	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width = 80
		height = 24
	}

	pty, err := xpty.NewPty(width, height)
	if err != nil {
		return "", fmt.Errorf("could not execute: %w", err)
	}
	defer func() { _ = pty.Close() }()

	cmd := exec.CommandContext(ctx, shell, shellFlag, config.Execute) //nolint: gosec
	if err := pty.Start(cmd); err != nil {
		return "", fmt.Errorf("could not execute: %w", err)
	}

	var out bytes.Buffer
	donec := make(chan struct{})
	go func() {
		_, _ = io.Copy(&out, pty)
		close(donec)
	}()

	err = xpty.WaitProcess(ctx, cmd)
	waitOutput(pty, donec)
	if err != nil {
		return out.String(), fmt.Errorf("could not execute: %w", err)
	}
	return out.String(), nil
}
