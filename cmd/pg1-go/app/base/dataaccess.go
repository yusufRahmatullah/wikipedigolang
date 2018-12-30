package base

import (
	"fmt"
	"log"
	"os"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"
	"github.com/globalsign/mgo"
)

var (
	mongodbURL string
)

// DataAccess holds information for data access to database
type DataAccess struct {
	Client *mgo.Session
}

func init() {
	mongodbURL = os.Getenv("MONGODB_URL")
	if mongodbURL == "" {
		log.Fatal("$MONGODB_URL must me set")
	}
}

// NewDataAccess return new DataAccess instance with instantiated Client
func NewDataAccess() *DataAccess {
	dbClient, err := mgo.Dial(mongodbURL)
	utils.HandleError(
		err,
		fmt.Sprintf("Failed to open MongoDB at URL: %v\n", mongodbURL),
		func() { dbClient.Close() },
	)
	return &DataAccess{Client: dbClient}
}

// GetCollection returns MongoDB collection from DataAccess
func (dataAccess *DataAccess) GetCollection(name string) *mgo.Collection {
	return dataAccess.Client.DB("").C(name)
}

// Close terminate DataAccess Client
func (dataAccess *DataAccess) Close() {
	dataAccess.Client.Close()
}
