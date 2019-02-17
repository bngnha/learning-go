package crawler

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
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
	Content     string   `bson:"content"`
	Author      string   `bson:"author_name"`
	AuthorImage string   `bson:"author_image"`
	Tags        []string `bson:"tags"`
	LikeNo      int      `bson:"liked"`
}

func (q *quote) String() string {
	return fmt.Sprintf("%s ― %s\n Tags: %v \t Like: %s\n %s\n\n", q.Content, q.Author, q.Tags, strconv.Itoa(q.LikeNo), q.AuthorImage)
}

// GetQuotes function
func GetQuotes() {
	c := colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
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

// Shopify function
func Shopify() {
	// Array containing all the known URLs in a sitemap
	knownUrls := []string{}

	// Create a Collector specifically for Shopify
	c := colly.NewCollector(colly.AllowedDomains("www.shopify.com"))

	// Create a callback on the XPath query searching for the URLs
	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		knownUrls = append(knownUrls, e.Text)
	})

	// Start the collector
	c.Visit("https://www.shopify.com/sitemap.xml")

	fmt.Println("All known URLs:")
	for _, url := range knownUrls {
		fmt.Println("\t", url)
	}
	fmt.Println("Collected", len(knownUrls), "URLs")
}

// CoinMarketCap function
func CoinMarketCap() {
	fName := "cryptocoinmarketcap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Name", "Symbol", "Price (USD)", "Volume (USD)", "Market capacity (USD)", "Change (1h)", "Change (24h)", "Change (7d)"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("#currencies-all tbody tr", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText(".currency-name-container"),
			e.ChildText(".col-symbol"),
			e.ChildAttr("a.price", "data-usd"),
			e.ChildAttr("a.volume", "data-usd"),
			e.ChildAttr(".market-cap", "data-usd"),
			e.ChildText(".percent-1h"),
			e.ChildText(".percent-24h"),
			e.ChildText(".percent-7d"),
		})
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
