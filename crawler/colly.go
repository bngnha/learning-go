package crawler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var (
	quoteString   string         = "https://www.goodreads.com/quotes"
	contentRegexp *regexp.Regexp = regexp.MustCompile("“(.+?)”")
)

type Quote struct {
	Content string
	Author  string
	Tags    []string
	LikeNo  int
}

func (q *Quote) String() string {
	return fmt.Sprintf("%s ― %s\n Tags: %v \t Like: %s\n\n", q.Content, q.Author, q.Tags, strconv.Itoa(q.LikeNo))
}

func GetQuotes() {
	c := colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains("www.goodreads.com"),
	)
	var quotes []Quote

	c.OnHTML(".quoteDetails", func(e *colly.HTMLElement) {
		res := contentRegexp.FindAllStringSubmatch(e.ChildText("div.quoteText"), -1)

		if len(res) < 1 {
			return
		}

		if len(res[0]) < 1 {
			return
		}

		likeNoStrs := strings.Split(e.ChildText("a.smallText"), " ")
		likeNo, err := strconv.Atoi(likeNoStrs[0])
		if err != nil {
			return
		}

		tags := []string{}
		tagStr := e.ChildText("div.smallText")
		if tagStr != "" {
			tagStrs := strings.Split(tagStr, ":")
			re := regexp.MustCompile(`\r?\n`)
			tags = strings.Split(re.ReplaceAllString(tagStrs[1], " "), ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		}

		quote := Quote{
			Content: res[0][0],
			Author:  e.ChildText(".authorOrTitle"),
			Tags:    tags,
			LikeNo:  likeNo,
		}
		fmt.Print(quote.String())

		quotes = append(quotes, quote)
	})

	// next page
	c.OnHTML(".next_page", func(e *colly.HTMLElement) {
		// limit to download first page
		if len(quotes) < 0 {
			e.Request.Visit(e.Attr("href"))
		}
	})

	fmt.Println("Launching Scraper !\n\n")
	c.Visit(quoteString)
}

func DownloadImages() {

}
