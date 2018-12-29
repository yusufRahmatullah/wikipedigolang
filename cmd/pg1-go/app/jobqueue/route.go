package jobqueue

import "github.com/gin-gonic/gin"

// DefineRoutes defines routes for producer
func DefineRoutes(router *gin.Engine, prefix string) {
	router.POST(prefix+"/job_queue", newJobQueueHandler)
}
