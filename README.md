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
| 3 | https://mangaschan.com | :brazil: |

## Manga Detail
Returns *one single result* which matches with the name
```
/[server]/manga/detail/[manga name]
```
example: https://mangahoot.onrender.com/0/manga/detail/berserk

example: https://mangahoot.onrender.com/3/manga/detail/berserk

## Search Manga
Returns *many results* which matches with the name
```
/[server]/manga/search/[manga name]
```
example: https://mangahoot.onrender.com/0/manga/search/jojo

example: https://mangahoot.onrender.com/3/manga/search/jojo

## TODO
- [ ] Fetch manga pages 
- [ ] Fetch manga by genre
- [ ] Get populars mangas
