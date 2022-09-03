package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Returns a manga which have all information about it
func FetchMangaData(link string) (Manga, error) {
	c := colly.NewCollector()

	var errorPage error = nil
	manga := Manga{}
	manga.Chapters = make(map[float64]string)
	page := 1

	// Entering on a site
	c.OnRequest(func(r *colly.Request) {
		const defaultPage = "comienzo"
		fmt.Println(r.URL.String())
		if strings.Contains(r.URL.String(), defaultPage) {
			errorPage = errors.New("data not found")
			fmt.Println("Error detected!!")
		}
	})
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
		if errorPage != nil {
			currPage, _ := strconv.Atoi(e.Text)
			if currPage == page {
				page++
				c.Visit(e.Attr("href"))
			}
		}
	})
	c.Visit(link)

	if manga.IsEmpty() {
		errorPage = errors.New("data not found")
	}

	return manga, errorPage
}

// Searches a manga on google and returns a slice of matching results
func Search(mangaName string) (links []string, err error) {
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
	if len(links) == 0 {
		err = errors.New("manga not found")
	}
	return
}

// Returns all manga pages from a chapter
func FetchImagesByName(name string, chapter int) (pages []string, err error) {
	const notAllowedSymbols = ":!@#$%&"
	const static = "https://img.meusmangas.net//image"
	/*
	 * Formats manga name to another one whose website is able to reach
	 * Example: From: "Huge: Stupid Large    NAME" to "huge-stupid-large-name"
	 */
	nameFormated := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	for _, symbol := range notAllowedSymbols {
		if strings.Contains(nameFormated, string(symbol)) {
			nameFormated = strings.ReplaceAll(nameFormated, string(symbol), "")
		}
	}
	// Visits first page
	i := 1
	req := fmt.Sprintf("%s/%s/%d/%d.jpg", static, nameFormated, chapter, i)
	resp, er := http.Get(req)
	err = er
	// Returns with the error
	if er != nil {
		return
	}
	// If ok, visits the remain of pages
	defer resp.Body.Close()
	for resp.StatusCode == http.StatusOK {
		pages = append(pages, req)
		i++
		req = fmt.Sprintf("%s/%s/%d/%d.jpg", static, nameFormated, chapter, i)
		log.Println(req)
		resp, err = http.Get(req)
	}
	return
}

func UselessLine() string {
	return "==================================================="
}
