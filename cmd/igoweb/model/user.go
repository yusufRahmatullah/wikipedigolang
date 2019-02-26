package model

import (
	"git.heroku.com/pg1-go-work/cmd/igoweb/util"
)

// UserRole represents User's role/privileges
type UserRole string

// User holds information about user that access the app
type User struct {
	TimeStamp
	MongoID
	Username string   `bson:"username" json:"username"`
	Password string   `bson:"password"`
	Role     UserRole `bson:"role" json:"role"`
}

const (
	// RoleAdmin has privilege to view admin pages
	RoleAdmin UserRole = "admin"
	// RoleCommon has privilege to view public pages
	RoleCommon UserRole = "common"
)

// NewUser instantiate User using given username and password
func NewUser(username, password string) (*User, error) {
	hashPass, err := util.HashAndSalt([]byte(password))
	var user User
	if err == nil {
		user = User{
			Username: username,
			Password: hashPass,
			Role:     RoleCommon,
		}
	}
	return &user, err
}
