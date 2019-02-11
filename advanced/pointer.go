package main

import "fmt"

func workWithPointer() {
	var value int
	value = 5

	fmt.Printf("Value: %v\n", value)

	plusOne(&value)
	fmt.Printf("Value: %v\n", value)

	plusOne(&value)
	plusOne(&value)
	fmt.Printf("Value: %v\n", value)
}

func plusOne(in *int) {
	*in++
}

func main() {
	workWithPointer()
}
