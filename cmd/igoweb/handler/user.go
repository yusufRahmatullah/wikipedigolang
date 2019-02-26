package handler

import (
	"git.heroku.com/pg1-go-work/cmd/igoweb/config"
	"git.heroku.com/pg1-go-work/cmd/igoweb/model"
	"git.heroku.com/pg1-go-work/cmd/igoweb/repository"
	"git.heroku.com/pg1-go-work/cmd/igoweb/service"
)

// UserHandler handles request about User
type UserHandler struct {
	Repository repository.UserRepository
	Logger     *service.LoggerService
}

// NewUserHandler instantiate UserHandler instance
func NewUserHandler(rep repository.UserRepository, logger *service.LoggerService) *UserHandler {
	handler := UserHandler{
		Repository: rep,
		Logger:     logger,
	}
	handler.initSuperAdmin()
	return &handler
}

func (handler *UserHandler) initSuperAdmin() {
	superAdmin, err := model.NewUser("superadmin", config.GetInstance().SuperAdminPass)
	if err != nil {
		handler.Logger.Fatal("Failed to create SuperAdmin", err)
		return
	}
	superAdmin.Role = model.RoleAdmin
	err = handler.Repository.Create(superAdmin)
	if err != nil {
		handler.Logger.Fatal("Failed to create SuperAdmin", err)
	}
}
