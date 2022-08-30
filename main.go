package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gotneb/manga_api/web"
)

func main() {
	var mangas []*web.Manga
	c := colly.NewCollector()

	// Dummy data read, it'll only be useful for very small tests
	var nome string
	fmt.Print("Digite o nome do manga: ")
	fmt.Scanln(&nome)
	nome = format(nome)
	// ---------------------------------------------------------
	links := web.Search(nome)

	web.SetupCollect(c, mangas, &links)
}

// Dummy functions only to execute (also dymmy) tests
func format(input string) (out string) {
	out = strings.ReplaceAll(input, "_", "+")
	return
}
