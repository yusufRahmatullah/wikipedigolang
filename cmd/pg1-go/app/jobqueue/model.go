package jobqueue

import (
	"fmt"
	"os"
	"strconv"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"
	"github.com/globalsign/mgo/bson"
)

const jobQueueCol = "job_queue"
const postponedJobCol = "postponed_job_queue"

var (
	// JobLimit is maximum Job to be processed and queried
	JobLimit    int
	modelLogger = logger.NewLogger("JobQueue", true, true)
)

func init() {
	var err error
	jobStr := os.Getenv("JOB_LIMIT")
	JobLimit, err = strconv.Atoi(jobStr)
	if err != nil {
		modelLogger.Warning("$JOB_LIMIT not found, using default")
		JobLimit = 10
	}
}

// JobQueue holds information about job function name
// and its params
type JobQueue struct {
	ID     bson.ObjectId          `json:"id" bson:"_id,omitempty"`
	Name   string                 `json:"name" bson:"name"`
	Params map[string]interface{} `json:"params" bson:"params"`
}

// NewJobQueue return new JobQueue instance
func NewJobQueue(name string, params map[string]interface{}) *JobQueue {
	return &JobQueue{Name: name, Params: params}
}

// Save writes JobQueue instance into database
// returns true if save success
func Save(jq *JobQueue) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	err := col.Insert(&jq)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to create JobQueue with name: %v", jq.Name))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to create JobQueue with name: %v", jq.Name))
	return false
}

// GetAll returns All JobQueue in database with
// maximum records defined by JobLimit
func GetAll() []JobQueue {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	var jobQueues []JobQueue
	err := col.Find(nil).Limit(JobLimit).All(&jobQueues)
	if err != nil {
		modelLogger.Fatal("Failed to get all JobQueue")
	}
	return jobQueues
}

// DeleteJobQueue remove JobQueue instance from database
// Returns true if remove is successful
func DeleteJobQueue(jobQueue *JobQueue) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	err := col.Remove(bson.M{"_id": jobQueue.ID})
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to delete JobQueue with name: %v", jobQueue.Name))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to delete JobQueue with name: %v", jobQueue.Name))
	return false
}

// PostponeJobQueue move the JobQueue to canceled job collections
// Returns true if postponing JobQueue is successful
func PostponeJobQueue(jobQueue *JobQueue) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(postponedJobCol)
	suc := DeleteJobQueue(jobQueue)
	if !suc {
		return suc
	}
	jobQueue.ID = ""
	err := col.Insert(jobQueue)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to postpone JobQueue with name: %v", jobQueue.Name))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to postpone JobQueue with name: %v", jobQueue.Name))
	return false
}

// GetPostponed returns Postponed JobQueue by its id
func GetPostponed(jid string) *JobQueue {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(postponedJobCol)
	var jq JobQueue
	err := col.FindId(bson.ObjectIdHex(jid)).One(&jq)
	if err != nil {
		modelLogger.Info(fmt.Sprintf("Failed to get postponed JobQueue with id: %v", jid))
	}
	return &jq
}

// GetAllPostponed returns all postponed JobQueue in database
// with maximum records defined by JobLimit
func GetAllPostponed() []JobQueue {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(postponedJobCol)
	var jobQueues []JobQueue
	err := col.Find(nil).Limit(JobLimit).All(&jobQueues)
	if err != nil {
		modelLogger.Fatal("Failed to get all JobQueue")
	}
	return jobQueues
}

// DeletePostponed remove postponed JobQueue instance from database
func DeletePostponed(jq *JobQueue) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(postponedJobCol)
	err := col.Remove(bson.M{"_id": jq.ID})
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to delete postponed JobQueue with name: %v", jq.Name))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to delete JobQueue with name: %v", jq.Name))
	return false
}

// RequeuePostponed move postponed JobQueue as new JobQueue
// returns true if success
func RequeuePostponed(jq *JobQueue) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	suc := DeletePostponed(jq)
	if !suc {
		return suc
	}
	jq.ID = ""
	err := col.Insert(jq)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to requeue JobQueue with name: %v", jq.Name))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to requeue JobQueue with name: %v", jq.Name))
	return false
}
