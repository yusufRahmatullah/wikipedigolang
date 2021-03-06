package jobqueue

import (
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/auth"
	"github.com/gin-gonic/gin"
)

// DefineAPIRoutes defines routes for JobQueue
func DefineAPIRoutes(router *gin.RouterGroup) {
	reqAdmin := router.Group("")
	reqAdmin.Use(auth.RequiredAdmin())
	{
		reqAdmin.GET("/available_jobs", getAvailableJobsHandler)
		reqAdmin.POST("/batch_add", batchAddHandler)
		reqAdmin.GET("/postponed_jobs/count", countPostponedJobsHandler)
		reqAdmin.GET("/postponed_jobs", getPostponedJobsHandler)
		reqAdmin.DELETE("/postponed_jobs/:job_id", deletePostponedHandler)
		reqAdmin.POST("/requeue_postponed_jobs", requeuePostponedHandler)
		reqAdmin.GET("/jobs", findJobQueueHandler)
		reqAdmin.POST("/jobs", newJobQueueHandler)
		reqAdmin.POST("/jobs/action", actionHandler)
	}
}

// DefineViewRoutes defines routes for JobQueue that contains view
func DefineViewRoutes(router *gin.Engine, prefix string) {
	reqAdmin := router.Group(prefix)
	reqAdmin.Use(auth.RequiredAdmin())
	{
		reqAdmin.GET("/batch_add", batchAddIndexView)
		reqAdmin.GET("/jobs", jobQueueIndexView)
		reqAdmin.GET("/postponed_jobs", postponedJobsView)
	}
}
