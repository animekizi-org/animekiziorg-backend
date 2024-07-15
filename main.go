package main

import (
	"fmt"
	"log"
	"main/prisma/db"
	"main/routes"
	"main/util"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func CheckVideos() {

	for {
		posts, err := util.ReturnDatabase().Post.FindMany().Exec(*util.ReturnContext())
		if err != nil {
			continue
		}

		for _, post := range posts {
			path := fmt.Sprintf("./tmp/%s.mp4", post.ID)

			if post.Date.Add(time.Hour*24*5).Unix() < time.Now().Unix() {
				if _, err := util.ReturnDatabase().Post.FindUnique(db.Post.ID.Equals(post.ID)).Delete().Exec(*util.ReturnContext()); err != nil {
					log.Printf("[ERROR] Error while deleting post: %s", err.Error())
				} else {
					log.Printf("[INFO] Post deleted: \"%s\"", post.ID)
					if err := os.Remove(path); err != nil {
						log.Printf("[ERROR] while deleting video: %s", err.Error())
					} else {
						log.Printf("[INFO] Video deleted: \"%s\"", post.ID)
					}
				}

			}
		}

		time.Sleep(time.Minute * 60)
	}
}

func main() {
	go CheckVideos() // run an asynchronous function that deleted videos automatically.

	app := echo.New()

	app.Use(middleware.Recover()) // recovers from errors
	app.Use(middleware.CORS())    // cors

	app.POST("/api/download", routes.DownloadVideo)             // POST /api/download
	app.GET("/api/videos/:videoId", routes.GetVideo)            //POST /api/videos/:videoId
	app.GET("/api/retrieveLatest", routes.RetrieveLatestVideos) //POST /api/retrieveLatest
	app.GET("/api/deleteVideo", routes.DeleteEntry)             //POST /api/deleteVideo

	app.Logger.Fatal(app.Start(":1337"))
}
