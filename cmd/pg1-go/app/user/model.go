package user

import (
	"time"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const userCol = "user"

func init() {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	index := mgo.Index{
		Key:        []string{"name", "phone"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := col.EnsureIndex(index)
	utils.HandleError(err, "", nil)
}

// User holds user information
type User struct {
	base.Model
	Name  string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
}

// NewUser returns new User instance
func NewUser(name string, phone string) *User {
	return &User{Name: name, Phone: phone}
}

// Save writes User instance to database
// or update if exist
// returns true if success
func (user *User) Save() bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	user.InitBase()
	selector := bson.M{"name": user.Name}
	update := bson.M{
		"$set": bson.M{
			"phone":       user.Phone,
			"created_at":  user.CreatedAt,
			"modified_at": time.Now(),
		},
	}
	_, err := col.Upsert(selector, update)
	utils.HandleError(err, "", nil)
	return err == nil
}

// GetAll returns All User in database
func GetAll() []User {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	var users []User
	err := col.Find(nil).All(&users)
	utils.HandleError(err, "", nil)
	return users
}

// FindUser find User instance in database by its name
func FindUser(name string) User {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	user := User{}
	err := col.Find(bson.M{"name": name}).One(&user)
	utils.HandleError(err, "", nil)
	return user
}

// DeleteUser removes User instance from database by its name
// returns true if success
func DeleteUser(name string) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	err := col.Remove(bson.M{"name": name})
	utils.HandleError(err, "", nil)
	return err == nil
}
