package igmedia

import (
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/auth"
	"github.com/gin-gonic/gin"
)

// DefineAPIRoutes defines routes for IgMedia
func DefineAPIRoutes(router *gin.RouterGroup) {
	router.GET("/igmedias", func(c *gin.Context) {
		findIgMediaHandler(c, StatusShown)
	})
	router.GET("/igmedias/by_id", func(c *gin.Context) {
		findIgMediaByIDHandler(c, StatusShown)
	})
	router.GET("/igmedias/count", func(c *gin.Context) {
		countIgMediaHandler(c, StatusShown)
	})
	reqAdmin := router.Group("")
	reqAdmin.Use(auth.RequiredAdmin())
	{
		router.GET("/igmedias/by_id_all", func(c *gin.Context) {
			findIgMediaByIDHandler(c, StatusAll)
		})
		reqAdmin.GET("igmedias/count_all", func(c *gin.Context) {
			countIgMediaHandler(c, StatusAll)
		})
		reqAdmin.GET("/igmedias/all", func(c *gin.Context) {
			filterStatus := c.Query("filterStatus")
			status := StatusAll
			switch filterStatus {
			case "shown":
				status = StatusShown
			case "hidden":
				status = StatusHidden
			}
			findIgMediaHandler(c, status)
		})
		reqAdmin.POST("/igmedias/action", updateStatusHandler)
	}
}

// DefineViewRoutes defines routes for IgMedia that contains view
func DefineViewRoutes(router *gin.Engine) {
	router.GET("/igmedias", func(c *gin.Context) {
		igMediasView(c, false)

	})
	adminPage := router.Group("/admin")
	adminPage.Use(auth.RequiredAdmin())
	{
		adminPage.GET("/igmedias", func(c *gin.Context) {
			igMediasView(c, true)
		})
	}
}
