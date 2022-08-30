package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Manga struct {
	title       string
	author      string
	tags        []string
	chapters    int
	description string
	thumbnail   string
	situation   string
}

func (m *Manga) Show() {
	fmt.Println("Title: " + m.title)
	fmt.Printf("Author: %s\n", m.author)
	fmt.Printf("Situation: %s\n", m.situation)
	fmt.Printf("Chapters: %d\n", m.chapters)
	fmt.Println("Thumbnail: " + m.thumbnail)
	fmt.Printf("Tags: ")
	for _, v := range m.tags {
		fmt.Printf("%s, ", v)
	}
	fmt.Println()
	fmt.Println("Description: " + m.description)
}

func main() {
	manga := &Manga{}
	c := colly.NewCollector()

	var nome string
	fmt.Print("Digite o nome do manga: ")
	fmt.Scanln(&nome)
	nome = Format(nome)
	//nome = "https://www.google.com/search?q=goblin+slayer+meus+mangas"
	links := Search(nome)

	// Fetch manga sinopse
	c.OnHTML("div.sinopse-page", func(e *colly.HTMLElement) {
		manga.description = e.Text
	})
	// Fetch description manga
	c.OnHTML("img.hGq41", func(e *colly.HTMLElement) {
		manga.thumbnail = e.Attr("src")
	})
	// Fetch title manga
	c.OnHTML("h1.kw-title", func(e *colly.HTMLElement) {
		manga.title = e.Text
	})
	// Fetch tags manga
	c.OnHTML("a.widget-btn", func(e *colly.HTMLElement) {
		manga.tags = append(manga.tags, e.Text)
	})
	// Fetch chapters
	c.OnHTML("div.jVBw-infos > span", func(e *colly.HTMLElement) {
		number := strings.Split(e.Text, " ")[0]
		if len(number) > 0 {
			chapters, err := strconv.Atoi(number)
			if err == nil {
				manga.chapters = chapters
			}
		}
	})
	// Fetch manga situation
	c.OnHTML("div.jVBw-infos span.mdq", func(e *colly.HTMLElement) {
		manga.situation = e.Text
	})
	// Fetch manga author
	c.OnHTML("div.jVBw-infos div", func(e *colly.HTMLElement) {
		manga.author = e.Text
	})
	// Finished
	c.OnScraped(func(r *colly.Response) {
		manga.Show()
	})
	// Entering on site
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	for _, link := range links {
		c.Visit(link)
	}
}

func Format(input string) (out string) {
	out = strings.ReplaceAll(input, "_", "+")
	return
}

func Search(mangaName string) (links []string) {
	c := colly.NewCollector()
	lookFor := mangaName + "+meus+mangas"
	staticLink := "https://www.google.com/search?q=" + lookFor

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		def := "https://meusmangas.net/manga/mango/"
		att := e.Attr("href")
		url := ""
		if strings.Contains(att, def) {
			url = att
		}
		if url != "" {
			url = url[7:strings.Index(url, "&")]
			links = append(links, url)
			fmt.Println(url)
		}
	})
	c.Visit(staticLink)
	return
}
