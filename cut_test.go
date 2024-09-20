package main

import (
	"testing"
)

const test = `1 package main
2 
3 import (
4   "fmt"
5 )
6
7 func main() {
8    fmt.Println("Hello World")
9 }`

func TestCut(t *testing.T) {
	tests := []struct {
		input    string
		lines    []int
		expected string
	}{
		{"", []int{0, 0}, ""},
		{"Hello World", []int{0, 0}, "Hello World"},
		{"Hello World", []int{0, 1}, "Hello World"},
		{"Hello World", []int{0, 10}, "Hello World"},
		{"Hello World", []int{10, 10}, ""},
		{test, []int{0, 10}, test},
		{test, []int{0, 1}, "1 package main\n2 "},
		{test, []int{1, 2}, "2 \n3 import ("},
		{test, []int{2, 4}, "3 import (\n4   \"fmt\"\n5 )"},
		{test, []int{6}, "7 func main() {\n8    fmt.Println(\"Hello World\")\n9 }"},
		{test, []int{-3}, "7 func main() {\n8    fmt.Println(\"Hello World\")\n9 }"},
		{test, []int{-2}, "8    fmt.Println(\"Hello World\")\n9 }"},
		{test, []int{-1}, "9 }"},
		{test, []int{0, 9}, test},
		{test, []int{6, 6}, "7 func main() {"},
		{test, []int{7, 7}, "8    fmt.Println(\"Hello World\")"},
	}

	for _, test := range tests {
		actual := cut(test.input, test.lines)
		if actual != test.expected {
			t.Log(actual)
			t.Log(test.expected)
			t.Errorf("cut(%s, %v)", test.input, test.lines)
		}
	}
}
