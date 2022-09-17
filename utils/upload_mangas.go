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
	case serv.MEUS_MANGAS:
		uploadFromMeusMangas()
	case serv.MANGAINN:
		uploadFromMangainn()
	default:
		panic(errors.New("server not found"))
	}
}

func uploadFromMeusMangas() {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	client := serv.GetClient(serv.MEUS_MANGAS)

	for _, letter := range alphabet {
		links := client.FetchAllMangaByLetter(string(letter))
		for _, link := range links {
			manga, stt := client.GetMangaDetail(link)
			if stt != http.StatusOK {
				log.Fatalln("Status:", stt)
			}
			db.AddManga(serv.MEUS_MANGAS, &manga)
		}
	}
}

func uploadFromMangainn() {

}
