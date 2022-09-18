package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/web"
)

type Mangainn struct{}

func (m *Mangainn) GetMangaDetail(mangaURL string) (manga web.Manga, statusCode int) {
	c := colly.NewCollector()
	index := 0

	// Detect errors on page
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Something went wrong:", err)
		statusCode = r.StatusCode
	})
	// Get response
	c.OnResponse(func(r *colly.Response) {
		statusCode = r.StatusCode
	})
	// Fetch manga title
	c.OnHTML("h5.widget-heading", func(e *colly.HTMLElement) {
		if len(manga.Title) == 0 {
			manga.Title = e.Text
		}
	})
	// Fetch manga status, tags and author
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
			genres := strings.Split(tagsText, ",")
			manga.Genres = genres
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
		manga.Summary = e.Text
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

func (m *Mangainn) GetMangaPages(mangaTitle, chapter string) (ch web.Chapter, err error) {
	// Because this site alraedy give the total of pages, I can get all of them without
	// using this slow function: return web.FetchImagesByName(pathImages[MANGAINN], mangaTitle, chapter)
	c := colly.NewCollector()
	path := pathImages[db.MANGAINN]
	fmtTitle := web.FormatedTitle(mangaTitle)
	link := fmt.Sprintf("%s/%s/%s/all-pages", path, fmtTitle, chapter)

	ch.Title = mangaTitle
	ch.Value = chapter

	// Entering on a site
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting:", r.URL)
	})
	// Detect errors on page
	c.OnError(func(_ *colly.Response, er error) {
		log.Println("Something went wrong:", er)
		err = er
	})
	// Fetch manga pages
	c.OnHTML("img.img-responsive[src]", func(e *colly.HTMLElement) {
		ch.Pages = append(ch.Pages, e.Attr("src"))
	})

	c.Visit(link)

	ch.TotalPages = len(ch.Pages)

	/*
		When Mangainn doesn't find the manga, the c.OnHTML only returns
		a single uri indicating the manga's poster
	*/
	if ch.TotalPages == 1 && strings.Contains(ch.Pages[0], "posters") {
		err = web.ErrMangaNotFound
	}

	return
}

func (m *Mangainn) FetchAllMangaByLetter(letter string) (links []string) {
	link := "https://www.mangainn.net/manga-list/" + letter
	c := colly.NewCollector()

	// Fetch manga sinopse
	c.OnHTML("ul.manga-list li a.manga-info-qtip", func(e *colly.HTMLElement) {
		links = append(links, e.Attr("href"))
	})

	c.Visit(link)
	return
}

func (m *Mangainn) GetManga(mangaTitle string) (manga web.Manga, err error) {
	return db.GetManga(db.MANGAINN, mangaTitle)
}
