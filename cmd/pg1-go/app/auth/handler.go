package auth

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	handleLogger = logger.NewLogger("AuthHandle", true, true)
	matcher      *regexp.Regexp
)

func init() {
	matcher = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_]+[A-Za-z0-9]$`)
}

// RequiredAdmin is the middleware that handles the
// request should be authenticated as user which has
// admin role
func RequiredAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("sid")
		if user != nil {
			username := user.(string)
			user := GetUser(username)
			if user.Role == adminRole {
				// Continue down the chain to handler etc
				c.Next()
			}
		}
		c.Redirect(http.StatusPermanentRedirect, "/login")
	}
}

func loginView(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("sid")
	if user != nil {
		c.Redirect(http.StatusPermanentRedirect, "/")
	}
	c.HTML(http.StatusOK, "login.tmpl.html", nil)
}

func validateLoginData(username, password string) string {
	if username == "" || password == "" {
		return "Parameters can't be empty"
	}
	if !matcher.MatchString(username) {
		return "Invalid username"
	}
	return ""
}

func login(c *gin.Context) {
	log.Println("==debug== login called")
	session := sessions.Default(c)
	loginData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "",
		Password: "",
	}
	c.BindJSON(&loginData)
	username := strings.Trim(loginData.Username, " ")
	password := strings.Trim(loginData.Password, " ")
	log.Println("==debug== up", username, password)
	errMsg := validateLoginData(username, password)
	log.Println("==debug== errMsg", errMsg)
	if errMsg != "" {
		data := base.ErrorJSON(errMsg, nil)
		c.JSON(http.StatusUnauthorized, data)
		return
	}
	user := GetUser(username)
	log.Println("==debug== user", user)
	if user == nil {
		handleLogger.Info(fmt.Sprintf("Failed to find username: %s", username))
		data := base.ErrorJSON("Invalid Username or Password", nil)
		c.JSON(http.StatusUnauthorized, data)
		return
	}
	suc := comparePasswords(user.Password, []byte(password))
	log.Println("==debug== suc", suc)
	if suc {
		session.Set("sid", username)
		err := session.Save()
		if err == nil {
			data := base.StandardJSON("Success", nil)
			c.JSON(http.StatusOK, data)
			return
		}
		handleLogger.Fatal("Failed to save session", err)
		data := base.ErrorJSON("Failed to save session", nil)
		c.JSON(http.StatusNotModified, data)
		return
	}
	data := base.ErrorJSON("Invalid Username or Password", nil)
	c.JSON(http.StatusUnauthorized, data)
	return
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("sid")
	if user != nil {
		session.Delete("sid")
		session.Save()
	}
	c.Redirect(http.StatusPermanentRedirect, "/")
}
