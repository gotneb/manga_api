package api

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/web"
)

func Init() {
	r := gin.Default()

	// A greeting
	r.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")
		if len(name) == 0 {
			c.String(http.StatusOK, "Hey there!")
		} else {
			c.String(http.StatusOK, "How you doing, "+name+"?")
		}
	})

	r.GET("/search/:mangaName", func(c *gin.Context) {
		name := c.Param("mangaName")
		manga, err := db.SeachManga(name)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, manga)
		}
	})

	r.GET("/pages/:mangaName/:chapter", func(c *gin.Context) {
		name := c.Param("mangaName")
		chapter, err := strconv.Atoi(c.Param("chapter"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			pages, err := web.FetchImagesByName(name, chapter)
			if err != nil {
				c.String(http.StatusNotFound, err.Error())
			} else {
				c.JSON(http.StatusOK, pages)
			}

		}
	})

	// Heroku
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	r.Run(":" + port)
}
