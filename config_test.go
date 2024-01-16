package main

import (
	"strings"
	"testing"

	"github.com/alecthomas/kong"
)

func TestConfig(t *testing.T) {
	for _, c := range configs {
		_, err := kong.JSON(strings.NewReader(c))
		if err != nil {
			t.Fatal(err)
		}
	}
}
