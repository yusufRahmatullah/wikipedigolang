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

// Collection returns Collection instance to interact with database
func (mc *MongoClient) Collection(name string) *mgo.Collection {
	return mc.Session.DB("").C(name)
}

// Prepare preparing MongoClient before can be used
func (mc *MongoClient) Prepare() error {
	var err error
	err = mc.prepareIgProfileIndex()
	err = mc.prepareLogIndex()
	err = mc.prepareUserIndex()
	if err != nil {
		return err
	}
	return err
}

// Close end the MongoClient connection
func (mc *MongoClient) Close() {
	mc.Session.Close()
}

func (mc *MongoClient) prepareIgProfileIndex() error {
	col := mc.Collection(IgProfileCollection)
	idIndex := mgo.Index{
		Key:        []string{"ig_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	return col.EnsureIndex(idIndex)
}

func (mc *MongoClient) prepareLogIndex() error {
	col := mc.Collection(LogCollection)
	index := mgo.Index{
		Key:        []string{"name", "level"},
		Background: true,
	}
	return col.EnsureIndex(index)
}

func (mc *MongoClient) prepareUserIndex() error {
	col := mc.Collection(LogCollection)
	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	return col.EnsureIndex(index)
}
