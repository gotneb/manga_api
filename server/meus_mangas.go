/*
Data scraped from site: https://meusmangas.net/comienzo
*/
package server

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/web"
)

var ErrMangaNotFound = errors.New("manga has't been found")
var ErrDataNotCollected = errors.New("manga data has't been collected")
var ErrPageNotFound = errors.New("manga page not found")
var ErrWebsiteWithoutContentData = errors.New("manga page null")
var ErrFirstPageNotFound = errors.New("first page not found")

type MeusMangas struct{}

func (m *MeusMangas) GetMangaDetail(mangaURL string) (manga web.Manga, err error) {
	c := colly.NewCollector()

	var errorPage error = nil
	manga = web.Manga{}
	collectData := true
	page := 1

	// Detect errors on page
	c.OnError(func(_ *colly.Response, err error) {
		errorPage = err
		log.Println("Something went wrong:", err)
	})
	// Fetch manga sinopse
	c.OnHTML("div.sinopse-page", func(e *colly.HTMLElement) {
		if desc := e.Text; len(desc) > 1 {
			manga.Description = desc
		}
	})
	// Fetch manga thumbnail
	c.OnHTML("img.hGq41", func(e *colly.HTMLElement) {
		if thumb := e.Attr("src"); collectData && len(thumb) > 1 {
			manga.Thumbnail = thumb
		}
	})
	// Fetch title manga
	c.OnHTML("h1.kw-title", func(e *colly.HTMLElement) {
		if title := e.Text; collectData && len(title) > 1 {
			manga.Title = e.Text
		}
	})
	// Fetch tags manga
	c.OnHTML("a.widget-btn", func(e *colly.HTMLElement) {
		if tag := e.Text; collectData && len(tag) > 1 {
			manga.Tags = append(manga.Tags, tag)
		}
	})
	// Fetch total chapters
	c.OnHTML("div.jVBw-infos > span", func(e *colly.HTMLElement) {
		number := strings.Split(e.Text, " ")[0]
		if collectData && len(number) > 0 {
			chapters, err := strconv.Atoi(number)
			if err == nil {
				manga.TotalChapters = chapters
			}
		}
	})
	// Fetch manga status. e.g "Em Andamento", "Finalizado"
	c.OnHTML("div.jVBw-infos span.mdq", func(e *colly.HTMLElement) {
		if sit := e.Text; collectData && len(sit) > 1 {
			manga.Status = e.Text
		}
	})
	// Fetch manga author
	c.OnHTML("div.jVBw-infos div", func(e *colly.HTMLElement) {
		if aut := e.Text; collectData && len(aut) > 1 {
			manga.Author = e.Text[1:]
			manga.Author = strings.TrimSpace(manga.Author)
		}
	})
	// Fetch chapters avaliable to read
	c.OnHTML("a.link-dark", func(e *colly.HTMLElement) {
		if len(e.Attr("title")) > 1 {
			// e.Attr("") returns "ler capitulo N"
			chTitle := strings.Split(e.Attr("title"), " ")[2]
			// For unknown reason, the chapter "0", isn't showed on the site
			if chTitle == "" {
				chTitle = "0"
			}
			// ==============================================================
			manga.Chapters = append(manga.Chapters, chTitle)
		}
	})
	// Visit all manga pages
	c.OnHTML("ul.content-pagination li a", func(e *colly.HTMLElement) {
		if errorPage == nil {
			currPage, _ := strconv.Atoi(e.Text)
			if currPage == page {
				collectData = false
				page++
				c.Visit(e.Attr("href"))
			}
		}
	})
	c.Visit(mangaURL)

	if manga.IsEmpty() {
		manga.Show()
		panic(ErrDataNotCollected)
	}
	manga.TotalChapters = len(manga.Chapters)
	return manga, errorPage
}

func (m *MeusMangas) GetMangaPages(mangaTitle, chapter string) (ch web.Chapter, err error) {
	return web.FetchImagesByName(pathImages[MEUS_MANGAS], mangaTitle, chapter)
}
