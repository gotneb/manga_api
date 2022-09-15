package main

import (
	"github.com/gotneb/manga_api/api"
)

func main() {
	api.Init()
	/*
		m := &server.MuitoManga{}
		m.FindManga("jojo")
	*/
	/*
		m, err := server.GetClient(server.MUITO_MANGA).GetMangaPages("akame ga kill!", "10")
		if err != nil {
			panic(err)
		}
		log.Println(m.Pages)
	*/
}
