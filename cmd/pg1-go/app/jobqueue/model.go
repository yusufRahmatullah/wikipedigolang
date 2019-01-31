package jobqueue

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/globalsign/mgo"

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
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	pcol := dataAccess.GetCollection(postponedJobCol)
	index := mgo.Index{
		Key:        []string{"unique_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = col.EnsureIndex(index)
	if err != nil {
		modelLogger.Warning("Failed to create index on JobQueue Collections")
	}
	err = pcol.EnsureIndex(index)
	if err != nil {
		modelLogger.Warning("Failed to create index on Postponed JobQueue Collections")
	}
}

// JobQueue holds information about job function name
// and its params
type JobQueue struct {
	ID       bson.ObjectId          `json:"id" bson:"_id,omitempty"`
	Name     string                 `json:"name" bson:"name"`
	Params   map[string]interface{} `json:"params" bson:"params"`
	UniqueID string                 `json:"unique_id" bson:"unique_id"`
}

func sortedKeys(params map[string]interface{}) []string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// GenerateUniqueID of JobQueue
func (jq *JobQueue) GenerateUniqueID() {
	keys := sortedKeys(jq.Params)
	var sb strings.Builder
	fmt.Println(sb.WriteString(jq.Name))
	for _, k := range keys {
		v := jq.Params[k]
		sb.WriteString("::")
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString(v.(string))
	}
	jq.UniqueID = sb.String()
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
	jq.GenerateUniqueID()
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
	modelLogger.Info(fmt.Sprintf("Failed to postpone JobQueue with name: %v cause: %v", jobQueue.Name, err))
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
