package main

import (
	"github.com/gotneb/manga_api/api"
	"github.com/gotneb/manga_api/db"
)

func main() {
	api.Init()
	/*
		section := "rst"
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
	*/
	db.GetManga("berserk")
}
