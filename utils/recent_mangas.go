package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/db"
	serv "github.com/gotneb/manga_api/server"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
UploadRecentMangasFrom goes to specified server and get all recent mangas updated
and stores it into the database.
*/
func UploadRecentMangasFrom(server int) {
	switch server {
	case db.MEUS_MANGAS:
		getFromMeusMangas()
	case db.MANGAINN:
		getFromMangainn()
	default:
		panic(errors.New("server not found"))
	}
}

// Get all recent mangas from "meusmangas.net"
func getRecentUpdates(page int) (links []string) {
	c := colly.NewCollector()
	link := fmt.Sprintf("https://meusmangas.net/comienzo/page/%d", page)

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting:", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	c.OnHTML("div#recent_releases  li.item_news-manga a.pull-left", func(e *colly.HTMLElement) {
		link := "https://meusmangas.net" + e.Attr("href")
		links = append(links, link)
		log.Println(link)
	})

	c.Visit(link)

	return
}

func getFromMeusMangas() {
	stop := false
	i := 1

	client := serv.GetClient(db.MEUS_MANGAS)
	serverCode := db.MEUS_MANGAS

	for !stop {
		links := getRecentUpdates(i)
		for _, link := range links {
			manga, stt := client.GetMangaDetail(link)
			if stt != http.StatusOK {
				log.Fatalln("manga not found")
			}
			// If do not exists, add it
			_, err := db.FindManga(serverCode, manga.Title)
			if err == mongo.ErrNoDocuments {
				db.AddManga(serverCode, &manga)
			}
			res, err := db.GetManga(serverCode, manga.Title)
			if err != nil {
				panic(err)
			}
			if manga.Title == res.Title {
				if manga.TotalChapters == res.TotalChapters {
					stop = true
					break
				} else {
					log.Printf("Updating: [%s]", manga.Title)
					db.UpdateManga(serverCode, &manga)
				}
			}
		}
	}
}

func getFromMangainn() {

}