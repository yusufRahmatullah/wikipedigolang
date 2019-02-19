package model

import (
	"github.com/globalsign/mgo/bson"
)

// SuccessJSON generates JSON for success response
func SuccessJSON(message string, data interface{}) interface{} {
	return bson.M{
		"status":  "OK",
		"message": message,
		"data":    data,
	}
}

// ErrorJSON generates JSON for error response
func ErrorJSON(message string, err error) interface{} {
	return bson.M{
		"status":  "error",
		"message": message,
		"error":   err,
	}
}
