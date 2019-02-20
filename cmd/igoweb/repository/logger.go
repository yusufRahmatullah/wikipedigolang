package repository

import (
	"git.heroku.com/pg1-go-work/cmd/igoweb/database"
	"git.heroku.com/pg1-go-work/cmd/igoweb/model"
)

// LogRepository access Log from repository
type LogRepository interface {
	Create(log *model.Log)
}

// MongoLogRepository is implementation of LogRepository using
// MongoDB as the database
type MongoLogRepository struct {
	DB *database.MongoClient
}

// Create store Log to database
func (rep *MongoLogRepository) Create(log *model.Log) {
	col := rep.DB.Collection(database.LogCollection)
	log.InitTimeStamp()
	col.Insert(log)
}
