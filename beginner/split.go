package main

import (
	"fmt"
	"strings"
)

var (
	varbl    string
	strarray []string
)

// Split function
func split() {
	varbl = "Hello World"
	fmt.Println(varbl)

	strarray = strings.Split(varbl, " ")

	for i := 0; i < len(strarray); i++ {
		fmt.Println(strarray[i])
	}
}

func main() {
	split()
}
