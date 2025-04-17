package input //nolint:revive

import (
	"io"
	"os"
)

// ReadFile returns the files content.
func ReadFile(file string) (string, error) {
	b, err := os.ReadFile(file)
	return string(b), err
}

// ReadInput reads some input.
func ReadInput(in io.Reader) (string, error) {
	b, err := io.ReadAll(in)
	return string(b), err
}

// IsPipe returns whether the stdin is a pipe.
func IsPipe(in *os.File) bool {
	stat, err := in.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}
