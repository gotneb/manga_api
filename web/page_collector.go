package web

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly"
)

// Get all manga links hosted on meus mangas at "section"
// Section -> "a" to "z"
func FetchPages(section string) (links []string) {
	link := "https://meusmangas.net/lista-de-mangas/" + section
	page := 1
	c := colly.NewCollector()

	// Fetch manga sinopse
	c.OnHTML("ul.seriesList li", func(e *colly.HTMLElement) {
		fmt.Println(e.ChildAttr("a", "href"))
		links = append(links, e.ChildAttr("a", "href"))
	})
	// Pass trough seciton page
	c.OnHTML("ul.content-pagination li a", func(e *colly.HTMLElement) {
		currPage, _ := strconv.Atoi(e.Text)
		if currPage == page {
			page++
			c.Visit(e.Attr("href"))
		}
	})

	c.Visit(link)
	return
}
