package jobqueue

import (
	"os"
	"strconv"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"
	"github.com/globalsign/mgo/bson"
)

const jobQueueCol = "job_queue"
const postponedJobCol = "postponed_job_queue"

var (
	// JobLimit is maximum Job to be processed and queried
	JobLimit int
)

func init() {
	var err error
	jobStr := os.Getenv("JOB_LIMIT")
	JobLimit, err = strconv.Atoi(jobStr)
	utils.HandleError(
		err,
		"$JOB_LIMIT not found. using default",
		func() {
			JobLimit = 10
		},
	)
}

// JobQueue holds information about job function name
// and its params
type JobQueue struct {
	ID     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string        `json:"name" bson:"name"`
	Params interface{}   `json:"params" bson:"params"`
}

// NewJobQueue return new JobQueue instance
func NewJobQueue(name string, params interface{}) JobQueue {
	return JobQueue{Name: name, Params: params}
}

// Save writes JobQueue instance into database
// returns true if save success
func (jq *JobQueue) Save() bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	err := col.Insert(&jq)
	utils.HandleError(err, "", nil)
	return err == nil
}

// GetAll returns All JobQueue in database with
// maximum records defined by JobLimit
func GetAll() []JobQueue {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	var jobQueues []JobQueue
	err := col.Find(nil).Limit(JobLimit).All(&jobQueues)
	utils.HandleError(err, "", nil)
	return jobQueues
}

// DeleteJobQueue remove JobQueue instance from database
// Returns true if remove is successful
func DeleteJobQueue(jobQueue JobQueue) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	err := col.Remove(bson.M{"_id": jobQueue.ID})
	utils.HandleError(err, "", nil)
	return err == nil
}

// DeleteJobQueueID remove JobQueue instance from database
// defined by its ID. Returns true if remove is successful
func DeleteJobQueueID(id bson.ObjectId) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	err := col.RemoveId(id)
	utils.HandleError(err, "", nil)
	return err == nil
}

// PostponeJobQueue move the JobQueue to canceled job collections
// Returns true if postponing JobQueue is successful
func PostponeJobQueue(jobQueue JobQueue) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(postponedJobCol)
	DeleteJobQueue(jobQueue)
	jobQueue.ID = ""
	err := col.Insert(jobQueue)
	utils.HandleError(err, "", nil)
	return err == nil
}
