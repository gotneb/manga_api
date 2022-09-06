package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gotneb/manga_api/api"
	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/web"
)

func main() {
	api.Init(UploadMangaPages)
	/*
		name := "https://meusmangas.net/manga/hunter-x-hunter/"
		manga, _ := web.FetchMangaData(name)
		ch, _ := web.FetchImagesByName(manga.Title, "390")
		fmt.Printf("Total pages: %d\n", ch.TotalPages)
	*/

	/*
		name := "https://meusmangas.net/manga/mango/the-landlords-little-girl"
		manga, _ := web.FetchMangaData(name)
		for _, v := range manga.Chapters {
			ch, _ := web.FetchImagesByName(manga.Title, v)
			fmt.Printf("Total pages: %d\n\n", ch.TotalPages)
			//db.AddChapter(&ch)
		}
	*/
}

func UploadMangaPages() {
	section := "abcdefghijklmnopqrstuvwxyz"
	for _, v := range section {
		links := web.FetchPages(string(v))

		for _, link := range links {
			manga, err := web.FetchMangaData(link)
			if err != nil {
				panic(err)
			}
			fmt.Println("\nStarting upload manga pages from", manga.Title)
			for _, v := range manga.Chapters {
				ch, _ := web.FetchImagesByName(manga.Title, v)
				fmt.Printf("Total pages: %d\n\n", ch.TotalPages)
				db.AddChapter(&ch)
			}
		}
	}
	log.Println("DONE!!!!")
}

// Upload all manga avaliable on website
func uploadMangaDetails() {
	// If hosted on heroku, automatically run gin framework
	port := os.Getenv("PORT")
	if len(port) == 0 {
		var option string
		fmt.Printf("Start web scraping? [S/N]: ")
		fmt.Scanf("%s", &option)

		if strings.ToLower(string(option[0])) == "s" {
			section := "abcdefghijklmnopqrstuvwxyz"
			for _, v := range section {
				links := web.FetchPages(string(v))

				for _, link := range links {
					manga, err := web.FetchMangaData(link)
					if err != nil {
						panic(err)
					}
					db.AddManga(&manga)
				}
			}
			log.Println("DONE!!!!")
		} else {
			fmt.Println("Runing API!")
			api.Init(nil)
		}
	} else {
		api.Init(nil)
	}
}
