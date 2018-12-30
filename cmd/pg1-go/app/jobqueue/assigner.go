package jobqueue

import "log"

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
	proc, exist := ja.ProcessorMap[name]
	if exist {
		suc := proc.Process(jobQueue)
		if suc {
			return DeleteJobQueue(*jobQueue)
		}
	}
	for k := range ja.ProcessorMap {
		log.Printf("%v\n", k)
	}
	suc := PostponeJobQueue(*jobQueue)
	if !suc {
		log.Printf(
			"Failed to postpone job queue: {Name: %v, Params: %v}\n",
			jobQueue.Name,
			jobQueue.Params,
		)
	}
	return false
}
