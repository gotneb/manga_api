package db

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/server"
	"go.mongodb.org/mongo-driver/mongo"
)

const uri = "https://meusmangas.net"

// Get all manga links updates on the website
func GetRecentUpdates(page int) (links []string) {
	c := colly.NewCollector()
	link := fmt.Sprintf("https://meusmangas.net/comienzo/page/%d", page)

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting:", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	c.OnHTML("div#recent_releases  li.item_news-manga a.pull-left", func(e *colly.HTMLElement) {
		link := uri + e.Attr("href")
		links = append(links, link)
		log.Println(link)
	})

	c.Visit(link)

	return
}

func AddAllRecentMangas() {
	stop := false
	i := 1
	for !stop {
		links := GetRecentUpdates(i)
		for _, link := range links {
			manga, err := server.GetClient(server.MEUS_MANGAS).GetMangaDetail(link)
			if err != nil {
				panic(err)
			}
			// If do not exists, add it
			_, err = FindManga(manga.Title)
			if err == mongo.ErrNoDocuments {
				AddManga(&manga)
			}
			res, err := GetManga(manga.Title)
			if err != nil {
				panic(err)
			}
			if manga.Title == res.Title {
				if manga.TotalChapters == res.TotalChapters {
					stop = true
					break
				} else {
					log.Printf("Updating: [%s]", manga.Title)
					UpdateManga(&manga)
				}
			}
		}
		i++
	}
	log.Println("All mangas are up to date!")
}
