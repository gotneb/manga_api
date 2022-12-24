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
		case db.MANGAS_CHAN, db.SEEMANGAS:
			manga, err = server.Client(serv).GetManga(name)
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

		infoCh, err := server.Client(serv).GetMangaPages(name, ch)
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

		switch serv {
		case db.SEEMANGAS, db.MANGAS_CHAN:
			name := c.Param("mangaName")
			listMangas, err := db.SearchManga(serv, name)
			if err != nil || len(listMangas) == 0 {
				c.String(http.StatusNotFound, err.Error())
			} else {
				c.JSON(http.StatusOK, listMangas)
			}
		default:
			c.String(http.StatusForbidden, db.ErrUnauthorized.Error())
		}
	})

	/*
	 * BELOW THERE ARE ONLY ENDPOINTS THAT IS GOING TO BE USED BY ME
	 * SO, DON'T CARE ABOUT THEM ;)
	 */

	// This endpoint will upload all mangas avaliable on specified `server`
	r.GET("/:server/backup/:auth", func(c *gin.Context) {
		// Authenthication
		auth := c.Param("auth")
		if auth != db.AuthUpload {
			c.String(http.StatusBadRequest, "Password key to upload is wrong.")
			return
		}
		// Get server
		serv, err := strconv.Atoi(c.Param("server"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		// Upload
		utils.UploadAllMangasFrom(serv)
	})

	/*
	 * TLDR; Don't worry about this, it's an endpoint where only I can use :).
	 *
	 * I know it's pretty bad passing auth's key into url... I'll fix it later.
	 *
	 * The majority of sites they use AJAX to load content, with this aproach, I'm not able to
	 * scrape data from them. In this case I can use another tool (Seleniun) to get data and send it for here.
	 * e.g: READM.ORG
	 */
	r.POST("/:server/add/recent-mangas/:auth", func(c *gin.Context) {
		// Check credentials
		auth := c.Param("auth")
		if auth != db.AuthUpload {
			c.String(http.StatusForbidden, db.ErrUnauthorized.Error())
			return
		}
		// Get server
		serv, err := strconv.Atoi(c.Param("server"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		// Bind link to manga
		var data Releases
		err = c.BindJSON(&data)
		if err != nil {
			c.String(http.StatusBadGateway, err.Error())
			return
		}
		// Upload to specified server
		switch serv {
		case db.MANGAINN:
			utils.UploadRecentMangasFrom(db.MANGAINN, data.Links)
			c.String(http.StatusOK, "Saved with sucess")
		default:
			c.String(http.StatusBadRequest, "Server not found")
		}
	})

	// Listen to port
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	r.Run(":" + port)
}
