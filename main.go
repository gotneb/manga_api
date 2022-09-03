package main

import (
	"github.com/gotneb/manga_api/api"
)

func main() {
	api.Init()
	/*
		// Dummy data read, it'll only be useful for very small tests
		var nome string
		fmt.Print("Digite o nome do manga: ")
		fmt.Scanln(&nome)
		nome = format(nome)
		// Search a manga
		links, err := web.Search(nome)
		if err != nil {
			panic(err)
		}
		// Choose between them
		opt := 0
		if len(links) > 1 {
			fmt.Println("AValiable mangas:")
			for i, link := range links {
				fmt.Printf("[%d]: %s\n", i+1, link)
			}
			fmt.Printf("Option: ")
			fmt.Scanf("%d", &opt)
		}
		manga, err := web.FetchMangaData(links[opt-1])
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\nMANGAS\n%s\n", web.UselessLine(), web.UselessLine())
		fmt.Println(manga.Title)

		var n int
		fmt.Printf("Numero do capitulo: ")
		fmt.Scanf("%d", &n)
		fmt.Println("Looking pages...")
		pages, err := web.FetchImagesByName(manga.Title, n)
		if err != nil {
			panic(err)
		}
		for _, p := range pages {
			fmt.Println(p)
		}
	*/
}
