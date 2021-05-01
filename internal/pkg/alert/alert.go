package alert

import (
	"fmt"

	ircColor "github.com/n7st/go-ircformat/color"
	ircEmphasis "github.com/n7st/go-ircformat/emphasis"
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
	prefix := ircColor.Green(fmt.Sprintf("%s post:", a.postType()))
	alert := fmt.Sprintf(`%s "%s" posted in /r/%s by %s. %s`,
		prefix, a.Post.Title, a.Post.SubredditName, a.Post.Author, a.Post.URL)

	if a.Post.NSFW {
		alert = alert + " " + ircEmphasis.Bold(ircColor.Red("NSFW"))
	}

	return alert
}
