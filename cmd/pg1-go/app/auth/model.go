package auth

import (
	"fmt"
	"os"
	"time"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

type userRole string

// User holds information about user that access the app
type User struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username   string        `bson:"username" json:"username"`
	Password   string        `bson:"password"`
	Role       userRole      `bson:"role" json:"role"`
	CreatedAt  time.Time     `json:"created_at,omitempty" bson:"created_at"`
	ModifiedAt time.Time     `json:"modified_at,omitempty" bson:"modified_at"`
}

const (
	userCol             = "user"
	adminRole  userRole = "admin"
	commonRole userRole = "common"
)

var (
	modelLogger = logger.NewLogger("User", true, true)
)

func init() {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := col.EnsureIndex(index)
	if err != nil {
		modelLogger.Warning("Failed to create index")
	}

	// Create superadmin
	sapass := os.Getenv("SUPER_ADMIN")
	if sapass == "" {
		modelLogger.Error("$SUPER_ADMIN must be set", nil)
	}
	superAdmin := NewUser("superadmin", sapass)
	superAdmin.Role = adminRole
	err = col.Insert(&superAdmin)
	if err != nil {
		modelLogger.Info(fmt.Sprintf("Failed to create superadmin. cause: %v", err))
	}
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		modelLogger.Fatal(fmt.Sprintf("Failed to generate hash from password"), err)
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {

	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		modelLogger.Info(fmt.Sprintf("Failed to compare hash with password cause: %v", err))
		return false
	}

	return true

}

func (model *User) initTime() {
	model.CreatedAt = time.Now()
	model.ModifiedAt = time.Now()
}

// NewUser return new User instance
func NewUser(username, password string) *User {
	hashPass := hashAndSalt([]byte(password))
	return &User{Username: username, Password: hashPass, Role: commonRole}
}

// Save writes User isntance to database
// returns true if success
func Save(user *User) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	user.initTime()
	err := col.Insert(user)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to create User with username: %v", user.Username))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to create User with username: %v", user.Username))
	return false
}

// Update modify User instance in database
// returns true if success
func Update(username string, changes map[string]interface{}) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	selector := bson.M{"username": username}
	changes["modified_at"] = time.Now()
	update := bson.M{"$set": changes}
	err := col.Update(selector, update)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to update User with username: %v", username))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to update User with username: %v", username))
	return false
}

// GenerateChanges build map of non-empty User's attribute
func GenerateChanges(user *User) map[string]interface{} {
	changes := gin.H{}
	if user.Password != "" {
		changes["password"] = user.Password
	}
	if user.Role != "" {
		changes["role"] = user.Role
	}
	return changes
}

// GetAll returns All User in database
// Require offset and limit number for pagination
func GetAll(offset, limit int, sortBy ...string) []User {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	var users []User
	if len(sortBy) == 0 {
		sortBy = []string{"_id"}
	}
	err := col.Find(nil).Sort(sortBy...).Skip(offset).Limit(limit).All(&users)
	if err != nil {
		modelLogger.Info("Failed to get all Users")
	}
	return users
}

// GetUser get User instance in database by username
func GetUser(username string) *User {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	var user User
	err := col.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		modelLogger.Info(fmt.Sprintf("Failed to get User with username: %v", username))
	}
	return &user
}

// DeleteUser removes User instance from database
// return true if success
func DeleteUser(username string) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(userCol)
	err := col.Remove(bson.M{"username": username})
	if err != nil {
		modelLogger.Info(fmt.Sprintf("Failed to delete User with username: %s", username))
		return false
	}
	return true
}
