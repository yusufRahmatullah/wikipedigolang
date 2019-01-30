package auth

import (
	"fmt"
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

// Required is the middleware that handles the
// request should be authenticated
func Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusUnauthorized, "/login")
		} else {
			// Continue down the chain to handler etc
			c.Next()
		}
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

func login(c *gin.Context) {
	session := sessions.Default(c)
	loginData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "",
		Password: "",
	}
	c.BindJSON(&loginData)
	fmt.Printf("==debug== loginData: %v", loginData)
	username := strings.Trim(loginData.Username, " ")
	password := strings.Trim(loginData.Password, " ")

	if username == "" || password == "" {
		data := base.ErrorJSON("Parameters can't be empty", nil)
		c.JSON(http.StatusUnauthorized, data)
		return
	}
	if !matcher.MatchString(username) {
		data := base.ErrorJSON("Invalid username", nil)
		c.JSON(http.StatusUnauthorized, data)
		return
	}
	user := GetUser(username)
	if user == nil {
		handleLogger.Info(fmt.Sprintf("Failed to find username: %s", username))
		data := base.ErrorJSON("Invalid Username or Password", nil)
		c.JSON(http.StatusUnauthorized, data)
	} else {
		suc := comparePasswords(user.Password, []byte(password))
		if suc {
			session.Set("sid", username)
			err := session.Save()
			if err == nil {
				data := base.StandardJSON("Success", nil)
				c.JSON(http.StatusOK, data)
			} else {
				handleLogger.Fatal("Failed to save session")
				data := base.ErrorJSON("Failed to save session", nil)
				c.JSON(http.StatusNotModified, data)
			}
		} else {
			data := base.ErrorJSON("Invalid Username or Password", nil)
			c.JSON(http.StatusUnauthorized, data)
		}
	}
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
