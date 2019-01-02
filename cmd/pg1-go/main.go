package main

import (
	"html/template"
	"net/http"
	"os"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"

	humanize "github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igprofile"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/jobqueue"
)

var mainLogger = logger.NewLogger("PG1-Go::Main", false, true)

func dec(num int) int {
	return num - 1
}

func inc(num int) int {
	return num + 1
}

func humInt(num int) string {
	return humanize.Comma(int64(num))
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		mainLogger.Error("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.SetFuncMap(template.FuncMap{
		"decrease": dec,
		"increase": inc,
		"humInt":   humInt,
	})
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	jobqueue.DefineAPIRoutes(router, "")
	jobqueue.DefineViewRoutes(router, "")
	igprofile.DefineAPIRoutes(router, "api")
	igprofile.DefineViewRoutes(router, "")
	router.Run(":" + port)
}
