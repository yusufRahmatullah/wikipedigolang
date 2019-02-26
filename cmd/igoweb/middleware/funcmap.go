package middleware

import (
	"git.heroku.com/pg1-go-work/cmd/igoweb/config"
)

// GinMode is Template Function that returns GinMode from config
func GinMode() string {
	return config.GetInstance().GinMode
}
