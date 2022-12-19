package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"unicode"
)

var ErrMangaNotFound = errors.New("manga has't been found")
var ErrDataNotCollected = errors.New("manga data has't been collected")
var ErrPageNotFound = errors.New("manga page not found")
var ErrWebsiteWithoutContentData = errors.New("manga page null")
var ErrFirstPageNotFound = errors.New("first page not found")

// Returns all manga pages from a chapter
func FetchImagesByName(domain, name, chapter string) (ch Chapter, err error) {
	ch.Title = name
	ch.Value = chapter

	nameFormated := FormatedTitle(name)

	initialLetter := string(unicode.ToLower(rune(nameFormated[0])))
	request := fmt.Sprintf("%s/%v/%s/capitulo-%s/1.jpg", domain, initialLetter, nameFormated, chapter)
	resp, _ := http.Get(request)

	log.Println("\n", chapter, "| Req:", resp.Request.URL, "\nStatus code:", resp.StatusCode)
	if err != nil {
		return
	}

	i := 1
	initialPage := i
	// Find all manga pages
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
		// If ok, visits the remain of pages
		defer resp.Body.Close()
		for {
			i += 10
			req := fmt.Sprintf("%s/%v/%s/capitulo-%s/%d.jpg", domain, initialLetter, nameFormated, chapter, i)
			resp, err = http.Get(req)
			if resp.StatusCode != http.StatusOK {
				i -= 9
				req = fmt.Sprintf("%s/%v/%s/capitulo-%s/%d.jpg", domain, initialLetter, nameFormated, chapter, i)
				resp, err = http.Get(req)
				for resp.StatusCode == http.StatusOK {
					i++
					req = fmt.Sprintf("%s/%v/%s/capitulo-%s/%d.jpg", domain, initialLetter, nameFormated, chapter, i)
					resp, err = http.Get(req)
				}
				break
			}
		}
		for j := initialPage; j < i; j++ {
			ch.Pages = append(ch.Pages, fmt.Sprintf("%s/%v/%s/capitulo-%s/%d.jpg", domain, initialLetter, nameFormated, chapter, j))
		}
		ch.TotalPages = len(ch.Pages)
		return
	} else {
		err = ErrMangaNotFound
		return
	}
}
