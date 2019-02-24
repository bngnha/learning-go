package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	pid := os.Getpid()
	str := strconv.Itoa(pid)
	fmt.Println(pid)
	fmt.Println(str)

	exec.Command("kill", "-9", str).Output()
}
