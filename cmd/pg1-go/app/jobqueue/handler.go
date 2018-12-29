package jobqueue

import (
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"github.com/gin-gonic/gin"
)

// newJobQueueHandler handles the request to post a new JobQueue
func newJobQueueHandler(c *gin.Context) {
	var jq JobQueue
	var data interface{}
	c.BindJSON(&jq)
	if jq.Name == "" {
		data = base.ErrorJSON("name is required", gin.H{})
		c.JSON(http.StatusBadRequest, data)
	} else {
		suc := jq.Save()
		var msg string
		var status int
		if suc {
			msg = "Create JobQueue successful"
			status = http.StatusCreated
		} else {
			msg = "Failed to create JobQueue"
			status = http.StatusOK
		}
		data := base.StandardJSON(msg, gin.H{})
		c.JSON(status, data)
	}
}
