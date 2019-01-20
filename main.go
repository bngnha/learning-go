package main

import (
	"fmt"

	a "github.com/bngnha/learn-golang/advanced"
	b "github.com/bngnha/learn-golang/beginner"
)

func main() {
	//====BEGINNER====
	b.Hello()
	fmt.Println("===============")

	//b.PrintVar()

	// Write content to file
	//b.WriteToFile("test.txt", "Hello World")

	// Input from stdio
	//b.Input()

	// Read input from sdtio
	//b.Scan()

	// Array
	//b.WorkWithArray()

	// Interface
	//b.FunctionWithInterface()

	// Work with http
	//b.WorkWithHttp()

	//====ADVANCED====
	//a.WorkWithPointer()

	//a.PlayWithVarDic()
	//a.PlayWithCallback()
	a.PlayWithClosure()
}
