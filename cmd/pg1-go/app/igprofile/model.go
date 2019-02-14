package igprofile

import (
	"fmt"
	"time"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"
	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// ProfileStatus is a status of IgProfile
type ProfileStatus string

const (
	igProfileCol = "ig_profile"
	// StatusActive means the IgProfile will be shown
	StatusActive ProfileStatus = "active"
	// StatusBanned means the IgProfile will not be shown
	StatusBanned ProfileStatus = "banned"
	// StatusMulti means the IgProfile will be shown on MultiAcc page
	// as active Multi Account
	StatusMulti ProfileStatus = "multi"
	// StatusBannedMulti means IgProfile will be shown on MultiAcc page
	// as inactive Multi Account
	StatusBannedMulti ProfileStatus = "banned_multi"
	// StatusAll means all IgProfile will be shown
	StatusAll ProfileStatus = ""
)

var modelLogger = logger.NewLogger("IGProfile", true, true)

func init() {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igProfileCol)
	idIndex := mgo.Index{
		Key:        []string{"ig_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := col.EnsureIndex(idIndex)
	if err != nil {
		modelLogger.Warning(fmt.Sprintf("Failed to create index on %v", igProfileCol))
	}
}

// IgProfile holds information about IG Profile
// include its IG ID, Name, followers number, following number,
// post number, and profile picture URL
type IgProfile struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CreatedAt  time.Time     `json:"created_at,omitempty" bson:"created_at"`
	ModifiedAt time.Time     `json:"modified_at,omitempty" bson:"modified_at"`
	IGID       string        `json:"ig_id" bson:"ig_id"`
	Name       string        `json:"name" bson:"name"`
	Followers  int           `json:"followers" bson:"followers"`
	Following  int           `json:"following" bson:"following"`
	Posts      int           `json:"posts" bson:"posts"`
	PpURL      string        `json:"pp_url" bson:"pp_url"`
	Status     ProfileStatus `josn:"status" bson:"status"`
}

func (model *IgProfile) initTime() {
	model.CreatedAt = time.Now()
	model.ModifiedAt = time.Now()
}

// Save writes IgProfile instance to database
// returns empty string if success
func Save(igp *IgProfile) string {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igProfileCol)
	igp.initTime()
	err := col.Insert(igp)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to create IgProfile with IG ID: %v", igp.IGID))
		return ""
	}
	modelLogger.Fatal(fmt.Sprintf("Failed to create IgProfile with IG ID: %v", igp.IGID), err)
	return "Failed to create IgProfile"
}

// Update modify IgProfile instance in database
// returns empty string if success
func Update(igID string, changes map[string]interface{}) string {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igProfileCol)
	selector := bson.M{"ig_id": igID}
	changes["modified_at"] = time.Now()
	update := bson.M{"$set": changes}
	err := col.Update(selector, update)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to update IgProfile with IG ID: %v", igID))
		return ""
	}
	modelLogger.Fatal(fmt.Sprintf("Failed to update IgProfile with IG ID: %v", igID), err)
	return "Failed to update IgProfile"
}

type changeBuilder struct {
	igp       *IgProfile
	name      string
	followers int
	following int
	posts     int
	ppURL     int
	status    ProfileStatus
	changes   gin.H
}

func (cb *changeBuilder) setIgp(igp *IgProfile) {
	cb.igp = igp
	cb.changes = gin.H{}
}

func (cb *changeBuilder) checkName() {
	if cb.igp.Name != "" {
		cb.changes["name"] = cb.igp.Name
	}
}

func (cb *changeBuilder) checkFollowers() {
	if cb.igp.Followers > 0 {
		cb.changes["followers"] = cb.igp.Followers
	}
}

func (cb *changeBuilder) checkFollowing() {
	if cb.igp.Following > 0 {
		cb.changes["following"] = cb.igp.Following
	}
}

func (cb *changeBuilder) checkPostss() {
	if cb.igp.Posts > 0 {
		cb.changes["posts"] = cb.igp.Posts
	}
}

func (cb *changeBuilder) checkURL() {
	if cb.igp.PpURL != "" {
		cb.changes["pp_url"] = cb.igp.PpURL
	}
}

func (cb *changeBuilder) checkStatus() {
	if cb.igp.Status != "" {
		cb.changes["status"] = cb.igp.Status
	}
}

func (cb *changeBuilder) getChanges() gin.H {
	return cb.changes
}

// GenerateChanges build hash map of non-empty igp's attributes
func GenerateChanges(igp *IgProfile) map[string]interface{} {
	var cb changeBuilder
	cb.setIgp(igp)
	cb.checkName()
	cb.checkFollowers()
	cb.checkFollowing()
	cb.checkURL()
	cb.checkStatus()
	return cb.getChanges()
}

// SaveOrUpdate writes IgProfile instance to database
// if IgProfile doesn't exists or update existing IgProfile
// with new data. Returns empty string if success
func SaveOrUpdate(igp *IgProfile) string {
	strdIgp := GetIgProfile(igp.IGID)
	if strdIgp.IGID == "" {
		return Save(igp)
	}
	return Update(igp.IGID, GenerateChanges(igp))
}

// GetAll returns All IgProfile in database
// Require offset and limit number for pagination
// Require status to define what status of the Profile
func GetAll(offset, limit int, status ProfileStatus, sortBy ...string) []IgProfile {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igProfileCol)
	var igps []IgProfile
	if len(sortBy) == 0 {
		sortBy = []string{"_id"}
	}
	err := col.Find(bson.M{
		"status": bson.M{"$regex": status, "$options": "i"},
	}).Sort(sortBy...).Skip(offset).Limit(limit).All(&igps)
	if err == nil {
		modelLogger.Debug("Success to get all IgProfile")
	} else {
		modelLogger.Fatal("Failed to get all IgProfiles", err)
	}
	return igps
}

func countIgProfiles(status ProfileStatus) (int, error) {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igProfileCol)
	return col.Find(bson.M{"status": bson.M{"$regex": status}}).Count()
}

// GetIgProfile get IgProfile instance in database by its IGID
func GetIgProfile(igID string) *IgProfile {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igProfileCol)
	igp := IgProfile{}
	err := col.Find(bson.M{"ig_id": igID}).One(&igp)
	if err == nil {
		modelLogger.Debug(fmt.Sprintf("Success to get IgProfile with IG ID: %v", igp.IGID))
	} else {
		modelLogger.Debug(fmt.Sprintf("Failed to get IgProfile with IG ID: %v", igp.IGID))
	}
	return &igp
}

// FindIgProfile find IgProfiles in database by its IGID or name
// Require offset and limit number for pagination
// Require status to define
func FindIgProfile(query string, offset, limit int, status ProfileStatus, sortBy ...string) []IgProfile {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igProfileCol)
	var igps []IgProfile
	if len(sortBy) == 0 {
		sortBy = []string{"-modified_at"}
	}
	err := col.Find(bson.M{
		"$or": []bson.M{
			bson.M{"ig_id": bson.M{"$regex": query, "$options": "i"}},
			bson.M{"name": bson.M{"$regex": query, "$options": "i"}},
		},
		"status": bson.M{"$regex": status, "$options": "i"},
	}).Sort(sortBy...).Skip(offset).Limit(limit).All(&igps)
	if err == nil {
		modelLogger.Debug(fmt.Sprintf("Success to find IgProfile with query: %v", query))
	} else {
		modelLogger.Fatal(fmt.Sprintf("Failed to find IgProfile with query: %v", query), err)
	}
	return igps
}

// DeleteIgProfile removes IgProfile instance from database by its IGID
// and add the deleted IG ID to another database
// returns empty string if success
func DeleteIgProfile(igID string, isMulti bool) string {
	status := StatusBanned
	if isMulti {
		status = StatusBannedMulti
	}
	return Update(igID, bson.M{"status": status})
}

// Builder instantiate the IgProfile using builder pattern
type Builder struct {
	IGID      string
	Name      string
	Followers int
	Following int
	Posts     int
	PpURL     string
	Status    ProfileStatus
}

// NewBuilder instante new IgProgile Builder
func NewBuilder() *Builder {
	return &Builder{
		IGID:      "",
		Name:      "",
		Followers: 0,
		Following: 0,
		Posts:     0,
		PpURL:     "",
		Status:    StatusActive,
	}
}

// Build instantiate new IgProfile instance with Builder's attribute
func (bd *Builder) Build() *IgProfile {
	return &IgProfile{
		IGID:      bd.IGID,
		Name:      bd.Name,
		Followers: bd.Followers,
		Following: bd.Following,
		Posts:     bd.Posts,
		PpURL:     bd.PpURL,
		Status:    bd.Status,
	}
}

// SetFollowers set Builder's Followers
func (bd *Builder) SetFollowers(fol int) *Builder {
	bd.Followers = fol
	return bd
}

// SetFollowing set Builder's Following
func (bd *Builder) SetFollowing(fol int) *Builder {
	bd.Following = fol
	return bd
}

// SetIGID set Builder's IGID
func (bd *Builder) SetIGID(igID string) *Builder {
	bd.IGID = igID
	return bd
}

// SetName set Builder's Name
func (bd *Builder) SetName(name string) *Builder {
	bd.Name = name
	return bd
}

// SetPosts set Builder's Posts
func (bd *Builder) SetPosts(posts int) *Builder {
	bd.Posts = posts
	return bd
}

// SetPpURL set Builder's PpURL
func (bd *Builder) SetPpURL(ppURL string) *Builder {
	bd.PpURL = ppURL
	return bd
}

// SetStatus set Builder's Status
func (bd *Builder) SetStatus(status ProfileStatus) *Builder {
	bd.Status = status
	return bd
}
