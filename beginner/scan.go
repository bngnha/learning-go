package beginner

import "fmt"

func Scan() {
	var s string
	fmt.Println("Please insert a string and press enter!")
	fmt.Scan(&s)
	fmt.Printf("Read string \"%v\" from stdin\n", s)
}
