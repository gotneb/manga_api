# Mangahoot :coffee:
A manga api written with Go using MongoDB :heart: 

##  About  
This project uses a web scrapping tool to get all manga data.

**API PATH** = [https://mangahoot.onrender.com](https://mangahoot.onrender.com) 

## Servers

| Server  |  Host  | Language |
| --- | --- | --- |
| 1 | https://www.mangainn.net | :us_outlying_islands: |

## Manga Detail
Returns *one single result* which matches with the name
```
/[server]/manga/detail/[manga name]
```
example: https://mangahoot.onrender.com/1/manga/detail/berserk

## Manga Pages
Returns all pages related to *manga name* on specified *chapter*

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
