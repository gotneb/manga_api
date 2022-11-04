package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/web"
)

type MangasChan struct{}

func (m *MangasChan) GetMangaDetail(mangaURL string) (manga web.Manga, statusCode int) {
	c := colly.NewCollector()

	counter := 1
	last := ""

	// Detect errors on page
	c.OnError(func(r *colly.Response, err error) {
		statusCode = r.StatusCode
		log.Println("Something went wrong:", err)
	})
	// Fetch title manga
	c.OnHTML("div.seriestuheader h1.entry-title", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})
	// Fetch manga thumbnail
	c.OnHTML("div.thumb img", func(e *colly.HTMLElement) {
		manga.Thumbnail = e.Attr("data-lazy-src")
	})
	// Fetch manga summary
	c.OnHTML("div.seriestuhead div.entry-content.entry-content-single p", func(e *colly.HTMLElement) {
		manga.Summary = e.Text
	})
	// Fetch manga genres
	c.OnHTML("div.seriestucontr div.seriestugenre a", func(e *colly.HTMLElement) {
		manga.Genres = append(manga.Genres, e.Text)
	})
	// Fetch manga author and status!
	c.OnHTML("div.seriestucontr table.infotable td", func(e *colly.HTMLElement) {
		if counter%2 != 0 {
			last = e.Text
		} else {
			if last == "Status" {
				manga.Status = e.Text
			} else if last == "Autor" {
				manga.Author = e.Text
			}
		}
		counter++
	})
	// Fetch manga chapters
	c.OnHTML("div.bixbox.bxcl.epcheck div#chapterlist.eplister ul li div.chbox div.eph-num a span.chapternum", func(e *colly.HTMLElement) {
		ch := strings.Split(e.Text, " ")[1]
		manga.Chapters = append(manga.Chapters, ch)
	})

	c.Visit(mangaURL)

	manga.TotalChapters = len(manga.Chapters)

	return
}

func (m *MangasChan) FetchAllMangaByPageNumber(page int) (links []string) {
	init := "https://mangaschan.com/mangas/page/" + fmt.Sprint(page) + "/"
	c := colly.NewCollector()

	c.OnHTML("div.page div.listupd.cp div.bs.styletere div.bsx a", func(e *colly.HTMLElement) {
		links = append(links, e.Attr("href"))
	})

	c.Visit(init)

	for _, link := range links {
		log.Println(link)
	}
	return
}

func (m *MangasChan) GetMangaPages(mangaTitle, chapter string) (ch web.Chapter, err error) {
	return web.FetchImagesByName(pathImages[db.MANGAS_CHAN], mangaTitle, chapter)
}
