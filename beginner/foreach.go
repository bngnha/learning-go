package main

import "fmt"

var strarrs = []string{"hello", "world", "from", "go"}

func foreachFunc() {
	for index, element := range strarrs {
		fmt.Println(index, element, strarrs[index])
	}
}

func main() {
	foreachFunc()
}
