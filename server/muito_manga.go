/*
Data scraped from site: https://muitomanga.com/
*/
package server

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/web"
)

type MuitoManga struct{}

func (m *MuitoManga) GetMangaDetail(mangaURL string) (manga web.Manga, err error) {
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
	c.OnHTML("div.widget-title h1", func(e *colly.HTMLElement) {
		// Remove useless "Ler " from manga title
		manga.Title = e.Text
	})
	// Fetch manga description
	c.OnHTML("div.boxAnimeSobreLast p", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})
	// Fetch manga thumbnail
	c.OnHTML("div.capaMangaInfo a img", func(e *colly.HTMLElement) {
		manga.Thumbnail = e.Attr("data-src")
	})
	// Fetch manga genres
	c.OnHTML("div.boxAnimeSobreLast ul.lancamento-list li a", func(e *colly.HTMLElement) {
		manga.Tags = append(manga.Tags, e.Text)
	})
	// Fetch manga author and manga status
	c.OnHTML("div.boxAnimeSobreLast span.series_autor2", func(e *colly.HTMLElement) {
		manga.Author = e.Text
		/*
			Fetch manga status
			Well, when I get Author, the status comes along it, beucase of this I handle it right here -.-
			E.g: "Finalizado[author name]", such a mess!
		*/
		const finished = "Finalizado"
		const ongoing = "Em andamento"
		if strings.Contains(manga.Author, finished) {
			manga.Status = finished
			manga.Author = strings.ReplaceAll(manga.Author, finished, "")
		} else {
			manga.Status = ongoing
		}
	})
	// Fetch manga chapters
	c.OnHTML("div.manga-chapters div.single-chapter a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Capítulo") {
			single := strings.ReplaceAll(e.Text, "Capítulo #", "")
			manga.Chapters = append(manga.Chapters, single)
		}
	})

	c.Visit(mangaURL)
	manga.TotalChapters = len(manga.Chapters)
	return
}

func (m *MuitoManga) GetMangaPages(mangaTitle, chapter string) (ch web.Chapter, err error) {
	return web.FetchImagesByName(pathImages[MUITO_MANGA], mangaTitle, chapter)
}
