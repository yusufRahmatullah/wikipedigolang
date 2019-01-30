package jobqueue

import (
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/auth"
	"github.com/gin-gonic/gin"
)

// DefineAPIRoutes defines routes for JobQueue
func DefineAPIRoutes(router *gin.Engine, prefix string) {
	reqAdmin := router.Group(prefix)
	reqAdmin.Use(auth.RequiredAdmin())
	{
		reqAdmin.POST("/job_queue", newJobQueueHandler)
		reqAdmin.GET("/available_jobs", getAvailableJobsHandler)
	}
}

// DefineViewRoutes defines routes for JobQueue that contains view
func DefineViewRoutes(router *gin.Engine, prefix string) {
	reqAdmin := router.Group(prefix)
	reqAdmin.Use(auth.RequiredAdmin())
	{
		reqAdmin.GET("/job_queue", jobQueueIndexView)
	}
}
