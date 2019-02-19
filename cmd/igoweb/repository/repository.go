package repository

import (
	"git.heroku.com/pg1-go-work/cmd/igoweb/model"
)

// IgProfileRepository access IgProfile from repository
type IgProfileRepository interface {
	Create(igProfile *model.IgProfile) error
	Find(findRequest *model.FindRequest, status model.ProfileStatus) ([]*model.IgProfile, error)
	Activate(igProfile *model.IgProfile) error
	Ban(igProfile *model.IgProfile) error
	AsMulti(igProfile *model.IgProfile) error
}
