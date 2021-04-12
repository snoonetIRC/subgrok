package alert

import (
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

const (
	postTypeSelf = "Self"
	postTypeLink = "Link"
)

// Alerts are pushed to the IRC bot when a new post is made
type Alert struct {
	Channels []string
	Post     *reddit.Post
}

func (a *Alert) postType() string {
	var postType string

	if a.Post.IsSelfPost {
		postType = postTypeSelf
	} else {
		postType = postTypeLink
	}

	return postType
}

// ToString formats the alert as a string.
func (a *Alert) ToString() string {
	prefix := fmt.Sprintf("\x0303%s post:\x03", a.postType())
	alert := fmt.Sprintf(`%s "%s" posted in /r/%s by %s. %s`, prefix, a.Post.Title, a.Post.SubredditName, a.Post.Author, a.Post.URL)

	if a.Post.NSFW {
		alert = alert + " \x0304NSFW"
	}

	return alert
}
