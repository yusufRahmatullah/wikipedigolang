package jobqueue

// Job is interface of base Job that process JobQueue
type Job interface {
	Name() string
	Process(jobQueue *JobQueue) string
}
