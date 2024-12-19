package main

import (
	"github.com/muesli/reflow/wordwrap"
	"strings"
)

func SoftWrap(input string, wrapLength int) []bool {
	var wrap []bool
	for _, line := range strings.Split(input, "\n") {
		wrappedLine := wordwrap.String(line, wrapLength)

		for i := range strings.Split(wrappedLine, "\n") {
			if i == 0 {
				// We want line number on the original line
				wrap = append(wrap, true)
			} else {
				// for wrapped line, we do not want line number
				wrap = append(wrap, false)
			}
		}
	}
	return wrap
}
