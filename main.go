package main

import (
	"fmt"
	//"os"

	//a "github.com/bngnha/learn-golang/advanced"
	//b "github.com/bngnha/learn-golang/beginner"
	c "github.com/bngnha/learn-golang/crawler"
)

func main() {
	//====BEGINNER====
	//b.Hello()
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
	//a.PlayWithClosure()

	//a.PlayWithOOPBasic()

	//====CRAWLER====
	/*
		if len(os.Args) > 0 {
			if os.Args[1] == "tag" {
				c.GetTags()
			} else if os.Args[1] == "quote" {
				c.GetQuotes()
			} else if os.Args[1] == "image" {
				c.DownloadImages()
			}
		}
	*/
	c.GetQuotes()
	//c.Shopify()
	//c.CoinMarketCap()
}
