//go:build windows
// +build windows

package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"syscall"

	"github.com/caarlos0/go-shellwords"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/exp/term/conpty"
	"golang.org/x/sys/windows"
)

func executeCommand(config Config) (string, error) {
	args, err := shellwords.Parse(config.Execute)
	if err != nil {
		log.Error(err)
		printErrorFatal("Something went wrong", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.ExecuteTimeout)
	defer cancel()

	cpty, err := conpty.New(80, 10, 0)
	if err != nil {
		return "", err
	}
	defer cpty.Close()

	pid, proc, err := cpty.Spawn(args[0], args, &syscall.ProcAttr{Env: os.Environ()})
	if err != nil {
		return "", err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		// If we can't find the process via os.FindProcess, terminate the
		// process as that's what we rely on for all further operations on the
		// object.
		if tErr := windows.TerminateProcess(windows.Handle(proc), 1); tErr != nil {
			return "", tErr
		}
		return "", err
	}

	type result struct {
		*os.ProcessState
		error
	}
	donec := make(chan result, 1)
	go func() {
		state, err := process.Wait()
		donec <- result{state, err}
	}()

	ctx, cancelFunc := context.WithTimeout(context.Background(), config.ExecuteTimeout)
	defer cancelFunc()
	var out bytes.Buffer
	go func() {
		_, _ = io.Copy(&out, cpty)
	}()

	select {
	case <-ctx.Done():
		err = windows.TerminateProcess(windows.Handle(proc), 1)
	case r := <-donec:
		err = r.error
	}

	if err != nil && !config.ExecuteNonZero {
		return "", err
	}
	return out.String(), nil
}
