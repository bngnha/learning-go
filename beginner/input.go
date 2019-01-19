package beginner

import (
	"bufio"
	"fmt"
	"os"
)

func Input() {
	if 1 < len(os.Args) {
		fmt.Print(len(os.Args) - 1)
		fmt.Println("Arguments provided")
	}

	fmt.Println("Enter a string...")
	bio := bufio.NewReader(os.Stdin)
	line, _, err := bio.ReadLine()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, char := range string(line) {
			fmt.Println(char)
		}
	}
}
