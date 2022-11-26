package api

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/server"
	"github.com/gotneb/manga_api/utils"
	"github.com/gotneb/manga_api/web"
)

type Releases struct {
	Links []string `json:"links" validate:"dive"`
}

func Init() {
	r := gin.Default()

	// An useless greeting :)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to my api :] !")
	})

	r.GET("/:server/manga/detail/:mangaName", func(c *gin.Context) {
		serv, err := strconv.Atoi(c.Param("server"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		name := c.Param("mangaName")
		var manga web.Manga

		switch serv {
		case db.MEUS_MANGAS:
			manga, err = server.GetClient(db.MEUS_MANGAS).GetManga(name)
		case db.MANGAINN:
			manga, err = server.GetClient(db.MANGAINN).GetManga(name)
		default:
			err = errors.New("server not found")
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

		infoCh, err := server.GetClient(serv).GetMangaPages(name, ch)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		}
		c.JSON(http.StatusOK, infoCh)
	})

	r.GET("/:server/manga/search/:mangaName", func(c *gin.Context) {
		serv, err := strconv.Atoi(c.Param("server"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		name := c.Param("mangaName")
		listMangas, err := db.SearchManga(serv, name)
		if err != nil || len(listMangas) == 0 {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, listMangas)
		}
	})

	/*
		r.POST("/add/manga/:password", func(c *gin.Context) {
			pass := c.Param("password")
			if pass != db.PasswordUpload {
				c.String(http.StatusBadRequest, "Password key to upload is wrong.")
				return
			}

			var manga web.Manga
			err := c.BindJSON(&manga)
			if err != nil {
				c.String(http.StatusBadGateway, err.Error())
			}
			db.AddManga(2, &manga)
			c.String(http.StatusOK, "Saved with sucess")
		})
	*/

	/*
	 * TLDR; Don't worry about this, it's an endpoint where only I can use :).
	 *
	 * The majority of sites they use AJAX to load content, with this aproach, I'm not able to
	 * scrape data from them. In this case I can use another tool (Seleniun) to get data and send it for here.
	 *
	 * [WARNING]: You might expect try this endpoint... But "utils" need authenthication, so if you try, then you fail.
	 */
	r.POST("/add/release-mangas/", func(c *gin.Context) {
		var data Releases
		err := c.BindJSON(&data)

		if err != nil {
			c.String(http.StatusBadGateway, err.Error())
			return
		}

		utils.UploadRecentMangasFrom(db.MEUS_MANGAS, data.Links)
		c.String(http.StatusOK, "Saved with sucess")
	})

	// Listen to port
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	r.Run(":" + port)
}
