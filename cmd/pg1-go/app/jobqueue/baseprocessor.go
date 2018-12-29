package jobqueue

import "log"

// BaseProcessor is base class of JobQueue processor
// Each JobQueue Processor has its own unique Name
type BaseProcessor struct {
	Name string
}

// NewBaseProcessor instantiate BaseProcessor instance
// and assign its Name with "BaseProcessor"
func NewBaseProcessor() BaseProcessor {
	return BaseProcessor{Name: "BaseProcessor"}
}

// Process executes job queue with the given params
// Returns true if process success
func (jp *BaseProcessor) Process(jobQueue JobQueue) bool {
	params := jobQueue.Params
	log.Printf("[base.JobProcessor] params: %v\n", params)
	return true
}
