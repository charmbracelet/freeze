package main

import "fmt"

// freeze/issues/14

func main() {
	fmt.Println("This is a really long line that is going to go over the 80 character limit. This is a really long line that is going to go over the 80 character limit. This is a really long line that is going to go over the 80 character limit. This is a really long line that is going to go over the 80 character limit.")
}
