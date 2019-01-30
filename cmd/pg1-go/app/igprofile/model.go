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

const deletedIDCol = "deleted_ig_id"
const igProfileCol = "ig_profile"

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
		modelLogger.Warning("Failed to create index")
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
}

func (model *IgProfile) initTime() {
	model.CreatedAt = time.Now()
	model.ModifiedAt = time.Now()
}

// Save writes IgProfile instance to database
// returns true if success
func Save(igp *IgProfile) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	delCol := dataAccess.GetCollection(deletedIDCol)
	var exsIgp IgProfile
	delCol.Find(bson.M{"ig_id": igp.IGID}).One(&exsIgp)
	if exsIgp.IGID != "" {
		modelLogger.Info(fmt.Sprintf("Failed to create IgProfile because IG ID: %v was banned", igp.IGID))
		return false
	}
	col := dataAccess.GetCollection(igProfileCol)
	igp.initTime()
	err := col.Insert(igp)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to create IgProfile with IG ID: %v", igp.IGID))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to create IgProfile with IG ID: %v", igp.IGID))
	return false
}

// Update modify IgProfile instance in database
// returns true if success
func Update(igID string, changes map[string]interface{}) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igProfileCol)
	selector := bson.M{"ig_id": igID}
	changes["modified_at"] = time.Now()
	update := bson.M{"$set": changes}
	err := col.Update(selector, update)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to update IgProfile with IG ID: %v", igID))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to update IgProfile with IG ID: %v", igID))
	return false
}

// GenerateChanges build hash map of non-empty igp's attributes
func GenerateChanges(igp *IgProfile) map[string]interface{} {
	changes := gin.H{}
	if igp.Name != "" {
		changes["name"] = igp.Name
	}
	if igp.Followers > 0 {
		changes["followers"] = igp.Followers
	}
	if igp.Following > 0 {
		changes["following"] = igp.Following
	}
	if igp.Posts > 0 {
		changes["posts"] = igp.Posts
	}
	if igp.PpURL != "" {
		changes["pp_url"] = igp.PpURL
	}
	return changes
}

// SaveOrUpdate writes IgProfile instance to database
// if IgProfile doesn't exists or update existing IgProfile
// with new data. Returns true if success
func SaveOrUpdate(igp *IgProfile) bool {
	strdIgp := GetIgProfile(igp.IGID)
	if strdIgp.IGID == "" {
		return Save(igp)
	}
	return Update(igp.IGID, GenerateChanges(igp))
}

// GetAll returns All IgProfile in database
// Require offset and limit number for pagination
func GetAll(offset, limit int, sortBy ...string) []IgProfile {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(igProfileCol)
	var igps []IgProfile
	if len(sortBy) == 0 {
		sortBy = []string{"_id"}
	}
	err := col.Find(nil).Sort(sortBy...).Skip(offset).Limit(limit).All(&igps)
	if err == nil {
		modelLogger.Debug("Success to get all IgProfile")
	} else {
		modelLogger.Fatal("Failed to get all IgProfiles")
	}
	return igps
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
func FindIgProfile(query string, offset, limit int, sortBy ...string) []IgProfile {
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
	}).Sort(sortBy...).Skip(offset).Limit(limit).All(&igps)
	if err == nil {
		modelLogger.Debug(fmt.Sprintf("Success to find IgProfile with query: %v", query))
	} else {
		modelLogger.Fatal(fmt.Sprintf("Failed to find IgProfile with query: %v", query))
	}
	return igps
}

// DeleteIgProfile removes IgProfile instance from database by its IGID
// and add the deleted IG ID to another database
// returns true if success
func DeleteIgProfile(igID string) bool {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	ipCol := dataAccess.GetCollection(igProfileCol)
	delCol := dataAccess.GetCollection(deletedIDCol)
	igp := GetIgProfile(igID)
	err := ipCol.RemoveId(igp.ID)
	if err != nil {
		modelLogger.Info(fmt.Sprintf("Failed to delete IgProfile with IG ID: %v", igp.IGID))
		return false
	}
	igp.ID = ""
	delCol.Insert(igp)
	if err == nil {
		modelLogger.Info(fmt.Sprintf("Success to delete IgProfile with IG ID: %v", igp.IGID))
		return true
	}
	modelLogger.Info(fmt.Sprintf("Failed to move deleted IgProfile with IG ID: %v", igp.IGID))
	return false
}

// Builder instantiate the IgProfile using builder pattern
type Builder struct {
	IGID      string
	Name      string
	Followers int
	Following int
	Posts     int
	PpURL     string
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
