package server

import "github.com/gotneb/manga_api/web"

// It represents a site where manga data will be scraped
const (
	MEUS_MANGAS = iota
	MUITO_MANGA
)

// The "GetMangaPages" method use these strings to request from server
var pathImages = map[int]string{
	MEUS_MANGAS: "https://img.meusmangas.net//image",
	MUITO_MANGA: "https://cdn.statically.io/img/imgs.muitomanga.com/f=auto/imgs",
}

var (
	meusMangas = &MeusMangas{}
	muitoManga = &MuitoManga{}
)

// A Host is any site where it's possible scrape data from it
type Host interface {
	GetMangaDetail(mangaURL string) (web.Manga, error)
	GetMangaPages(mangaTitle, chapter string) (ch web.Chapter, err error)
}

/*
Returns the "site representation" where users will be able to request data
*/
func GetClient(server int) Host {
	switch server {
	case MEUS_MANGAS:
		return meusMangas
	case MUITO_MANGA:
		return muitoManga
	default:
		panic("this server doesn't exist")
	}
}
