package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/web"
)

func Init() {
	r := gin.Default()

	// A greeting
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to my api :] !")
	})

	r.GET("/manga/detail/:mangaName", func(c *gin.Context) {
		name := c.Param("mangaName")
		manga, err := db.SeachManga(name)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, manga)
		}
	})

	r.GET("/manga/pages/:mangaName/:chapter", func(c *gin.Context) {
		name := c.Param("mangaName")
		ch := c.Param("chapter")
		pages, err := web.FetchImagesByName(name, ch)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, pages)
		}
	})

	// Heroku
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	r.Run(":" + port)
}
