package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmdOutput, err := exec.Command("cmd").Output()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Cmd output string:", string(cmdOutput))
}
