package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/auth"
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
		mainLogger.Error("$PORT must be set", nil)
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		mainLogger.Warning("$SESSION_SECRET must be set, set as default")
		sessionSecret = "LongEnoughSecretKeyTooAvoidBruteForce"
	}

	store := sessions.NewCookieStore([]byte(sessionSecret))

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(sessions.Sessions("defaultSession", store))
	router.SetFuncMap(template.FuncMap{
		"decrease": dec,
		"increase": inc,
		"humInt":   humInt,
	})
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		// c.HTML(http.StatusOK, "index.tmpl.html", nil)
		c.Redirect(http.StatusTemporaryRedirect, "/igprofiles")
	})

	jobqueue.DefineAPIRoutes(router, "api")
	jobqueue.DefineViewRoutes(router, "")
	igprofile.DefineAPIRoutes(router, "api")
	igprofile.DefineViewRoutes(router, "")
	auth.DefineViewRoutes(router, "")
	auth.DefineAPIRoutes(router, "api")
	router.Run(":" + port)
}
