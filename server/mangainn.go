package server

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/web"
)

type Mangainn struct{}

func (m *Mangainn) GetMangaDetail(mangaURL string) (manga web.Manga, err error) {
	c := colly.NewCollector()
	index := 0
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
	c.OnHTML("h5.widget-heading", func(e *colly.HTMLElement) {
		if len(manga.Title) == 0 {
			manga.Title = e.Text
		}
	})
	// Fetch manga status, tags and suthor
	/*
		Because there is no pattern in "dd" tag, I noticed that if I put an "index" when data were scraped
		It would be possible to have a pattern there
	*/
	c.OnHTML("div.col-md-8 dl.dl-horizontal dd", func(e *colly.HTMLElement) {
		switch index {
		case 1:
			manga.Status = e.Text
		case 2:
			tagsText := strings.TrimSpace(strings.ReplaceAll(e.Text, " ", ""))
			tags := strings.Split(tagsText, ",")
			manga.Tags = tags
		case 4:
			manga.Author = strings.TrimSpace(e.Text)
		}
		index++
	})
	// Fetch manga thumbnail
	c.OnHTML("div.col-md-4 img.img-responsive", func(e *colly.HTMLElement) {
		manga.Thumbnail = e.Attr("src")
	})
	// Fetch manga description
	c.OnHTML("div.note", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})
	// Fetch manga chapters
	c.OnHTML("ul.chapter-list li a span.val", func(e *colly.HTMLElement) {
		ch := strings.Split(strings.ReplaceAll(e.Text, " ", ""), "-")
		manga.Chapters = append(manga.Chapters, ch[len(ch)-1])
	})

	c.Visit(mangaURL)
	manga.TotalChapters = len(manga.Chapters)
	return
}
