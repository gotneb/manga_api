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
	// If hospeded heroku, automatically run gin framework
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
			api.Init()
		}
	} else {
		api.Init()
	}
}
