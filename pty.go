package main

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// runInPty opens a new pty and runs the given command in it.
// The returned file is the pty's file descriptor and must be closed by the
// caller.
func runInPty(c *exec.Cmd) (*os.File, error) {
	return pty.Start(c)
}
