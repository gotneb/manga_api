package server

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/web"
)

type KissManga struct{}

func (k *KissManga) GetMangaDetail(mangaURL string) (manga web.Manga, err error) {
	c := colly.NewCollector()

	// Entering on a site
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting:", r.URL)
	})
	// Detect errors on page
	c.OnError(func(_ *colly.Response, er error) {
		log.Println("Something went wrong:", err)
		err = er
	})
	// Fetch manga title
	c.OnHTML("h1.page-title", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})
	// Fetch manga status
	c.OnHTML("span.series-status", func(e *colly.HTMLElement) {
		manga.Status = e.Text
	})
	// Fetch manga author
	c.OnHTML("span#first_episode a small", func(e *colly.HTMLElement) {
		manga.Author = e.Text
	})
	// Fetch manga thumbnail
	c.OnHTML("img.series-profile-thumb", func(e *colly.HTMLElement) {
		manga.Thumbnail = "https://readm.org" + e.Attr("src")
	})
	// Fetch manga genres
	c.OnHTML("div.item a[title]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), "category") {
			manga.Tags = append(manga.Tags, e.Attr("title"))
		}
	})
	// Fetch manga description
	c.OnHTML("div.series-summary-wrapper p", func(e *colly.HTMLElement) {
		if len(e.Text) > 0 {
			manga.Description = e.Text
		}
	})

	c.Visit(mangaURL)
	manga.Show()
	return
}
