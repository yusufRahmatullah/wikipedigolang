package igmedia

import (
	"fmt"
	"time"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"github.com/globalsign/mgo/bson"
)

// MediaStatus is a status of IgMedia
type MediaStatus string

const (
	igMediaCol = "ig_media"
	// StatusShown means the IgMedia will be shown
	StatusShown MediaStatus = "shown"
	// StatusHidden means the IgMedia will be hid
	StatusHidden MediaStatus = "hidden"
	// StatusAll means all IgMedia will be shown
	StatusAll MediaStatus = ""
)

var (
	modelLogger = logger.NewLogger("IgMedia", true, true)
)

// IgMedia holds information about media posted by IgProfile
// The information includes ig_id, media's url and status
type IgMedia struct {
	ID         string      `json:"id" bson:"_id,omitempty"`
	CreatedAt  time.Time   `json:"created_at,omitempty" bson:"created_at"`
	ModifiedAt time.Time   `json:"modified_at,omitempty" bson:"modified_at"`
	IGID       string      `json:"ig_id" bson:"ig_id"`
	URL        string      `json:"url" bson:"url"`
	Status     MediaStatus `json:"status" bson:"status"`
}

// NewIgMedia instantiate IgMedia with given ID, IGID and URL
// Status default as StatusShown
func NewIgMedia(id, igID, url string) *IgMedia {
	return &IgMedia{ID: id, IGID: igID, URL: url, Status: StatusShown}

}

func (model *IgMedia) initTime() {
	model.CreatedAt = time.Now()
	model.ModifiedAt = time.Now()
}

// Save writes IgMedia instance to database
// returns true if success
func Save(igm *IgMedia) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igMediaCol)
	igm.initTime()
	err := col.Insert(igm)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to create IgMedia with ID: %v", igm.ID))
		return true
	}
	modelLogger.Fatal(fmt.Sprintf("Failed to create IgMedia with ID: %v", igm.ID), err)
	return false
}

// FindIgMedia find IgMedia in database by IG ID
// Require offset and limit number for pagination
// Require status to define
func FindIgMedia(fr *utils.FindRequest, status MediaStatus) []IgMedia {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igMediaCol)
	var igms []IgMedia
	err := col.Find(bson.M{
		"ig_id":  bson.M{"$regex": fr.Query, "$options": "i"},
		"status": bson.M{"$regex": status, "$options": "i"},
	}).Sort(fr.Sort).Skip(fr.Offset).Limit(fr.Limit).All(&igms)
	if err != nil {
		modelLogger.Fatal(fmt.Sprintf("Failed to find IgMedia with igID: %v", fr.Query), err)
	}
	return igms
}

func countIgMedia(status MediaStatus) (int, error) {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igMediaCol)
	return col.Find(bson.M{"status": bson.M{"$regex": status}}).Count()
}

// UpdateStatus update IgMedia status
// returns true if success
func UpdateStatus(id string, status MediaStatus) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igMediaCol)
	err := col.UpdateId(id, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		modelLogger.Fatal(fmt.Sprintf("Failed to update IgMedia status with id: %v", id), err)
		return false
	}
	return true
}

// UpdateStatusAll update IgMedia status of all igID
// returns empty string if success
func UpdateStatusAll(igID string, status MediaStatus) string {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igMediaCol)
	info, err := col.UpdateAll(
		bson.M{"ig_id": igID},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err != nil {
		modelLogger.Fatal(fmt.Sprintf("Failed to update all IgMedia status with IG ID: %v", igID), err)
		return "Failed to update all IgMedia"
	}
	modelLogger.Info(fmt.Sprintf("Success to update %v IgMedia to %v", info.Updated, status))
	return ""
}
