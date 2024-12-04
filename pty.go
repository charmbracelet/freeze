package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/caarlos0/go-shellwords"
	"github.com/charmbracelet/x/term"
	"github.com/charmbracelet/x/xpty"
)

func executeCommand(config Config) (string, error) {
	args, err := shellwords.Parse(config.Execute)
	if err != nil {
		return "", err //nolint: wrapcheck
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ExecuteTimeout)
	defer cancel()

	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width = 80
		height = 24
	}

	pty, err := xpty.NewPty(width, height)
	if err != nil {
		return "", err
	}
	defer func() { _ = pty.Close() }()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...) //nolint: gosec
	if err := pty.Start(cmd); err != nil {
		return "", err
	}

	var out bytes.Buffer
	var errorOut bytes.Buffer
	go func() {
		_, _ = io.Copy(&out, pty)
		errorOut.Write(out.Bytes())
	}()

	if err := xpty.WaitProcess(ctx, cmd); err != nil {
		return errorOut.String(), err //nolint: wrapcheck
	}
	return out.String(), nil
}
