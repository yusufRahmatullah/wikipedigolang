package base

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Model is basic model which has base attributes
// such as ID and Timestamp
type Model struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CreatedAt  time.Time     `json:"created_at,omitempty" bson:"created_at"`
	ModifiedAt time.Time     `json:"modified_at,omitempty" bson:"modified_at"`
}

// InitBase initiates base attributes such as CreatedAt
// and ModifiedAt
func (model *Model) InitBase() {
	model.CreatedAt = time.Now()
	model.ModifiedAt = time.Now()
}

// NewModel returns new Model instance
func NewModel() Model {
	return Model{}
}
