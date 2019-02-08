package igmedia

import (
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/auth"
	"github.com/gin-gonic/gin"
)

// DefineAPIRoutes defines routes for IgMedia
func DefineAPIRoutes(router *gin.RouterGroup) {
	router.GET("/igmedias", func(c *gin.Context) {
		findIgMediaHandler(c, StatusShown)
	})
	reqAdmin := router.Group("")
	reqAdmin.Use(auth.RequiredAdmin())
	{
		reqAdmin.GET("/igmedias/all", func(c *gin.Context) {
			findIgMediaHandler(c, StatusAll)
		})
		reqAdmin.POST("/igmedias/action", updateStatusHandler)
	}
}

// DefineViewRoutes defines routes for IgMedia that contains view
func DefineViewRoutes(router *gin.Engine) {
	router.GET("/igmedias", func(c *gin.Context) {
		c.HTML(http.StatusOK, "igmedia.tmpl.html", nil)
	})
	adminPage := router.Group("/admin")
	adminPage.Use(auth.RequiredAdmin())
	{
		adminPage.GET("/igmedias", func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin_igmedia.tmpl.html", nil)
		})
	}
}
