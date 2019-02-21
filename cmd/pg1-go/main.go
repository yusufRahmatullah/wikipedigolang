package main

import (
	"html/template"
	"os"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igmedia"

	"github.com/gin-contrib/cors"
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

func ginMode() string {
	gm := os.Getenv("GIN_MODE")
	return gm
}

func noCacheHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Next()
	}
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
	router.Use(noCacheHeader())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://yusufRahmatullah.github.io", "https://yusufRahmatullah.github.io"},
		AllowMethods: []string{"GET"},
	}))

	router.SetFuncMap(template.FuncMap{
		"decrease": dec,
		"increase": inc,
		"humInt":   humInt,
		"ginMode":  ginMode,
	})
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	jobqueue.DefineViewRoutes(router, "")
	igprofile.DefineViewRoutes(router, "")
	auth.DefineViewRoutes(router, "")
	igmedia.DefineViewRoutes(router)

	api := router.Group("/api")
	jobqueue.DefineAPIRoutes(api)
	igprofile.DefineAPIRoutes(api)
	auth.DefineAPIRoutes(api)
	igmedia.DefineAPIRoutes(api)
	router.Run(":" + port)
}
