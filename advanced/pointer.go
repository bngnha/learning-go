package advanced

import "fmt"

func WorkWithPointer() {
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
