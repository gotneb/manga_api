package api

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/server"
	"github.com/gotneb/manga_api/web"
)

func Init() {
	r := gin.Default()

	// An useless greeting :)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to my api :] !")
	})

	// Some people can have trouble forgeting to include "server" option
	r.GET("/manga/detail/:mangaName", func(c *gin.Context) {
		c.String(http.StatusBadRequest, "Request should include server option")
	})

	r.GET("/:server/manga/detail/:mangaName", func(c *gin.Context) {
		serv, err := strconv.Atoi(c.Param("server"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		name := c.Param("mangaName")
		var manga web.Manga

		if serv == server.MEUS_MANGAS {
			manga, err = db.GetManga(name)
		} else if serv == server.MUITO_MANGA {
			x := &server.MuitoManga{}
			manga, err = x.FindManga(name)
		}

		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, manga)
		}
	})

	r.GET("/:server/manga/pages/:mangaName/:chapter", func(c *gin.Context) {
		serv, err := strconv.Atoi(c.Param("server"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		name := c.Param("mangaName")
		ch := c.Param("chapter")

		//infoCh, err := web.FetchImagesByName(name, ch)
		infoCh, err := server.GetClient(serv).GetMangaPages(name, ch)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}
		c.JSON(http.StatusOK, infoCh)
	})

	r.GET("/manga/search/:mangaName", func(c *gin.Context) {
		name := c.Param("mangaName")
		listMangas, err := db.SearchManga(name)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, listMangas)
		}
	})

	// Heroku
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	r.Run(":" + port)
}
