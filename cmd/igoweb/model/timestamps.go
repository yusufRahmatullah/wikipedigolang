package model

import "time"

// TimeStamp is the base model that has timestamp
// attributes such as CreatedAt and ModifiedAt
type TimeStamp struct {
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at"`
	ModifiedAt time.Time `json:"modified_at,omitempty" bson:"modified_at"`
}

// InitTimeStamp instantiate timestamp attributes
func (ts *TimeStamp) InitTimeStamp() {
	ts.CreatedAt = time.Now()
	ts.ModifiedAt = time.Now()
}
