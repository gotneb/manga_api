/*
Data scraped from site: https://meusmangas.net/comienzo
*/
package server

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/web"
)

var ErrMangaNotFound = errors.New("manga not found")
var ErrDataNotCollected = errors.New("manga data data not collected")
var ErrPageNotFound = errors.New("manga page not found")
var ErrSiteWithoutContentData = errors.New("manga page without content")
var ErrFirstPageNotFound = errors.New("first page not found")

type MeusMangas struct{}

func (m *MeusMangas) GetMangaDetail(mangaURL string) (manga web.Manga, statusCode int) {
	c := colly.NewCollector()

	var errorPage error = nil
	manga = web.Manga{}
	collectData := true
	page := 1

	// Detect errors on page
	c.OnError(func(r *colly.Response, err error) {
		statusCode = r.StatusCode
		errorPage = err
		log.Println("Something went wrong:", err)
	})
	// Get response
	c.OnResponse(func(r *colly.Response) {
		statusCode = r.StatusCode
	})
	// Fetch manga sinopse
	c.OnHTML("div.sinopse-page", func(e *colly.HTMLElement) {
		if summ := e.Text; len(summ) > 1 {
			manga.Summary = summ
		}
	})
	// Fetch manga thumbnail
	c.OnHTML("img.hGq41", func(e *colly.HTMLElement) {
		if thumb := e.Attr("data-lazy-src"); collectData && len(thumb) > 1 {
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
		if genre := e.Text; collectData && len(genre) > 1 {
			manga.Genres = append(manga.Genres, genre)
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
			// Sometimes author's name can be already formated. In this case it's ok return
			if strings.Contains(manga.Author, " ") {
				return
			}
			// Although it can happens: "NameSurname", I wanna this => "Name Surname"
			for index, letter := range manga.Author {
				if index > 0 {
					if unicode.IsUpper(letter) {
						manga.Author = fmt.Sprintf("%s %s", manga.Author[:index], manga.Author[index:])
						return
					}
				}
			}
		}
	})
	// Fetch chapters avaliable to read
	c.OnHTML("a.link-dark", func(e *colly.HTMLElement) {
		if len(e.Attr("title")) > 1 {
			// Get chapter date
			//date := e.ChildText("div.chapter-options span.chapter-date")

			// e.Attr("title") returns "ler capitulo N"
			chTitle := strings.Split(e.Attr("title"), " ")[2]

			// For unknown reason, the chapter "0" isn't displayed on the site
			if chTitle == "" {
				chTitle = "0"
			}

			//sc := web.SimpleChapter{Value: chTitle, Date: date}
			//manga.Chapters = append(manga.Chapters, sc)
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
	return
}

func (m *MeusMangas) GetMangaPages(mangaTitle, chapter string) (ch web.Chapter, err error) {
	return web.FetchImagesByName(pathImages[db.MEUS_MANGAS], mangaTitle, chapter)
}

func (m *MeusMangas) FetchAllMangaByLetter(letter string) (links []string) {
	link := "https://seemangas.com/lista-de-mangas/" + letter
	page := 1
	c := colly.NewCollector()

	// Fetch manga sinopse
	c.OnHTML("ul.seriesList li", func(e *colly.HTMLElement) {
		//fmt.Println(e.ChildAttr("a", "href"))
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

func (m *MeusMangas) GetManga(mangaTitle string) (manga web.Manga, err error) {
	return db.GetManga(db.MEUS_MANGAS, mangaTitle)
}
