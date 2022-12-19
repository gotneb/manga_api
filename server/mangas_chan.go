package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/web"
)

type MangasChan struct{}

func (m *MangasChan) GetMangaDetail(mangaURL string) (manga web.Manga, statusCode int) {
	c := colly.NewCollector()
	statusCode = http.StatusOK
	counter := 1

	isFirstParagraph := true

	// Detect errors on page
	c.OnError(func(r *colly.Response, err error) {
		statusCode = r.StatusCode
		log.Println("Something went wrong:", err)
	})
	// Fetch title manga
	c.OnHTML("div.info-desc.bixbox h1.entry-title", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})
	// Fetch manga thumbnail
	c.OnHTML("div.thumb img", func(e *colly.HTMLElement) {
		manga.Thumbnail = e.Attr("data-lazy-src")
	})
	// Fetch manga summary
	c.OnHTML("div.entry-content.entry-content-single p", func(e *colly.HTMLElement) {
		if isFirstParagraph {
			manga.Summary = e.Text
			isFirstParagraph = !isFirstParagraph
		} else {
			manga.Summary += "\n\n" + e.Text
		}
	})
	// Fetch manga genres
	c.OnHTML("div.wd-full span.mgen a", func(e *colly.HTMLElement) {
		manga.Genres = append(manga.Genres, e.Text)
	})
	// Fetch manga author and status!
	c.OnHTML("div.tsinfo.bixbox div.imptdt i", func(e *colly.HTMLElement) {
		if counter == 1 {
			manga.Status = e.Text
		}
		if counter == 3 {
			manga.Author = e.Text
		}
		counter++
	})
	// Fetch manga chapters
	c.OnHTML("div#chapterlist.eplister ul.clstyle li", func(e *colly.HTMLElement) {
		manga.Chapters = append(manga.Chapters, e.Attr("data-num"))
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

func (m *MangasChan) FetchAllMangaByLetter(letter string) (links []string) {
	link := fmt.Sprintf("https://mangaschan.com/lista-de-a-z/page/1/?show=%v", strings.ToUpper(letter))
	page := 1
	c := colly.NewCollector()

	// Fetch manga sinopse
	c.OnHTML("div.listo div.bs a", func(e *colly.HTMLElement) {
		links = append(links, e.Attr("href"))
	})
	// Pass trough seciton page
	c.OnHTML("div.pagination a", func(e *colly.HTMLElement) {
		page++
		c.Visit(e.Attr("href"))
	})

	c.Visit(link)
	return
}

func (m *MangasChan) GetMangaPages(mangaTitle, chapter string) (ch web.Chapter, err error) {
	return web.FetchImagesByName(pathImages[db.MANGAS_CHAN], mangaTitle, chapter)
}

func (m *MangasChan) GetManga(mangaTitle string) (manga web.Manga, err error) {
	return db.GetManga(db.MANGAS_CHAN, mangaTitle)
}
