//go:build windows

package main

import (
	"github.com/charmbracelet/x/xpty"
)

// waitOutput waits for the copy goroutine to finish. On Windows, the read
// end of the ConPty output pipe is not closed when the process exits, so
// the copy goroutine stays blocked in ReadFile. Closing the pty unblocks it
// so the process output collected so far can be returned.
func waitOutput(pty xpty.Pty, donec <-chan struct{}) {
	_ = pty.Close()
	<-donec
}
