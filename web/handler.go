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

var ErrMangaNotFound = errors.New("manga has't been found")
var ErrDataNotCollected = errors.New("manga data has't been collected")
var ErrPageNotFound = errors.New("manga page not found")
var ErrWebsiteWithoutContentData = errors.New("manga page null")
var ErrFirstPageNotFound = errors.New("first page not found")

/*
Get all details about the manga.
Link: page link where scrapper can get data
*/
func FetchMangaData(link string) (Manga, error) {
	c := colly.NewCollector()

	var errorPage error = nil
	manga := Manga{}
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
	c.Visit(link)

	if manga.IsEmpty() {
		manga.Show()
		panic(ErrDataNotCollected)
	}
	manga.TotalChapters = len(manga.Chapters)
	return manga, errorPage
}

// Returns all manga pages from a chapter
func FetchImagesByName(name, chapter string) (ch Chapter, err error) {
	const static = "https://img.meusmangas.net//image"
	ch.Title = name
	ch.Value = chapter

	nameFormated := FormatedTitle(name)

	/*
		This is embarrassing, but I was too tired when I wrote this so, it's a messing from here to below (-.-')
	*/

	/*
		This function searches the correct page index. Hence this website sometimes give a "jpeg" or "jpg" file
		Or even, weirdly it starts a page with "2" instead "1"
	*/
	findIndex := func(index int) (string, bool, int, error) {
		resp, _ := http.Get(fmt.Sprintf("%s/%s/%s/%d.jpg", static, nameFormated, chapter, index))
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
			return "jpg", false, index, nil
		}
		resp, _ = http.Get(fmt.Sprintf("%s/%s/%s/0%d.jpg", static, nameFormated, chapter, index))
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
			return "jpg", true, index, nil
		}
		resp, _ = http.Get(fmt.Sprintf("%s/%s/%s/%d.jpeg", static, nameFormated, chapter, index))
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
			return "jpeg", false, index, nil
		} else {
			resp, _ = http.Get(fmt.Sprintf("%s/%s/%s/0%d.jpeg", static, nameFormated, chapter, index))
			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
				return "jpeg", true, index, nil
			}
			// Sometimes is possible to happen to acess a url loaded and without content (I really don't know why it happens)
			// fmt.Println("Error fetching chapter\nURL without content (-.-') . Adding chapter", chapter, "...")
			return "", false, 0, ErrFirstPageNotFound
		}
	}
	ext, hasZero, i, err := findIndex(1)
	if err != nil {
		ext, hasZero, i, err = findIndex(2)
		if err == ErrFirstPageNotFound {
			fmt.Println("Error fetching chapter\nURL without content (-.-') . Adding chapter", chapter, "...")
			err = ErrWebsiteWithoutContentData
			return
		}
	}
	/*
		On latest manga upates an url is given like that: "01.jpg", "07.jpg" instead of "1.jpg", "7.jpg"
		This function will add zero when it's needed
	*/
	format := func(count int, putZero bool) string {
		if putZero && count < 10 {
			return fmt.Sprintf("%s/%s/%s/0%d.%s", static, nameFormated, chapter, count, ext)
		}
		return fmt.Sprintf("%s/%s/%s/%d.%s", static, nameFormated, chapter, count, ext)
	}

	req := format(i, false)
	resp, err := http.Get(req)
	if resp.StatusCode == http.StatusNotFound {
		hasZero = true
		req = format(i, hasZero)
		resp, err = http.Get(req)
	}

	log.Println("\n", chapter, "| Req:", resp.Request.URL, "\nStatus code:", resp.StatusCode)
	if err != nil {
		return
	}
	initialPage := i
	// Find all manga pages
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
		// If ok, visits the remain of pages
		defer resp.Body.Close()
		for {
			i += 10
			req = format(i, hasZero)
			resp, err = http.Get(req)
			if resp.StatusCode != http.StatusOK {
				i -= 9
				req = format(i, hasZero)
				resp, err = http.Get(req)
				for resp.StatusCode == http.StatusOK {
					i++
					req = format(i, hasZero)
					resp, err = http.Get(req)
				}
				break
			}
		}
		for j := initialPage; j < i; j++ {
			ch.Pages = append(ch.Pages, format(j, hasZero))
		}
		ch.TotalPages = len(ch.Pages)
		return
	} else {
		err = ErrMangaNotFound
		return
	}
}
