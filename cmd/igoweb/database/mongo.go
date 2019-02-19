package database

import (
	"fmt"
	"log"

	"git.heroku.com/pg1-go-work/cmd/igoweb/config"
	"github.com/globalsign/mgo"
)

// MongoClient is database client that use MongoDB
type MongoClient struct {
	Session *mgo.Session
}

// NewMongoClient instantiate MongoClient instance
func NewMongoClient() *MongoClient {
	mongoURL := config.GetInstance().MongoDBURL
	sess, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Panic(fmt.Sprintf("Failed to open MongoDB at URL: %v", mongoURL))
	}
	return &MongoClient{Session: sess}
}
