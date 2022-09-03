package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotneb/manga_api/web"
)

func Init() {
	r := gin.Default()

	r.GET("/search/:mangaName", func(c *gin.Context) {
		mangaName := c.Param("mangaName")
		mangaName = strings.ReplaceAll(mangaName, " ", "+")
		links, err := web.Search(mangaName)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			manga, err := web.FetchMangaData(links[0])
			if err != nil {
				c.String(http.StatusNotFound, err.Error())
			} else {
				c.JSON(http.StatusOK, manga)
			}
		}
	})

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	r.Run()
}
