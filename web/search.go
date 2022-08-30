package web

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func SetupCollect(c *colly.Collector, mangas []*Manga, links *[]string) {
	manga := &Manga{}
	// Entering on site
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	// Fetch manga sinopse
	c.OnHTML("div.sinopse-page", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})
	// Fetch description manga
	c.OnHTML("img.hGq41", func(e *colly.HTMLElement) {
		manga.Thumbnail = e.Attr("src")
	})
	// Fetch title manga
	c.OnHTML("h1.kw-title", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})
	// Fetch tags manga
	c.OnHTML("a.widget-btn", func(e *colly.HTMLElement) {
		manga.Tags = append(manga.Tags, e.Text)
	})
	// Fetch chapters
	c.OnHTML("div.jVBw-infos > span", func(e *colly.HTMLElement) {
		number := strings.Split(e.Text, " ")[0]
		if len(number) > 0 {
			chapters, err := strconv.Atoi(number)
			if err == nil {
				manga.Chapters = chapters
			}
		}
	})
	// Fetch manga situation
	c.OnHTML("div.jVBw-infos span.mdq", func(e *colly.HTMLElement) {
		manga.Situation = e.Text
	})
	// Fetch manga author
	c.OnHTML("div.jVBw-infos div", func(e *colly.HTMLElement) {
		manga.Author = e.Text
	})
	// Finished
	c.OnScraped(func(r *colly.Response) {
		manga.Show()
		mangas = append(mangas, manga)
		manga = new(Manga)
	})

	for _, link := range *links {
		c.Visit(link)
	}
}

// Searches a manga on google and returns a slice of matching results
func Search(mangaName string) (links []string) {
	c := colly.NewCollector()
	lookFor := mangaName + "+meus+mangas"
	staticLink := "https://www.google.com/search?q=" + lookFor

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		defAdresses := [2]string{
			"https://meusmangas.net/manga/mango/",
			"https://meusmangas.net/manga/hd/",
		}
		url := e.Attr("href")
		//fmt.Printf("All <a> founds: %s\n", url)
		if strings.Contains(url, defAdresses[0]) {
			url = url[7:strings.Index(url, "&")]
			links = append(links, url)
			//fmt.Println("Found: " + url)
		} else if strings.Contains(url, defAdresses[1]) {
			url = url[7:strings.Index(url, "&")]
			links = append(links, url)
		}
	})
	c.Visit(staticLink)
	return
}
