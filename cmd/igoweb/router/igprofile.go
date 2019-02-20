package router

import (
	"git.heroku.com/pg1-go-work/cmd/igoweb/handler"
	"git.heroku.com/pg1-go-work/cmd/igoweb/model"
	"github.com/gin-gonic/gin"
)

// IgProfileRouter redirect user request to corresponding handler
type IgProfileRouter struct {
	Handler *handler.IgProfileHandler
}

// DefineAPIRoutes defines routes for API call that return JSON response
func (router *IgProfileRouter) DefineAPIRoutes(rg *gin.RouterGroup) {
	rg.GET("/igprofiles", func(c *gin.Context) {
		router.Handler.FindIgProfiles(c, model.StatusActive)
	})
	rg.GET("/igprofiles/count", func(c *gin.Context) {
		router.Handler.CountIgProfiles(c, model.StatusActive)
	})
}
