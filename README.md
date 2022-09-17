# Mangahoot :coffee:
A manga api written with Go using MongoDB :heart: 

##  About
This project uses a web scrapping tool to get all manga data.

**API PATH** = [https://mangahoot.herokuapp.com/](https:/mangahoot.herokuapp.com/) 

## Servers

| Server  |  Host  |  |
| --- | --- | --- |
|  0 |  https://meusmangas.net/ | :brazil: |
|  1 |  https://www.mangainn.net/ | :us: |

## Manga Detail
Returns *one single result* which matches with the name
```
/[server]/manga/detail/[manga name]
```
example: https://mangahoot.herokuapp.com/0/manga/detail/berserk

## Chapter Pages
```
/[server]/manga/pages/[name of manga]/[chapter number]
```
example: https://mangahoot.herokuapp.com/1/manga/pages/one%20punch%20man/70

## Search Manga
Returns *many results* which matches with the name
```
/manga/search/[manga name]
```
https://mangahoot.herokuapp.com/manga/search/jojo

## :sparkles: Contact
I'm looking for a job. Do you want hire-me? [My email](gabriel_origenstdb@gmail.com). 
