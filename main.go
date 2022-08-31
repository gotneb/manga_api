package main

import (
	"fmt"
	"strings"

	"github.com/gotneb/manga_api/web"
)

func main() {
	// Dummy data read, it'll only be useful for very small tests
	var nome string
	fmt.Print("Digite o nome do manga: ")
	fmt.Scanln(&nome)
	nome = format(nome)
	// ---------------------------------------------------------
	links := web.Search(nome)
	mangas := web.SetupCollect(&links)
	fmt.Printf("%s\nMANGAS\n%s\n", web.UselessLine(), web.UselessLine())
	for _, m := range mangas {
		fmt.Println(m.Title)
	}

	var n int
	fmt.Printf("Numero do capitulo: ")
	fmt.Scanf("%d", &n)
	web.FetchImagesByName(mangas[0].Title, n)
}

// Dummy functions only to execute (also dymmy) tests
func format(input string) (out string) {
	out = strings.ReplaceAll(input, "_", "+")
	return
}
