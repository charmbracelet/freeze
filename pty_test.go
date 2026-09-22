package main

import (
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestExecuteCommandSupportsPipes is a regression test for
// https://github.com/charmbracelet/freeze/issues/67: --execute used to
// tokenize the command string with shellwords and exec the first token
// directly, so shell operators like pipes were passed through as literal
// argv tokens to that first command instead of being interpreted.
func TestExecuteCommandSupportsPipes(t *testing.T) {
	execute := "echo hello world | grep world"
	if runtime.GOOS == "windows" {
		execute = "echo hello world | findstr world"
	}

	out, err := executeCommand(Config{
		Execute:        execute,
		ExecuteTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("executeCommand returned an error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "world") {
		t.Fatalf("expected output to contain %q, got: %q", "world", out)
	}
}

func TestShellCommand(t *testing.T) {
	shell, flag := shellCommand()
	if runtime.GOOS == "windows" {
		if shell != "cmd" || flag != "/c" {
			t.Fatalf("expected cmd /c on windows, got %s %s", shell, flag)
		}
		return
	}
	if flag != "-c" {
		t.Fatalf("expected -c flag on unix, got %s", flag)
	}
	if shell == "" {
		t.Fatal("expected a non-empty shell")
	}
}
