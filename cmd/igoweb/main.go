package main

import (
	"html/template"

	"git.heroku.com/pg1-go-work/cmd/igoweb/config"
	"git.heroku.com/pg1-go-work/cmd/igoweb/database"
	"git.heroku.com/pg1-go-work/cmd/igoweb/handler"
	"git.heroku.com/pg1-go-work/cmd/igoweb/middleware"
	"git.heroku.com/pg1-go-work/cmd/igoweb/repository"
	"git.heroku.com/pg1-go-work/cmd/igoweb/router"
	"git.heroku.com/pg1-go-work/cmd/igoweb/service"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func setupMiddlewares(engine *gin.Engine) {
	engine.Use(gin.Logger())
	store := sessions.NewCookieStore([]byte(config.GetInstance().SessionSecret))
	engine.Use(sessions.Sessions("defaultSession", store))
	engine.Use(middleware.NoCacheHeader())
	engine.Use(cors.New(cors.Options{
		AllowedOrigins: config.GetInstance().Cors,
		AllowedMethods: []string{"GET"},
		Debug:          config.GetInstance().GinMode != "release",
	}))
}

func setupRouters(engine *gin.Engine) {
	mongo := database.NewMongoClient()
	loggerRepository := repository.MongoLogRepository{
		DB: mongo,
	}
	userRepository := repository.MongoUserRepository{
		DB: mongo,
	}
	igpLogger := service.LoggerService{
		Name:       "IG Profile",
		Repository: &loggerRepository,
	}
	igpRepository := repository.MongoIgProfileRepository{
		DB: mongo,
	}
	igpHandler := handler.IgProfileHandler{
		Logger:     &igpLogger,
		Repository: &igpRepository,
	}
	igpRouter := router.IgProfileRouter{
		Handler:        &igpHandler,
		UserRepository: &userRepository,
	}
	api := engine.Group("/api")
	igpRouter.DefineAPIRoutes(api)
}

func setupFuncMap(engine *gin.Engine) {
	engine.SetFuncMap(template.FuncMap{
		"ginMode": middleware.GinMode,
	})
}

func main() {
	engine := gin.New()
	setupMiddlewares(engine)
	setupRouters(engine)
	setupFuncMap(engine)
	engine.LoadHTMLGlob("templates/*.tmpl.html")
	engine.Static("/static", "static")
	engine.Run(":" + config.GetInstance().Port)
}
