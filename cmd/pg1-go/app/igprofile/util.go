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

// FetchIgProfile to fetch Ig Profile information from IG
func FetchIgProfile(igID string) *IgProfile {
	igID = strings.Trim(igID, " ")
	if igID == "" {
		utilLogger.Fatal("IG ID cannot be empty")
		return nil
	}
	r := req.New()
	resp, err := r.Get(fmt.Sprintf("https://www.instagram.com/%s", igID))
	if err == nil {
		code := resp.Response().StatusCode
		if code == http.StatusNotFound {
			utilLogger.Fatal(fmt.Sprintf("IG ID %v not exist", igID))
			return nil
		}
		bodyText := resp.String()
		matches := matcher.FindStringSubmatch(bodyText)
		if len(matches) < 2 {
			utilLogger.Fatal(fmt.Sprintf("Failed to match sharedData on IG ID: %s", igID))
			return nil
		}
		sharedData := matches[1]
		if sharedData == "" {
			utilLogger.Fatal("sharedData is empty")
			return nil
		}
		sharedData = sharedData[:len(sharedData)-1]
		var data map[string]interface{}
		json.Unmarshal([]byte(sharedData), &data)
		if data == nil {
			utilLogger.Fatal(fmt.Sprintf("Failed to parse sharedData on IG ID: %s", igID))
			return nil
		}
		entryData := data["entry_data"].(map[string]interface{})
		pps := entryData["ProfilePage"].([]interface{})
		pp := pps[0].(map[string]interface{})
		graph := pp["graphql"].(map[string]interface{})
		user := graph["user"].(map[string]interface{})
		edgeFollow := user["edge_follow"].(map[string]interface{})
		following := int(edgeFollow["count"].(float64))
		edgeFollowed := user["edge_followed_by"].(map[string]interface{})
		followers := int(edgeFollowed["count"].(float64))
		name := user["full_name"].(string)
		posts := user["edge_owner_to_timeline_media"].(map[string]interface{})
		postsCount := int(posts["count"].(float64))
		ppHD := user["profile_pic_url_hd"].(string)

		builder := NewBuilder()
		builder = builder.SetIGID(igID).SetFollowers(followers).SetName(name)
		builder = builder.SetFollowing(following).SetPosts(postsCount).SetPpURL(ppHD)
		return builder.Build()
	}
	utilLogger.Fatal(fmt.Sprintf("Failed to fetch IG ID: %s", igID))
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
