package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var ErrMangaNotFound = errors.New("manga has't been found")
var ErrDataNotCollected = errors.New("manga data has't been collected")
var ErrPageNotFound = errors.New("manga page not found")
var ErrWebsiteWithoutContentData = errors.New("manga page null")
var ErrFirstPageNotFound = errors.New("first page not found")

// Returns all manga pages from a chapter
func FetchImagesByName(hostImages, name, chapter string) (ch Chapter, err error) {
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
		resp, _ := http.Get(fmt.Sprintf("%s/%s/%s/%d.jpg", hostImages, nameFormated, chapter, index))
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
			return "jpg", false, index, nil
		}
		resp, _ = http.Get(fmt.Sprintf("%s/%s/%s/0%d.jpg", hostImages, nameFormated, chapter, index))
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
			return "jpg", true, index, nil
		}
		resp, _ = http.Get(fmt.Sprintf("%s/%s/%s/%d.jpeg", hostImages, nameFormated, chapter, index))
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent {
			return "jpeg", false, index, nil
		} else {
			resp, _ = http.Get(fmt.Sprintf("%s/%s/%s/0%d.jpeg", hostImages, nameFormated, chapter, index))
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
			return fmt.Sprintf("%s/%s/%s/0%d.%s", hostImages, nameFormated, chapter, count, ext)
		}
		return fmt.Sprintf("%s/%s/%s/%d.%s", hostImages, nameFormated, chapter, count, ext)
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
