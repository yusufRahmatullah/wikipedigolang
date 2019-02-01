package igprofile

import (
	"fmt"
	"net/http"
	"strconv"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"github.com/gin-gonic/gin"
)

const (
	defaultOffset = 0
	defaultLimit  = 24
)

var handlerLogger = logger.NewLogger("IgProfileHandler", false, true)

func newIgProfileHandler(c *gin.Context) {
	var igp IgProfile
	c.BindJSON(&igp)
	if igp.IGID == "" {
		data := base.ErrorJSON("ig_id is required", nil)
		c.JSON(http.StatusBadRequest, data)
	} else {
		suc := Save(&igp)
		var msg string
		var status int
		if suc {
			msg = "Create IgProfile successful"
			status = http.StatusCreated
		} else {
			msg = "Failed to create IgProfile"
			status = http.StatusOK
		}
		data := base.StandardJSON(msg, nil)
		c.JSON(status, data)
	}
}

func generateSortOrder(c *gin.Context) string {
	sort := c.Query("sort")
	orderStr := c.Query("order")
	order := convertIntOrDefault(orderStr, -1)
	if sort == "" {
		sort = "_id"
	}
	if order == -1 {
		sort = "-" + sort
	}
	return sort
}

func getAllIgProfileHandler(c *gin.Context) {
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	offset := convertIntOrDefault(offsetStr, defaultOffset)
	limit := convertIntOrDefault(limitStr, defaultLimit)
	sort := generateSortOrder(c)
	igps := GetAll(offset, limit, sort)
	data := base.StandardJSON("", igps)
	c.JSON(http.StatusOK, data)
}

func getIgProfileHandler(c *gin.Context) {
	igID := c.Param("ig_id")
	igp := GetIgProfile(igID)
	var data map[string]interface{}
	if igp.IGID == "" {
		data = base.StandardJSON("User not found", nil)
	} else {
		data = base.StandardJSON("", igp)
	}
	c.JSON(http.StatusOK, data)
}

func modifyIgProfileHandler(c *gin.Context) {
	var igp IgProfile
	c.BindJSON(&igp)
	igID := c.Param("ig_id")
	changes := GenerateChanges(&igp)

	suc := Update(igID, changes)
	var msg string
	var status int
	if suc {
		msg = "Update IgProfile successful"
		status = http.StatusCreated
	} else {
		msg = "Failed to update IgProfile"
		status = http.StatusOK
	}
	data := base.StandardJSON(msg, nil)
	c.JSON(status, data)
}

func convertIntOrDefault(text string, def int) int {
	if text == "" {
		return def
	}
	num, err := strconv.Atoi(text)
	if err != nil {
		handlerLogger.Warning(fmt.Sprintf("Failed to convert text to int: %v", text))
		return def
	}
	return num
}

func findIgProfileHandler(c *gin.Context) {
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	query := c.Query("query")
	offset := convertIntOrDefault(offsetStr, defaultOffset)
	limit := convertIntOrDefault(limitStr, defaultLimit)
	sort := generateSortOrder(c)
	igps := FindIgProfile(query, offset, limit, sort)
	data := base.StandardJSON("", igps)
	c.JSON(http.StatusOK, data)
}

func deleteIgProfileHandler(c *gin.Context) {
	igID := c.Param("ig_id")
	suc := DeleteIgProfile(igID)
	var msg string
	if suc {
		msg = "Delete IgProfile successful"
	} else {
		msg = "Failed to delete IgProfile"
	}
	data := base.StandardJSON(msg, nil)
	c.JSON(http.StatusOK, data)
}

/////////////////////////////////
// IgProfile Views
/////////////////////////////////

func igProfilesView(c *gin.Context) {
	c.HTML(http.StatusOK, "igprofiles.tmpl.html", nil)
}
