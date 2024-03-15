package main

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// runInPty opens a new pty and runs the given command in it.
// The returned file is the pty's file descriptor and must be closed by the
// caller.
func (cfg Config) runInPty(c *exec.Cmd) (*os.File, error) {
	return pty.StartWithSize(c, &pty.Winsize{
		Cols: 80,
		Rows: 10,
		X:    uint16(cfg.Width),
	})
}
