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

// JobStatus is a status of JobQueue
type JobStatus string

const (
	jobQueueCol = "job_queue"
	// StatusActive means JobQueue will be executed
	StatusActive JobStatus = "active"
	// StatusFinished means JobQueue has been executed
	StatusFinished JobStatus = "finished"
	// StatusPostponed means JobQueue will be executed later
	StatusPostponed JobStatus = "postponed"
)

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
	index := mgo.Index{
		Key:        []string{"unique_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = col.EnsureIndex(index)
	if err != nil {
		modelLogger.Warning(fmt.Sprintf("Failed to create index on %v", jobQueueCol))
	}
}

// JobQueue holds information about job function name
// and its params
type JobQueue struct {
	ID       bson.ObjectId          `json:"id" bson:"_id,omitempty"`
	Name     string                 `json:"name" bson:"name"`
	Params   map[string]interface{} `json:"params" bson:"params"`
	UniqueID string                 `json:"unique_id" bson:"unique_id"`
	Status   JobStatus              `json:"status" bson:"status"`
	Reason   string                 `json:"reason" bson:"reason"`
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
	sb.WriteString(jq.Name)
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
	return &JobQueue{Name: name, Params: params, Status: StatusActive}
}

// Save writes JobQueue instance into database
// returns empty string if save success
func Save(jq *JobQueue) string {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	jq.GenerateUniqueID()
	jq.Status = StatusActive
	err := col.Insert(&jq)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to create JobQueue with name: %v", jq.Name))
		return ""
	}
	modelLogger.Fatal(fmt.Sprintf("Failed to create JobQueue with name: %v", jq.Name), err)
	return "Failed to create JobQueue"
}

// Update modify JobQueue instance in database
// returns true if success
func Update(jq *JobQueue, changes bson.M) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	update := bson.M{"$set": changes}
	err := col.UpdateId(jq.ID, update)
	if err != nil {
		modelLogger.Fatal(fmt.Sprintf("Failed to update JobQueue with id: %v", jq.ID), err)
	}
	return err == nil
}

// GetAll returns All active JobQueue in database with
// maximum records defined by JobLimit
func GetAll() []JobQueue {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	var jobQueues []JobQueue
	err := col.Find(bson.M{
		"status": StatusActive,
	}).Limit(JobLimit).All(&jobQueues)
	if err != nil {
		modelLogger.Fatal("Failed to get all JobQueue", err)
	}
	return jobQueues
}

// DeleteJobQueue remove JobQueue instance from database
// Returns true if remove is successful
func DeleteJobQueue(jobQueue *JobQueue) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	err := col.RemoveId(jobQueue.ID)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to delete JobQueue with name: %v", jobQueue.Name))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to delete JobQueue with name: %v", jobQueue.Name))
	return false
}

// PostponeJobQueue move the JobQueue to canceled job collections
// Requires reason why JobQueue is postponed
// Returns true if postponing JobQueue is successful
func PostponeJobQueue(jq *JobQueue, reason string) bool {
	return Update(jq, bson.M{
		"status": StatusPostponed,
		"reason": reason,
	})
}

// GetPostponed returns Postponed JobQueue by its id
func GetPostponed(id string) *JobQueue {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	var jq JobQueue
	err := col.Find(bson.M{
		"_id":    bson.ObjectIdHex(id),
		"status": StatusPostponed,
	}).One(&jq)
	if err != nil {
		modelLogger.Info(fmt.Sprintf("Failed to get postponed JobQueue with id: %v", id))
	}
	return &jq
}

// GetAllPostponed returns all postponed JobQueue in database
// with maximum records defined by JobLimit
func GetAllPostponed() []JobQueue {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	var jobQueues []JobQueue
	err := col.Find(bson.M{
		"status": StatusPostponed,
	}).Limit(JobLimit).All(&jobQueues)
	if err != nil {
		modelLogger.Fatal("Failed to get all JobQueue", err)
	}
	return jobQueues
}

func countPostponedJobs() (int, error) {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	return col.Find(bson.M{"status": StatusPostponed}).Count()
}

// DeletePostponed remove postponed JobQueue instance from database
func DeletePostponed(jq *JobQueue) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	update := bson.M{"$set": bson.M{"status": StatusFinished}}
	err := col.UpdateId(jq.ID, update)
	if err != nil {
		modelLogger.Fatal(fmt.Sprintf("Failed to delete JobQueue with id: %v", jq.ID), err)
	}
	return err == nil
}

// RequeuePostponed move postponed JobQueue as new JobQueue
// returns true if success
func RequeuePostponed(jq *JobQueue) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(jobQueueCol)
	err := col.UpdateId(jq.ID, bson.M{"$set": bson.M{"status": StatusActive}})
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to requeue JobQueue with uid: %v", jq.UniqueID))
		return true
	}
	modelLogger.Fatal(fmt.Sprintf("Failed to requeue JobQueue with uid: %v", jq.UniqueID), err)
	return false
}
