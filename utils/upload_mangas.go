package utils

import (
	"errors"
	"log"
	"net/http"

	"github.com/gotneb/manga_api/db"
	serv "github.com/gotneb/manga_api/server"
)

// UploadAllMangasFrom upload all mangas from specified server and stores into database.
func UploadAllMangasFrom(server int) {
	switch server {
	case db.MEUS_MANGAS:
		uploadMangas(db.MEUS_MANGAS)
	case db.MANGAINN:
		uploadMangas(db.MANGAINN)
	case db.MANGAS_CHAN:
		uploadMangas(db.MANGAS_CHAN)
	default:
		panic(errors.New("server not found"))
	}
}

func uploadMangas(serverCode int) {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	client := serv.Client(serverCode)

	for _, letter := range alphabet {
		links := client.FetchAllMangaByLetter(string(letter))
		for _, link := range links {
			manga, stt := client.GetMangaDetail(link)
			if stt != http.StatusOK {
				log.Println("Error on Status:", stt)
			} else {
				/*
					Weirdly "Mangainn" adds mangas without even chapters, only title
					So, the above line guarantees that only mangas with chapters will be added
				*/
				if manga.TotalChapters > 2 && manga.Status != "Completed" {
					db.AddManga(serverCode, &manga)
				}
			}
		}
	}
}
