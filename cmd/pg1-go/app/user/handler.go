package user

import (
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"github.com/gin-gonic/gin"
)

// usersHandler handles the request to get all the Users
func usersHandler(c *gin.Context) {
	users := GetAll()
	data := base.StandardJSON("", users)
	c.JSON(http.StatusOK, data)
}

// findUserHandler handles the request to get a User by its name
func findUserHandler(c *gin.Context) {
	name := c.Param("name")
	user := FindUser(name)
	var data interface{}
	if user.Name == "" {
		data = base.StandardJSON("User not found", gin.H{})
	} else {
		data = base.StandardJSON("", user)
	}
	c.JSON(http.StatusOK, data)
}

// newUserHandler handles the request to post a new User
func newUserHandler(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	if user.Name == "" {
		data := base.ErrorJSON("name is required", gin.H{})
		c.JSON(http.StatusBadRequest, data)
	} else {
		suc := user.Save()
		var msg string
		var status int
		if suc {
			msg = "Create User successful"
			status = http.StatusCreated
		} else {
			msg = "Failed to create User"
			status = http.StatusOK
		}
		data := base.StandardJSON(msg, gin.H{})
		c.JSON(status, data)
	}
}

// deleteUserHandler handles the request to delete an exsiting User
func deleteUserHandler(c *gin.Context) {
	name := c.Param("name")
	suc := DeleteUser(name)
	var msg string
	if suc {
		msg = "Delete User successful"
	} else {
		msg = "Failed to delete User"
	}
	data := base.StandardJSON(msg, gin.H{})
	c.JSON(http.StatusOK, data)
}
