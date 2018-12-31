package jobqueue

import (
	"fmt"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
)

var assignerLogger = logger.NewLogger("JobAssigner", true, true)

// JobAssigner has objective to assign the JobQueue
// to corresponding Processor
type JobAssigner struct {
	ProcessorMap map[string]Job
}

// NewJobAssigner returns new JobAssigner instance
func NewJobAssigner() *JobAssigner {
	return &JobAssigner{ProcessorMap: make(map[string]Job)}
}

// Register add the Processor to JobAssigner's ProcessorMap
func (ja *JobAssigner) Register(proc Job) {
	ja.ProcessorMap[proc.Name()] = proc
}

// ProcessJobQueue process JobQueue by assign it into
// corresponding registered Processor
// Returns true if JobQueue processed succesfully
func (ja *JobAssigner) ProcessJobQueue(jobQueue *JobQueue) bool {
	name := jobQueue.Name
	params := jobQueue.Params
	proc, exist := ja.ProcessorMap[name]
	if exist {
		suc := proc.Process(jobQueue)
		if suc {
			assignerLogger.Info(fmt.Sprintf("Success to process %v", name))
			return DeleteJobQueue(jobQueue)
		}
		assignerLogger.Info(fmt.Sprintf("Failed to process %v with params: %v", name, params))

	}
	assignerLogger.Info(fmt.Sprintf("%v not exist", name))
	PostponeJobQueue(jobQueue)
	return false
}
