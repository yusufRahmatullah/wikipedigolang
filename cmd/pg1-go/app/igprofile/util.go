package igprofile

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igmedia"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"github.com/imroc/req"
)

const (
	// GraphSideCar is TypeName of ShortcodeMedia that represents multiple media in one post
	GraphSideCar = "GraphSideCar"
	// GraphImage is TypeName of ShortcodeMedia that represents one image in one post
	GraphImage = "GraphImage"
)

var (
	matcher    *regexp.Regexp
	accMatcher *regexp.Regexp
	utilLogger = logger.NewLogger("IgProfileUtil", false, true)
)

func init() {
	matcher = regexp.MustCompile(`<script type="text/javascript">\s?window._sharedData\s?=\s?([^<>]*)</script>`)
	accMatcher = regexp.MustCompile(`@[\w_]+`)
}

// IgData root of ig data
type IgData struct {
	EntryData HasProfilePage `json:"entry_data"`
}

// HasProfilePage node with key ProfilePage
type HasProfilePage struct {
	ProfilePage []HasGraphql `json:"ProfilePage"`
	PostPage    []HasGraphql `json:"PostPage"`
}

// HasGraphql node with key graphql
type HasGraphql struct {
	Graphql HasUser `json:"graphql"`
}

// HasUser node with key user
type HasUser struct {
	User           UserData       `json:"user"`
	ShortcodeMedia ShortcodeMedia `json:"shortcode_media"`
}

// ShortcodeMedia node with key shortcode_media
type ShortcodeMedia struct {
	Caption              MediaData `json:"edge_media_to_caption"`
	TypeName             string    `json:"__typename"`
	Media                MediaData `json:"edge_sidecar_to_children"`
	Tagged               MediaData `json:"edge_media_to_tagged_user"`
	ID                   string    `json:"id"`
	DisplayURL           string    `json:"display_url"`
	AccessibilityCaption string    `json:"accessibility_caption"`
	Text                 string    `json:"text"`
}

// UserData node with many keys
type UserData struct {
	EdgeFollow     HasCount  `json:"edge_follow"`
	EdgeFollowedBy HasCount  `json:"edge_followed_by"`
	FullName       string    `json:"full_name"`
	EdgeOwnerMedia MediaData `json:"edge_owner_to_timeline_media"`
	ProfPic        string    `json:"profile_pic_url_hd"`
	IsPrivate      bool      `json:"is_private"`
}

// MediaData node holds list of edges
type MediaData struct {
	Count int       `json:"count"`
	Edges []HasNode `json:"edges"`
}

// HasNode node with key node
type HasNode struct {
	Node NodeData `json:"node"`
}

// NodeData holds IgProfile's media data
type NodeData struct {
	ID                   string      `json:"id"`
	DisplayURL           string      `json:"display_url"`
	AccessibilityCaption string      `json:"accessibility_caption"`
	Text                 string      `json:"text"`
	User                 HasUsername `json:"user"`
}

// HasUsername node with key username
type HasUsername struct {
	Username string `json:"username"`
}

// HasCount node with key count
type HasCount struct {
	Count int `json:"count"`
}

func getDataFromResponse(resp *req.Resp) (*IgData, string) {
	code := resp.Response().StatusCode
	if code == http.StatusNotFound {
		return nil, "Post ID not exist"
	}
	bodyText := resp.String()
	matches := matcher.FindStringSubmatch(bodyText)
	if len(matches) < 2 {
		return nil, "Failed to match sharedData"
	}
	sharedData := matches[1]
	if sharedData == "" {
		return nil, "sharedData is empty"
	}
	var data IgData
	sharedData = sharedData[:len(sharedData)-1]
	err := json.Unmarshal([]byte(sharedData), &data)
	if err != nil {
		return nil, "Failed to parse shared data"
	}
	return &data, ""
}

func processGraphSideCar(sc ShortcodeMedia, igID string) []*igmedia.IgMedia {
	var retVals []*igmedia.IgMedia
	for _, media := range sc.Media.Edges {
		node := media.Node
		caption := node.AccessibilityCaption
		if strings.Contains(caption, "people") || strings.Contains(caption, "person") {
			igm := igmedia.NewIgMedia(node.ID, igID, node.DisplayURL)
			retVals = append(retVals, igm)
		}
	}
	return retVals
}

func processGraphImage(sc ShortcodeMedia, igID string) []*igmedia.IgMedia {
	var retVals []*igmedia.IgMedia
	if strings.Contains(sc.AccessibilityCaption, "people") || strings.Contains(sc.AccessibilityCaption, "person") {
		igm := igmedia.NewIgMedia(sc.ID, igID, sc.DisplayURL)
		retVals = append(retVals, igm)
	}
	return retVals
}

func getMediasFromData(data *IgData, igID string) []*igmedia.IgMedia {
	var retVals []*igmedia.IgMedia
	sc := data.EntryData.PostPage[0].Graphql.ShortcodeMedia
	if sc.TypeName == GraphSideCar {
		retVals = processGraphSideCar(sc, igID)
	} else if sc.TypeName == GraphImage {
		retVals = processGraphImage(sc, igID)
	}
	return retVals
}

// FetchMediaFromPost fetch Ig Media from post of IG ID
// Returns Ig Medias and empty string if success
// otherwise returns empty array and error message
func FetchMediaFromPost(igID string, postID string) ([]*igmedia.IgMedia, string) {
	var retVals []*igmedia.IgMedia
	postID = strings.Trim(postID, " ")
	if postID == "" {
		return retVals, "Post ID is empty"
	}
	r := req.New()
	resp, err := r.Get(fmt.Sprintf("https://www.instagram.com%s", postID))
	if err == nil {
		data, errStr := getDataFromResponse(resp)
		if errStr != "" {
			utilLogger.Fatal("Failed to get data from response", errors.New(errStr))
			return retVals, errStr
		}
		retVals = getMediasFromData(data, igID)
		return retVals, ""
	}
	return retVals, err.Error()
}

func getIDsFromData(data *IgData) []string {
	var retVals []string
	sc := data.EntryData.PostPage[0].Graphql.ShortcodeMedia
	edges := sc.Caption.Edges
	for _, edge := range edges {
		accs := accMatcher.FindAllString(edge.Node.Text, -1)
		retVals = append(retVals, accs...)
	}
	if len(sc.Tagged.Edges) > 0 {
		for _, edge := range sc.Tagged.Edges {
			retVals = append(retVals, edge.Node.User.Username)
		}
	}
	return retVals
}

// FetchAccountFromPost fetch IG IDs from post
// Returns IG IDs and empty string if success
// otherwise returns empty array and error message
func FetchAccountFromPost(postID string) ([]string, string) {
	var retVals []string
	postID = strings.Trim(postID, " ")
	if postID == "" {
		return retVals, "Post ID is empty"
	}
	r := req.New()
	resp, err := r.Get(fmt.Sprintf("https://www.instagram.com%s", postID))
	if err == nil {
		data, errStr := getDataFromResponse(resp)
		if errStr != "" {
			utilLogger.Fatal("Failed to get data from response", errors.New(errStr))
			return retVals, errStr
		}
		retVals = getIDsFromData(data)
		return retVals, ""
	}
	return retVals, err.Error()
}

// FetchIgProfile to fetch Ig Profile information from IG
func FetchIgProfile(igID string) *IgProfile {
	igID = strings.Trim(igID, " ")
	if igID == "" {
		utilLogger.Fatal("IG ID cannot be empty", nil)
		return nil
	}
	r := req.New()
	resp, err := r.Get(fmt.Sprintf("https://www.instagram.com/%s", igID))
	if err == nil {
		data, errStr := getDataFromResponse(resp)
		if errStr != "" {
			utilLogger.Fatal("Failed to get data from response", errors.New(errStr))
			return nil
		}
		user := data.EntryData.ProfilePage[0].Graphql.User
		following := user.EdgeFollow.Count
		followers := user.EdgeFollowedBy.Count
		name := user.FullName
		if name == "" {
			name = "@" + igID
		}
		postsCount := user.EdgeOwnerMedia.Count
		ppHD := user.ProfPic

		builder := NewBuilder()
		builder = builder.SetIGID(igID).SetFollowers(followers).SetName(name)
		builder = builder.SetFollowing(following).SetPosts(postsCount).SetPpURL(ppHD)
		return builder.Build()
	}
	utilLogger.Fatal(fmt.Sprintf("Failed to fetch IG ID: %s", igID), err)
	return nil
}

// CleanIgIDParams clean igID params from JobQueue which may copied
// from complete URL
func CleanIgIDParams(igID string) string {
	splts := strings.Split(igID, "/")
	splts = strings.Split(splts[len(splts)-1], "?")
	cleanID := splts[0]
	return cleanID
}

// FindByName find IG ID by name
// returns ig id if exist, otherwise empty string
func FindByName(name string) string {
	cname := strings.Trim(name, " ")
	r := req.New()
	params := req.Param{
		"query": cname,
	}
	resp, err := r.Get("https://www.instagram.com/web/search/topsearch", params)
	if err == nil {
		code := resp.Response().StatusCode
		if code != http.StatusOK {
			utilLogger.Fatal(fmt.Sprintf("Failed to find by name: '%v' with code: %v", cname, code), nil)
			return ""
		}
		bodyText := resp.String()
		var data map[string]interface{}
		json.Unmarshal([]byte(bodyText), &data)
		if data == nil {
			utilLogger.Fatal(fmt.Sprintf("Failed to unmarshall response on name: '%v, bodyText: %v'", cname, bodyText), nil)
			return ""
		}
		users := data["users"].([]interface{})
		if len(users) == 0 {
			utilLogger.Fatal(fmt.Sprintf("Failed, users is empty on name: '%v'", cname), nil)
			return ""
		}
		topUser := users[0].(map[string]interface{})
		userData := topUser["user"].(map[string]interface{})
		return userData["username"].(string)
	}
	utilLogger.Fatal(fmt.Sprintf("Failed to find by name: '%v'", cname), err)
	return ""
}
