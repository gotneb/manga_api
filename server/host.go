package server

import "github.com/gotneb/manga_api/web"

// A Host is any site where it's possible scrape data from it
type Host interface {
	GetMangaDetail(mangaURL string) (web.Manga, error)
	GetMangaPages(mangaTitle string, chapter int) (ch web.Chapter, err error)
}
