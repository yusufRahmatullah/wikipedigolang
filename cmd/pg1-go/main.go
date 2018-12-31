package main

import (
	"net/http"
	"os"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
)

var mainLogger = logger.NewLogger("PG1-Go::Main", false, true)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		mainLogger.Error("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	jobqueue.DefineAPIRoutes(router, "")
	jobqueue.DefineViewRoutes(router, "")
	igprofile.DefineAPIRoutes(router, "")

	router.Run(":" + port)
}
