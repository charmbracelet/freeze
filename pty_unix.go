//go:build !windows

package main

import (
	"time"

	"github.com/charmbracelet/x/xpty"
)

// drainTimeout is how long to wait for the copy goroutine to drain the pty
// after the process exits. The slave end can be held open by grandchildren
// of the executed command, in which case the read never reaches EOF, so the
// wait must be bounded.
const drainTimeout = 200 * time.Millisecond

// waitOutput waits for the copy goroutine to finish. On Unix, the slave end
// of the pty is usually closed when the process exits, so the copy goroutine
// receives an EOF (or EIO) once it has drained all buffered output. If the
// slave is held open past drainTimeout, the pty is closed to unblock the
// read; any output not yet read is discarded.
func waitOutput(pty xpty.Pty, donec <-chan struct{}) {
	select {
	case <-donec:
	case <-time.After(drainTimeout):
		_ = pty.Close()
		<-donec
	}
}
