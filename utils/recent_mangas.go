package utils

import (
	"errors"
	"log"
	"net/http"

	"github.com/gotneb/manga_api/db"
	serv "github.com/gotneb/manga_api/server"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
UploadRecentMangasFrom goes to specified server and get all recent mangas updated
and stores it into the database.
*/
func UploadRecentMangasFrom(server int, links []string) {
	switch server {
	case db.SEEMANGAS, db.MANGAINN:
		addLatests(server, links)
	default:
		panic(errors.New("server not found"))
	}
}

func addLatests(serverCode int, latestLinks []string) {
	stop := false

	client := serv.Client(serverCode)

	for !stop {
		links := latestLinks
		for _, link := range links {
			manga, stt := client.GetMangaDetail(link)
			if stt != http.StatusOK {
				log.Fatalln("manga not found")
			}
			// If do not exist, add it
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
