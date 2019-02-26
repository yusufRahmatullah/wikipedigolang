package middleware

import (
	"net/http"

	"git.heroku.com/pg1-go-work/cmd/igoweb/model"

	"git.heroku.com/pg1-go-work/cmd/igoweb/repository"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RequireAdmin is the middleware that handles the
// request should be authenticated as user which has
// admin role
func RequireAdmin(userRep repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("sid")
		if user != nil {
			username := user.(string)
			user, err := userRep.Get(username)
			if err == nil && user.Role == model.RoleAdmin {
				c.Next()
			}
		}
		c.Redirect(http.StatusPermanentRedirect, "/login")
	}
}
