package handler

import (
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/igoweb/model"
	"git.heroku.com/pg1-go-work/cmd/igoweb/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// IgProfileHandler handles request about IgProfile
type IgProfileHandler struct {
	Repository *repository.IgProfileRepository
	Logger     *zap.Logger
}

// CreateIgProfile handles request to create new IgProfile
func (handler *IgProfileHandler) CreateIgProfile(c *gin.Context) {
	var igProfile *model.IgProfile
	c.BindJSON(&igProfile)
	if igProfile.IGID == "" {
		resp := model.ErrorJSON("Param ig_id is required", nil)
		c.JSON(http.StatusBadRequest, resp)
	} else {
		err := handler.Repository.Create(igProfile)
		if err != nil {
			resp := model.ErrorJSON("Failed to create IgProfile", err)
			c.JSON(http.StatusBadRequest, resp)
		} else {
			resp := model.SuccessJSON("Create IgProfile successful", nil)
			c.JSON(http.StatusCreated, resp)
		}
	}
}
