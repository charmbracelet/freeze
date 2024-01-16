package main

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestConfig(t *testing.T) {
	dir := "configurations"

	entries, err := configs.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 2 {
		t.Fatal(entries)
	}

	for _, entry := range entries {
		f, err := configs.Open(dir + "/" + entry.Name())
		if err != nil {
			t.Fatal(err)
		}
		_, err = kong.JSON(f)
		if err != nil {
			t.Fatal(err)
		}
	}
}
