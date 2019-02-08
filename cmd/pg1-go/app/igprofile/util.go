package igprofile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
	"github.com/imroc/req"
)

var (
	matcher    *regexp.Regexp
	utilLogger = logger.NewLogger("IgProfileUtil", false, true)
)

func init() {
	matcher = regexp.MustCompile(`<script type="text/javascript">window._sharedData = ([^<>]*)</script>`)
}

// IgData root of ig data
type IgData struct {
	EntryData HasProfilePage `json:"entry_data"`
}

// HasProfilePage node with key ProfilePage
type HasProfilePage struct {
	ProfilePage []HasGraphql `json:"ProfilePage"`
}

// HasGraphql node with key graphql
type HasGraphql struct {
	Graphql HasUser `json:"graphql"`
}

// HasUser node with key user
type HasUser struct {
	User UserData `json:"user"`
}

// UserData node with many keys
type UserData struct {
	EdgeFollow     HasCount  `json:"edge_follow"`
	EdgeFollowedBy HasCount  `json:"edge_followed_by"`
	FullName       string    `json:"full_name"`
	EdgeOwnerMedia MediaData `json:"edge_owner_to_timeline_media"`
	ProfPic        string    `json:"profile_pic_url_hd"`
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
}

// HasCount node with key count
type HasCount struct {
	Count int `json:"count"`
}

// TopTwelveMedia to get top 12 media's URL of IgProfile
// with acessibility caption contains people
func TopTwelveMedia(igID string) []NodeData {
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
		var retVals []NodeData
		edges := user.EdgeOwnerMedia.Edges
		for _, edge := range edges {
			retVals = append(retVals, edge.Node)
		}
		return retVals
	}
	utilLogger.Fatal(fmt.Sprintf("Failed to fetch IG ID: %s", igID), err)
	return nil
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
