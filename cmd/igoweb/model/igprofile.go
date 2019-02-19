package model

// ProfileStatus is a status of IgProfile
type ProfileStatus string

const (
	// StatusActive means the IgProfile will be shown
	StatusActive ProfileStatus = "active"
	// StatusBanned means the IgProfile will not be shown
	StatusBanned ProfileStatus = "banned"
	// StatusMulti means the IgProfile will be shown on MultiAcc page
	// as active Multi Account
	StatusMulti ProfileStatus = "multi"
	// StatusAll means all IgProfile will be shown
	StatusAll ProfileStatus = ""
)

// IgProfile holds information about IG Profile
// include its IG ID, Name, followers number, following number,
// post number, and profile picture URL
type IgProfile struct {
	TimeStamp
	MongoID
	IGID      string        `json:"ig_id" bson:"ig_id"`
	Name      string        `json:"name" bson:"name"`
	Followers int           `json:"followers" bson:"followers"`
	Following int           `json:"following" bson:"following"`
	Posts     int           `json:"posts" bson:"posts"`
	PpURL     string        `json:"pp_url" bson:"pp_url"`
	Status    ProfileStatus `json:"status" bson:"status"`
}

// IgProfileBuilder instantiate IgProfile using Builder pattern
type IgProfileBuilder struct {
	IGID      string
	Name      string
	Followers int
	Following int
	Posts     int
	PpURL     string
	Status    ProfileStatus
}

// NewIgProfileBuilder instantiate new IgProfileBuilder
func NewIgProfileBuilder() *IgProfileBuilder {
	return &IgProfileBuilder{
		IGID:      "",
		Name:      "",
		Followers: 0,
		Following: 0,
		Posts:     0,
		PpURL:     "",
		Status:    StatusActive,
	}
}

// Build instantiate IgProfile instance with builder's attributes
func (bd *IgProfileBuilder) Build() *IgProfile {
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

// SetIGID set builder's IGID
func (bd *IgProfileBuilder) SetIGID(igID string) *IgProfileBuilder {
	bd.IGID = igID
	return bd
}

// SetName set builder's Name
func (bd *IgProfileBuilder) SetName(name string) *IgProfileBuilder {
	bd.Name = name
	return bd
}

// SetFollowers set builder's Followers
func (bd *IgProfileBuilder) SetFollowers(followers int) *IgProfileBuilder {
	bd.Followers = followers
	return bd
}

// SetFollowing set builder's Following
func (bd *IgProfileBuilder) SetFollowing(following int) *IgProfileBuilder {
	bd.Following = following
	return bd
}

// SetPosts set builder's Posts
func (bd *IgProfileBuilder) SetPosts(posts int) *IgProfileBuilder {
	bd.Posts = posts
	return bd
}

// SetPpURL set builder's PpURL
func (bd *IgProfileBuilder) SetPpURL(ppURL string) *IgProfileBuilder {
	bd.PpURL = ppURL
	return bd
}

// SetStatus set builder's Status
func (bd *IgProfileBuilder) SetStatus(status ProfileStatus) *IgProfileBuilder {
	bd.Status = status
	return bd
}
