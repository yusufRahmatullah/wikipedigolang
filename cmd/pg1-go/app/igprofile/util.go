package igprofile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/igmedia"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"github.com/imroc/req"
)

const (
	GraphSideCar = "GraphSideCar"
	GraphImage   = "GraphImage"
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
	PostPage    []HasGraphql `json:"PostPafe`
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
	Caption  MediaData `json:"edge_media_to_caption"`
	TypeName string    `json:"__typename"`
	Media    MediaData `json:"edge_sidecar_to_children"`
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
	ID                   string `json:"id"`
	DisplayURL           string `json:"display_url"`
	AccessibilityCaption string `json:"accessibility_caption"`
	Text                 string `json:"text"`
}

// HasCount node with key count
type HasCount struct {
	Count int `json:"count"`
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
		code := resp.Response().StatusCode
		if code == http.StatusNotFound {
			return retVals, "Post ID not exist"
		}
		bodyText := resp.String()
		matches := matcher.FindStringSubmatch(bodyText)
		if len(matches) < 2 {
			utilLogger.Fatal(fmt.Sprintf("Failed to match sharedData on Post ID: %s", postID), err)
			return retVals, "Failed to match sharedData"
		}
		sharedData := matches[1]
		if sharedData == "" {
			return retVals, "sharedData is empty"
		}
		sharedData = sharedData[:len(sharedData)-1]
		var data IgData
		err := json.Unmarshal([]byte(sharedData), &data)
		if err != nil {
			utilLogger.Fatal(fmt.Sprintf("Failed to parse sharedData on Post ID: %s", postID), err)
			return retVals, "Failed to parse shared data"
		}
		sc := data.EntryData.PostPage[0].Graphql.ShortcodeMedia
		if sc.TypeName == GraphSideCar {
			for _, media := range sc.Media.Edges {
				node := media.Node
				caption := node.AccessibilityCaption
				if strings.Contains(caption, "people") || strings.Contains(caption, "person") {
					igm := igmedia.NewIgMedia(node.ID, igID, node.DisplayURL)
					retVals = append(retVals, igm)
				}
			}
		} else if sc.TypeName == GraphImage {

		}
		return retVals, ""
	}
	return retVals, err.Error()
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
		code := resp.Response().StatusCode
		if code == http.StatusNotFound {
			return retVals, "Post ID not exist"
		}
		bodyText := resp.String()
		matches := matcher.FindStringSubmatch(bodyText)
		if len(matches) < 2 {
			utilLogger.Fatal(fmt.Sprintf("Failed to match sharedData on Post ID: %s", postID), err)
			return retVals, "Failed to match sharedData"
		}
		sharedData := matches[1]
		if sharedData == "" {
			return retVals, "sharedData is empty"
		}
		sharedData = sharedData[:len(sharedData)-1]
		var data IgData
		err := json.Unmarshal([]byte(sharedData), &data)
		if err != nil {
			utilLogger.Fatal(fmt.Sprintf("Failed to parse sharedData on Post ID: %s", postID), err)
			return retVals, "Failed to parse shared data"
		}
		sc := data.EntryData.PostPage[0].Graphql.ShortcodeMedia
		edges := sc.Caption.Edges
		for _, edge := range edges {
			accs := accMatcher.FindAllString(edge.Node.Text, -1)
			retVals = append(retVals, accs...)
		}
		return retVals, ""
	}
	return retVals, err.Error()
}

// TopTwelveMedia to get top 12 media's URL of IgProfile
// with acessibility caption contains people
// Return array of NodeData and error message if error occurred
func TopTwelveMedia(igID string) ([]NodeData, string) {
	igID = strings.Trim(igID, " ")
	if igID == "" {
		return nil, "IG ID is empty"
	}
	r := req.New()
	resp, err := r.Get(fmt.Sprintf("https://www.instagram.com/%s", igID))
	if err == nil {
		code := resp.Response().StatusCode
		if code == http.StatusNotFound {
			return nil, "IG ID not exist"
		}
		bodyText := resp.String()
		matches := matcher.FindStringSubmatch(bodyText)
		if len(matches) < 2 {
			utilLogger.Fatal(fmt.Sprintf("Failed to match sharedData on IG ID: %s", igID), err)
			return nil, "Failed to match sharedData"
		}
		sharedData := matches[1]
		if sharedData == "" {
			return nil, "sharedData is empty"
		}
		sharedData = sharedData[:len(sharedData)-1]
		var data IgData
		err := json.Unmarshal([]byte(sharedData), &data)
		if err != nil {
			utilLogger.Fatal(fmt.Sprintf("Failed to parse sharedData on IG ID: %s", igID), err)
			return nil, "Failed to parse shared data"
		}
		user := data.EntryData.ProfilePage[0].Graphql.User
		if user.IsPrivate {
			return nil, "IG is private"
		}
		var retVals []NodeData
		edges := user.EdgeOwnerMedia.Edges
		for _, edge := range edges {
			retVals = append(retVals, edge.Node)
		}
		return retVals, ""
	}
	return nil, err.Error()
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
		code := resp.Response().StatusCode
		if code == http.StatusNotFound {
			utilLogger.Fatal(fmt.Sprintf("IG ID %v not exist", igID), err)
			return nil
		}
		bodyText := resp.String()
		matches := matcher.FindStringSubmatch(bodyText)
		if len(matches) < 2 {
			utilLogger.Fatal(fmt.Sprintf("Failed to match sharedData on IG ID: %s", igID), err)
			return nil
		}
		sharedData := matches[1]
		if sharedData == "" {
			utilLogger.Fatal("sharedData is empty", nil)
			return nil
		}
		sharedData = sharedData[:len(sharedData)-1]
		var data IgData
		err := json.Unmarshal([]byte(sharedData), &data)
		if err != nil {
			utilLogger.Fatal(fmt.Sprintf("Failed to parse sharedData on IG ID: %s", igID), err)
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
