package main

import (
	"fmt"
)

var variable string

func printVar() {
	variable = "Testing..."

	fmt.Println(variable)
}

func main() {
	printVar()
}
