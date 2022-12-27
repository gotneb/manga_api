# Mangahoot :coffee:
A manga api written with Go using MongoDB :heart: 

##  About  
This project uses a web scrapping tool to get all manga data.

:warning: You may face some problems using it. It's not done yet.

**API PATH** = [https://mangahoot.onrender.com](https://mangahoot.onrender.com) 

## Servers

| Server  |  Host  | Language |
| --- | --- | --- |
| 0 | https://seemangas.com | :brazil: |
| 1 | https://www.mangainn.net/ | :us_outlying_islands: |
| 3 | https://mangaschan.com | :brazil: |

## Manga Detail
Returns *one single result* which matches with the name
```
/[server]/manga/detail/[manga name]
```
example: https://mangahoot.onrender.com/1/manga/detail/berserk

## Manga Pages
Returns all pages related to *manga name* on specified *chapter*

:warning: **Only works to Mangainn**
```
/[server]/manga/pages/[manga name]/[chapter]
```
example: https://mangahoot.onrender.com/1/manga/pages/akame-ga-kill/30

## Search Manga
Returns *many results* which matches with the name
```
/[server]/manga/search/[manga name]
```
example: https://mangahoot.onrender.com/1/manga/search/jojo

## TODO
- [ ] Fetch manga pages 
- [ ] Fetch manga by genre
- [ ] Get populars mangas

# Golang Side
If you're interested about how to use it on Go, this is how:

```

import (
	"log"
	"net/http"

	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/server"
)

func getMangaDetail(url string) {
	/*
	 * First you will have to identify which server you are going to get data from
	 * In 2022-12-27, there are 3 clients avaliable:
	 * Mangainn(en), MangasChan(pt-br), Seemangas(pt-br)
	 *
	 * Let's choose Mangainn:
	 */
	client := server.Client(db.MANGAINN)

	/*
	 * The 'url' param must be the adress of somewhere where scrapper can get:
	 * title, author, status, description etc...
	 * Like that:
	 * - https://www.mangainn.net/chainsaw-man
	 * - https://seemangas.com/manga/one-piece-my9771
	 * - https://mangaschan.com/manga/god-of-martial-arts/
	 */
	manga, status := client.GetMangaDetail(url)

	/*
	 * Finally when you go get details from a manga, the above method will give you the
	 * data requested and the status (a http status)
	 * Sometimes can occur that web-scrapping tool failing to get data
	 */
	if status != http.StatusOK {
		log.Fatalln("An unexpected error has ocurred on getting data...")
	}
	manga.Show()
}

```
