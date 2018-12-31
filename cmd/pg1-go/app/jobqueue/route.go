package jobqueue

import "github.com/gin-gonic/gin"

// DefineAPIRoutes defines routes for JobQueue
func DefineAPIRoutes(router *gin.Engine, prefix string) {
	router.POST(prefix+"/job_queue", newJobQueueHandler)
}

// DefineViewRoutes defines routes for JobQueue that contains view
func DefineViewRoutes(router *gin.Engine, prefix string) {
	router.GET(prefix+"/job_queue", jobQueueIndexView)
}
