package repository

import (
	"git.heroku.com/pg1-go-work/cmd/igoweb/database"
	"git.heroku.com/pg1-go-work/cmd/igoweb/model"
	"github.com/globalsign/mgo/bson"
)

// IgProfileRepository access IgProfile from repository
type IgProfileRepository interface {
	Create(igProfile *model.IgProfile) error
	Update(igProfile *model.IgProfile) error
	CreateOrUpdate(igProfile *model.IgProfile) error
	Count(findRequest *model.FindRequest, status model.ProfileStatus) (int, error)
	Find(findRequest *model.FindRequest, status model.ProfileStatus) ([]*model.IgProfile, error)
}

// MongoIgProfileRepository is implementation of IgProfileRepository using
// MongoDB as the database
type MongoIgProfileRepository struct {
	DB *database.MongoClient
}

// Create store igProfile to database
func (rep *MongoIgProfileRepository) Create(igProfile *model.IgProfile) error {
	col := rep.DB.Collection(database.IgProfileCollection)
	igProfile.InitTimeStamp()
	return col.Insert(igProfile)
}

// Update modify igProfile data on database
func (rep *MongoIgProfileRepository) Update(igProfile *model.IgProfile) error {
	col := rep.DB.Collection(database.IgProfileCollection)
	selector := bson.M{"ig_id": igProfile.IGID}
	return col.Update(selector, igProfile)
}

// CreateOrUpdate store igProfile to database, if exist, update the data instead
func (rep *MongoIgProfileRepository) CreateOrUpdate(igProfile *model.IgProfile) error {
	col := rep.DB.Collection(database.IgProfileCollection)
	selector := bson.M{"ig_id": igProfile.IGID}
	_, err := col.Upsert(selector, igProfile)
	return err
}

// Count returns number of data which meets FindRequest and ProfileStatus criteria
func (rep *MongoIgProfileRepository) Count(findRequest *model.FindRequest, status model.ProfileStatus) (int, error) {
	col := rep.DB.Collection(database.IgProfileCollection)
	query := generateFindQuery(findRequest, status)
	return col.Find(query).Count()
}

// Find search IgProfile with FindRequest and ProfileStatus criteria
func (rep *MongoIgProfileRepository) Find(findRequest *model.FindRequest, status model.ProfileStatus) ([]*model.IgProfile, error) {
	var igProfiles []*model.IgProfile
	col := rep.DB.Collection(database.IgProfileCollection)
	query := generateFindQuery(findRequest, status)
	q := col.Find(query).Sort(findRequest.Sort).Limit(findRequest.Limit)
	err := q.Skip(findRequest.Offset).All(&igProfiles)
	return igProfiles, err
}

func generateFindQuery(fr *model.FindRequest, status model.ProfileStatus) map[string]interface{} {
	return bson.M{
		"$or": []bson.M{
			bson.M{"ig_id": bson.M{"$regex": fr.Query, "$options": "i"}},
			bson.M{"name": bson.M{"$regex": fr.Query, "$options": "i"}},
		},
		"status": bson.M{"$regex": status, "$options": "i"},
	}
}
