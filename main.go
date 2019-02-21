package main

import (
	"fmt"
	//"os"
	//c "github.com/bngnha/learn-golang/crawler"
	"github.com/bngnha/learn-golang/videos"
)

func main() {
	fmt.Println("===============")
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
	//c.GetQuotes()
	//c.Shopify()
	//c.CoinMarketCap()
	videos.ReupYt()

	fmt.Scanln()
	fmt.Println("Done")
}
