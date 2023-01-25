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
