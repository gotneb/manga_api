package server

import (
	"github.com/gotneb/manga_api/web"
)

// It represents a site where manga data will be scraped
const (
	MEUS_MANGAS = iota
	MANGAINN
)

// The "GetMangaPages" method use these strings to request from server
var pathImages = map[int]string{
	MEUS_MANGAS: "https://img.meusmangas.net//image",
	MANGAINN:    "https://www.mangainn.net",
}

var (
	meusMangas = &MeusMangas{}
	mangainn   = &Mangainn{}
)

// A Host is any site where it's possible scrape data from it.
type Host interface {
	// GetMangaDetail returns a manga hosted on mangaURL.
	GetMangaDetail(mangaURL string) (manga web.Manga, statusCode int)
	// GetMangaPages returns all chapter pages found on that site.
	GetMangaPages(mangaTitle, chapter string) (ch web.Chapter, err error)
	// FetchAllMangaByLetter returns all manga links hosted on site at specified letter.
	FetchAllMangaByLetter(letter string) (links []string)
}

// GetClient returns the site representation where users will be able to request data
func GetClient(server int) Host {
	switch server {
	case MEUS_MANGAS:
		return meusMangas
	case MANGAINN:
		return mangainn
	default:
		panic("this server doesn't exist")
	}
}
