package model

import "github.com/globalsign/mgo/bson"

// MongoID is the base model which ID attribute is MongoDB's ID
type MongoID struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
}
