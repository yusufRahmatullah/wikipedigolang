package igpost

import (
	"fmt"
	"time"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	delPostCol = "deleted_ig_post"
	postCol    = "ig_post"
)

var modelLogger = logger.NewLogger("IgPost", true, true)

func init() {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(postCol)
	idIndex := mgo.Index{
		Key:        []string{"post_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := col.EnsureIndex(idIndex)
	if err != nil {
		modelLogger.Warning("Failed to create index")
	}
}

// IgPost holds information about post of an IG Profile
type IgPost struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CreatedAt  time.Time     `json:"created_at,omitempty" bson:"created_at"`
	ModifiedAt time.Time     `json:"modified_at,omitempty" bson:"modified_at"`
	PostID     string        `json:"post_id" bson:"post_id"`
}

// NewIgPost instantiate new IgPost
func NewIgPost(pID string) *IgPost {
	return &IgPost{PostID: pID}
}

func (model *IgPost) initTime() {
	model.CreatedAt = time.Now()
	model.ModifiedAt = time.Now()
}

// Save writes IgPost instance to database
// returns true if success
func Save(ip *IgPost) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(postCol)
	ip.initTime()
	err := col.Insert(ip)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to create IgPost with Post ID: %v", ip.PostID))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to create IgPost with Post ID: %v", ip.PostID))
	return false
}

// GetAll returns all IgPost in database
// Require offset and limit number for pagination
func GetAll(offset, limit int) []IgPost {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(postCol)
	var ips []IgPost
	err := col.Find(nil).Skip(offset).Limit(limit).All(&ips)
	if err == nil {
		modelLogger.Debug("Success to get all IgPost")
	} else {
		modelLogger.Fatal("Failed to get all IgPost")
	}
	return ips
}

// GetIgPost returns IgPost instance in database by its PostID
func GetIgPost(pID string) *IgPost {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(postCol)
	var igp IgPost
	err := col.Find(bson.M{"post_id": pID}).One(&igp)
	if err == nil {
		modelLogger.Debug(fmt.Sprintf("Success to get IgPost with Post ID: %v", pID))
	} else {
		modelLogger.Debug(fmt.Sprintf("Failed to get IgPost with Post ID: %v", pID))
	}
	return &igp
}

// Delete remove IgPost instance in database by Post ID
// returns true if success
func Delete(pID string) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(postCol)
	err := col.Remove(bson.M{"post_id": pID})
	if err != nil {
		modelLogger.Info(fmt.Sprintf("Failed to delete IgPost with Post ID: %v", pID))
		return false
	}
	modelLogger.Info(fmt.Sprintf("Success to delete IgPost with Post ID: %v", pID))
	return true
}
