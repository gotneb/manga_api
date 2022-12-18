# Mangahoot :coffee:
A manga api written with Go using MongoDB :heart: 

##  About  
This project uses a web scrapping tool to get all manga data.

:warning: You may face some problems using it. It's not done yet.

**API PATH** = [https://mangahoot.onrender.com](https://mangahoot.onrender.com) 

## Servers

| Server  |  Host  | Language |
| --- | --- | --- |
|  0 |  https://seemangas.com/ | :brazil: |

## Manga Detail
Returns *one single result* which matches with the name
```
/[server]/manga/detail/[manga name]
```
example: https://mangahoot.onrender.com/0/manga/detail/berserk

## Chapter Pages
```
/[server]/manga/pages/[name of manga]/[chapter number]
```
example: https://mangahoot.onrender.com/0/manga/pages/martial%20peak/500

## Search Manga
Returns *many results* which matches with the name
```
/[server]/manga/search/[manga name]
```
example: https://mangahoot.onrender.com/0/manga/search/jojo
