package igmedia

import (
	"fmt"
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"github.com/gin-gonic/gin"
)

var (
	handlerLogger = logger.NewLogger("IgMediaHandler", false, true)
)

func findIgMediaHandler(c *gin.Context, status MediaStatus) {
	fr := utils.GetFindRequest(c)
	igms := FindIgMedia(fr, status)
	compData := struct {
		Medias []IgMedia `json:"medias"`
		IGID   string    `json:"query"`
	}{
		Medias: igms,
		IGID:   fr.Query,
	}
	data := base.StandardJSON("", compData)
	c.JSON(http.StatusOK, data)
}

func findIgMediaByIDHandler(c *gin.Context, status MediaStatus) {
	fr := utils.GetFindRequest(c)
	igID := c.Query("ig_id")
	igms := findIgMediaByID(igID, fr, status)
	compData := struct {
		Medias []IgMedia `json:"medias"`
		IGID   string    `json:"query"`
	}{
		Medias: igms,
		IGID:   fr.Query,
	}
	data := base.StandardJSON("", compData)
	c.JSON(http.StatusOK, data)
}

func updateStatusHandler(c *gin.Context) {
	jd := struct {
		ID     string `json:"id"`
		Action string `json:"action"`
	}{
		ID:     "",
		Action: "",
	}
	c.BindJSON(&jd)
	if jd.ID == "" {
		data := base.ErrorJSON("Param ig_id can't be empty", nil)
		c.JSON(http.StatusBadRequest, data)
	}
	status := StatusShown
	switch jd.Action {
	case "activate":
		status = StatusShown
	case "hide":
		status = StatusHidden
	}
	suc := UpdateStatus(jd.ID, status)
	if suc {
		data := base.StandardJSON(fmt.Sprintf("Success to %v IgMedia with id: %v", jd.Action, jd.ID), nil)
		c.JSON(http.StatusOK, data)
	} else {
		data := base.ErrorJSON(fmt.Sprintf("Failed to %v IgMedia with id: %v", jd.Action, jd.ID), nil)
		c.JSON(http.StatusBadGateway, data)
	}
}

func countIgMediaHandler(c *gin.Context, status MediaStatus) {
	n, err := countIgMedia(status)
	if err == nil {
		data := base.StandardJSON("", n)
		c.JSON(http.StatusOK, data)
	} else {
		data := base.ErrorJSON("", err.Error())
		c.JSON(http.StatusInternalServerError, data)
	}
}

func igMediasView(c *gin.Context, admin bool) {
	c.HTML(http.StatusOK, "igpage.tmpl.html", gin.H{
		"admin":     admin,
		"title":     "IGO Medias",
		"search_ph": "Search IG ID here...",
		"sort_criteria": [][]string{
			[]string{"_id", "ID"},
			[]string{"ig_id", "IG ID"},
			[]string{"created_at", "Created At"},
			[]string{"modified_at", "Modified At"},
		},
		"filter_criteria": [][]string{
			[]string{"", "All"},
			[]string{"shown", "Shown"},
			[]string{"hidden", "Hidden"},
		},
		"script": gin.H{
			"name":         "igmedia",
			"page":         "igmedias",
			"subtitle":     "Medias",
			"search":       "",
			"count":        "/count",
			"admin_search": "/all",
			"admin_count":  "/count_all",
		},
	})
}
