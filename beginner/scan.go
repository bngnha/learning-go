package main

import "fmt"

func scanFunc() {
	var s string
	fmt.Println("Please insert a string and press enter!")
	fmt.Scan(&s)
	fmt.Printf("Read string \"%v\" from stdin\n", s)
}

func main() {
	scanFunc()
}
