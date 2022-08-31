package web

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func SetupCollect(links *[]string) (mangas []Manga) {
	c := colly.NewCollector()

	manga := Manga{}
	manga.Chapters = make(map[float64]string)
	page := 1

	// Entering on site
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	// Fetch manga sinopse
	c.OnHTML("div.sinopse-page", func(e *colly.HTMLElement) {
		if desc := e.Text; len(desc) > 1 {
			manga.Description = desc
		}
	})
	// Fetch description manga
	c.OnHTML("img.hGq41", func(e *colly.HTMLElement) {
		if thumb := e.Attr("src"); len(thumb) > 1 {
			manga.Thumbnail = thumb
		}
	})
	// Fetch title manga
	c.OnHTML("h1.kw-title", func(e *colly.HTMLElement) {
		if title := e.Text; len(title) > 1 {
			manga.Title = e.Text
		}
	})
	// Fetch tags manga
	c.OnHTML("a.widget-btn", func(e *colly.HTMLElement) {
		if tag := e.Text; len(tag) > 1 {
			manga.Tags = append(manga.Tags, tag)
		}
	})
	// Fetch total chapters
	c.OnHTML("div.jVBw-infos > span", func(e *colly.HTMLElement) {
		number := strings.Split(e.Text, " ")[0]
		if len(number) > 0 {
			chapters, err := strconv.Atoi(number)
			if err == nil {
				manga.TotalChapters = chapters
			}
		}
	})
	// Fetch manga situation
	c.OnHTML("div.jVBw-infos span.mdq", func(e *colly.HTMLElement) {
		if sit := e.Text; len(sit) > 1 {
			manga.Situation = e.Text
		}
	})
	// Fetch manga author
	c.OnHTML("div.jVBw-infos div", func(e *colly.HTMLElement) {
		if aut := e.Text; len(aut) > 1 {
			manga.Author = e.Text
		}
	})
	// Fetch chapters avaliable to read
	c.OnHTML("a.link-dark", func(e *colly.HTMLElement) {
		if len(e.Attr("title")) > 1 {
			// e.Attr("") returns "ler capitulo N"
			chTitle := strings.Split(e.Attr("title"), " ")[2]
			chNumber, _ := strconv.ParseFloat(chTitle, 32)
			manga.Chapters[chNumber] = e.Attr("href")
		}
	})
	// Visit all manga pages
	c.OnHTML("ul.content-pagination li a", func(e *colly.HTMLElement) {
		currPage, _ := strconv.Atoi(e.Text)
		if currPage == page {
			page++
			c.Visit(e.Attr("href"))
		}
	})
	// Finished
	c.OnScraped(func(r *colly.Response) {
		if manga.Title != "" {
			// manga.Show()
			mangas = append(mangas, manga)
			// Clear data for a new search :D
			manga = Manga{}
			manga.Chapters = make(map[float64]string)
			page = 1
		}
	})

	for _, link := range *links {
		c.Visit(link)
	}
	return
}

// Searches a manga on google and returns a slice of matching results
func Search(mangaName string) (links []string) {
	c := colly.NewCollector()
	lookFor := mangaName + "+meus+mangas"
	staticLink := "https://www.google.com/search?q=" + lookFor
	// Every time it enters in a website
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	// On every <a> whose has "href" attribute
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		defAdresses := [2]string{
			"https://meusmangas.net/manga/mango/",
			"https://meusmangas.net/manga/hd/",
		}
		url := e.Attr("href")
		if strings.Contains(url, defAdresses[0]) || strings.Contains(url, defAdresses[1]) {
			url = url[7:strings.Index(url, "&")]
			links = append(links, url)
		}
	})

	c.Visit(staticLink)
	return
}

func FetchImagesByName(name string, chapter int) {
	const notAllowedSymbols = "!@#$%&"
	const static = "https://img.meusmangas.net//image"
	// From: "Huge: Stupid Large    NAME" to "huge-stupid-large-name"
	req := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	for _, symbol := range notAllowedSymbols {
		if strings.Contains(req, string(symbol)) {
			req = strings.ReplaceAll(req, string(symbol), "")
		}
	}
	req = fmt.Sprintf("%s/%s/%d/1.jpg", static, req, chapter)
	fmt.Printf("%s\nLooking for images...\n", UselessLine())
	fmt.Println(req)
}

func UselessLine() string {
	return "==================================================="
}
