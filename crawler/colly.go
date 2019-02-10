package crawler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var (
	quoteString   = "https://www.goodreads.com/quotes?page=92"
	contentRegexp = regexp.MustCompile("“(.+?)”")
)

type quote struct {
	Content     string
	Author      string
	AuthorImage string
	Tags        []string
	LikeNo      int
}

func (q *quote) String() string {
	return fmt.Sprintf("%s ― %s\n Tags: %v \t Like: %s\n %s\n\n", q.Content, q.Author, q.Tags, strconv.Itoa(q.LikeNo), q.AuthorImage)
}

// GetQuotes function
func GetQuotes() {
	c := colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains("www.goodreads.com"),
	)
	var quotes []quote

	c.OnHTML(".quoteDetails", func(e *colly.HTMLElement) {
		res := contentRegexp.FindAllStringSubmatch(e.ChildText("div.quoteText"), -1)

		if len(res) < 1 {
			return
		}

		if len(res[0]) < 1 {
			return
		}

		// tag
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

		// liked number
		likeNoStrs := strings.Split(e.ChildText("a.smallText"), " ")
		likeNo, err := strconv.Atoi(likeNoStrs[0])
		if err != nil {
			likeNo = 0
		}

		//author image
		// authorImageLink, exist := e.DOM.Find("a.leftAlignedImage > img").Attr("src")
		// if !exist {
		// 	authorImageLink = ""
		// }
		authorImageLink := e.ChildAttr("a.leftAlignedImage > img", "src")

		q := quote{
			Content:     res[0][0],
			Author:      e.ChildText(".authorOrTitle"),
			AuthorImage: authorImageLink,
			Tags:        tags,
			LikeNo:      likeNo,
		}
		fmt.Print(q.String())

		quotes = append(quotes, q)
	})

	// next page
	c.OnHTML(".next_page", func(e *colly.HTMLElement) {
		// limit to download first page
		if len(quotes) < 0 {
			e.Request.Visit(e.Attr("href"))
		}
	})

	fmt.Println("Launching Scraper !")
	c.Visit(quoteString)
}

// DownloadImages function
func DownloadImages() {

}
