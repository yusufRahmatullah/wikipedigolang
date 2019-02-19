package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// ProfileStatus is a status of IgProfile
type ProfileStatus string

const (
	// StatusActive means the IgProfile will be shown
	StatusActive ProfileStatus = "active"
	// StatusBanned means the IgProfile will not be shown
	StatusBanned ProfileStatus = "banned"
	// StatusMulti means the IgProfile will be shown on MultiAcc page
	// as active Multi Account
	StatusMulti ProfileStatus = "multi"
	// StatusAll means all IgProfile will be shown
	StatusAll ProfileStatus = ""
)

// IgProfile holds information about IG Profile
// include its IG ID, Name, followers number, following number,
// post number, and profile picture URL
type IgProfile struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CreatedAt  time.Time     `json:"created_at,omitempty" bson:"created_at"`
	ModifiedAt time.Time     `json:"modified_at,omitempty" bson:"modified_at"`
	IGID       string        `json:"ig_id" bson:"ig_id"`
	Name       string        `json:"name" bson:"name"`
	Followers  int           `json:"followers" bson:"followers"`
	Following  int           `json:"following" bson:"following"`
	Posts      int           `json:"posts" bson:"posts"`
	PpURL      string        `json:"pp_url" bson:"pp_url"`
	Status     ProfileStatus `json:"status" bson:"status"`
}
