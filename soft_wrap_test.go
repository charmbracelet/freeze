package main

import (
	"reflect"
	"testing"
)

// Mock the dependency if needed, assuming wordwrap.String works correctly.
func TestSoftWrap(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wrapLength int
		expected   []bool
	}{
		{
			name:       "Single short line, no wrapping",
			input:      "Hello",
			wrapLength: 10,
			expected:   []bool{true},
		},
		{
			name:       "Single long line, wrapping",
			input:      "Hello World, this is a long line",
			wrapLength: 10,
			expected:   []bool{true, false, false, false},
		},
		{
			name:       "Multiple lines, some wrapped",
			input:      "Short\nThis is a long line",
			wrapLength: 10,
			expected:   []bool{true, true, false},
		},
		{
			name:       "Multiple lines, multiple wraps",
			input:      "This is an long line\nThis is an other long line\nThis is the last long line\nShort line",
			wrapLength: 10,
			expected:   []bool{true, false, true, false, false, true, false, false, true},
		},
		{
			name:       "Empty input",
			input:      "",
			wrapLength: 10,
			expected:   []bool{true},
		},
		{
			name:       "Lines with spaces only",
			input:      "    \n      ",
			wrapLength: 5,
			expected:   []bool{true, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SoftWrap(tt.input, tt.wrapLength)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SoftWrap() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
