package jobqueue

import (
	"fmt"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
)

const avaJobCol = "available_jobs"

var assignerLogger = logger.NewLogger("JobAssigner", true, true)

func init() {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(avaJobCol)
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := col.EnsureIndex(index)
	if err != nil {
		assignerLogger.Warning("Failed to create index")
	}
}

type avaJob struct {
	ID   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string        `json:"name" bson:"name"`
}

func getAvailableJobs() []avaJob {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(avaJobCol)
	var avaJobs []avaJob
	err := col.Find(nil).All(&avaJobs)
	if err != nil {
		assignerLogger.Warning("Failed to get all available jobs")
	}
	return avaJobs
}

// JobAssigner has objective to assign the JobQueue
// to corresponding Processor
type JobAssigner struct {
	ProcessorMap map[string]Job
}

// NewJobAssigner returns new JobAssigner instance
func NewJobAssigner() *JobAssigner {
	// Reset available job first
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(avaJobCol)
	ci, err := col.RemoveAll(bson.M{})
	if err == nil {
		assignerLogger.Info(fmt.Sprintf("Deleted jobs: %v", ci.Removed))
	} else {
		assignerLogger.Warning("Failed to remove available job")
	}
	return &JobAssigner{ProcessorMap: make(map[string]Job)}
}

// Register add the Processor to JobAssigner's ProcessorMap
func (ja *JobAssigner) Register(proc Job) {
	ja.ProcessorMap[proc.Name()] = proc
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(avaJobCol)
	err := col.Insert(&avaJob{Name: proc.Name()})
	if err == nil {
		assignerLogger.Warning(fmt.Sprintf("Success to insert available job: %s", proc.Name()))
	} else {
		assignerLogger.Warning(fmt.Sprintf("Failed to insert available job: %s", proc.Name()))
	}
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

	} else {
		assignerLogger.Info(fmt.Sprintf("%v not exist", name))
	}
	PostponeJobQueue(jobQueue)
	return false
}
