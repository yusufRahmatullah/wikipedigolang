package repository

import (
	"git.heroku.com/pg1-go-work/cmd/igoweb/database"
	"git.heroku.com/pg1-go-work/cmd/igoweb/model"
	"github.com/globalsign/mgo/bson"
)

// UserRepository access User from repository
type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Get(username string) (*model.User, error)
}

// MongoUserRepository is implementation of UserRepository using MongoDB
type MongoUserRepository struct {
	DB *database.MongoClient
}

// Create store User to database
func (rep *MongoUserRepository) Create(user *model.User) error {
	col := rep.DB.Collection(database.UserCollection)
	user.InitTimeStamp()
	return col.Insert(user)
}

// Update update User in database
func (rep *MongoUserRepository) Update(user *model.User) error {
	panic("Not Implemented")
}

// Get read User from database by username
func (rep *MongoUserRepository) Get(username string) (*model.User, error) {
	col := rep.DB.Collection(database.UserCollection)
	var user model.User
	err := col.Find(bson.M{"username": username}).One(&user)
	return &user, err
}
