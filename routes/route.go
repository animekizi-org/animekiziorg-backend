package routes

import (
	"fmt"
	"log"
	"main/util"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

type DownloadBody struct {
	Url string `json:"videoUri"`
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}

func DownloadVideo(c echo.Context) error {
	var downloadBody DownloadBody

	if err := c.Bind(&downloadBody); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": false,
			"path":   "not_found.html",
		})
	}

	err, path := util.DownloadRedditVideo(downloadBody.Url, c.RealIP())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": false,
			"path":   "not_found.html",
			"err":    err.Error(),
		})
	}

	err, post := util.GetPost(path)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": false,
			"path":   "not_found.html",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  true,
		"message": "Video downloaded successfully.",
		"level":   "success",
		"path":    path,
		"post": echo.Map{
			"title":     post.PostTitle,
			"date":      post.Date,
			"url":       post.PostURL,
			"thumbnail": post.Thumbnail,
		},
	})
}

func GetVideo(c echo.Context) error {
	videoId := c.Param("videoId")
	if videoId == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  false,
			"message": "You need to use 'videoId' param.",
			"level":   "error",
		})
	}

	_, err := os.Stat(fmt.Sprintf("./tmp/%s.mp4", videoId))
	if err != nil || videoId == "undefined" {
		return c.JSON(http.StatusOK, echo.Map{
			"status":  false,
			"message": "Video not found.",
			"level":   "error",
		})
	}

	return c.Attachment(fmt.Sprintf("./tmp/%s.mp4", videoId), fmt.Sprintf("%s.mp4", videoId))
}

func DeleteEntry(c echo.Context) error {
	auth := c.QueryParam("auth")
	videoId := c.QueryParam("videoId")

	if videoId == "" || auth == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  false,
			"message": "You need to use 'videoId' and 'auth' query param..",
			"level":   "error",
		})
	}

	if auth != os.Getenv("PROD_AUTH") {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  false,
			"message": "Wrong auth key.",
			"level":   "error",
		})
	}

	err := util.DeletePost(videoId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  false,
			"message": err.Error(),
			"level":   "error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  true,
		"message": "Successfully deleted the video.",
		"level":   "error",
	})
}

func RetrieveLatestVideos(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}

	err, list := util.RetrieveLatestVideos(page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  false,
			"message": err.Error(),
			"level":   "error",
		})
	}

	err, count := util.RetrieveCount()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  false,
			"message": err.Error(),
			"level":   "error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":          true,
		"videos":          list,
		"totalDownloaded": count,
	})
}
