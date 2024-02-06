package main

import "strings"

func cut(input string, window []int) string {
	if len(window) == 0 {
		return input
	}
	if len(window) == 1 && window[0] == 0 {
		return input
	}
	if len(window) == 2 && window[0] == 0 && window[1] == -1 {
		return input
	}

	lines := strings.Split(input, "\n")

	start := 0
	end := len(lines)

	switch len(window) {
	case 1:
		if window[0] > 0 {
			start = window[0]
		} else {
			start = len(lines) + window[0] // add negative = subtract
		}
	case 2:
		start = window[0]
		end = window[1]
	}

	start = clamp(start, 0, len(lines))
	end = clamp(end, start, len(lines))

	if start == end && start < len(lines) {
		return lines[start]
	}

	return strings.Join(lines[start:end], "\n")
}

func clamp(n, low, high int) int {
	return min(max(n, low), high)
}
