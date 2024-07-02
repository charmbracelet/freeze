//go:build !windows
// +build !windows

package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/caarlos0/go-shellwords"
	"github.com/creack/pty"
)

// runInPty opens a new pty and runs the given command in it.
// The returned file is the pty's file descriptor and must be closed by the
// caller.
func (cfg Config) runInPty(c *exec.Cmd) (*os.File, error) {
	return pty.StartWithAttrs(c, &pty.Winsize{
		Cols: 80,
		Rows: 10,
		X:    uint16(cfg.Width),
	}, &syscall.SysProcAttr{})
}

func executeCommand(config Config) (string, error) {
	args, err := shellwords.Parse(config.Execute)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.ExecuteTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	pty, err := config.runInPty(cmd)
	if err != nil {
		return "", err
	}
	defer pty.Close()
	var out bytes.Buffer
	go func() {
		_, _ = io.Copy(&out, pty)
	}()

	err = cmd.Wait()
	if err != nil && !config.ExecuteNonZero {
		return "", err
	}
	return out.String(), nil
}
