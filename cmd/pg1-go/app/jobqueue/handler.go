package jobqueue

import (
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"github.com/gin-gonic/gin"
)

func getAvailableJobsHandler(c *gin.Context) {
	avaJobs := getAvailableJobs()
	data := base.StandardJSON("", avaJobs)
	c.JSON(200, data)
}

// newJobQueueHandler handles the request to post a new JobQueue
func newJobQueueHandler(c *gin.Context) {
	var jq JobQueue
	var data interface{}
	c.BindJSON(&jq)
	if jq.Name == "" {
		data = base.ErrorJSON("name is required", nil)
		c.JSON(http.StatusBadRequest, data)
	} else {
		suc := Save(&jq)
		var msg string
		var status int
		if suc {
			msg = "Create JobQueue successful"
			status = http.StatusCreated
		} else {
			msg = "Failed to create JobQueue"
			status = http.StatusOK
		}
		data := base.StandardJSON(msg, nil)
		c.JSON(status, data)
	}
}

// jobQueueIndexView render JobQueue form
func jobQueueIndexView(c *gin.Context) {
	c.HTML(http.StatusOK, "jobqueue.tmpl.html", nil)
}
