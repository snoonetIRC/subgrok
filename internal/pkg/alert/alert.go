package alert

import (
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Alerts are pushed to the IRC bot when a new post is made
type Alert struct {
	Channels []string
	Post     *reddit.Post
}

// ToString formats the alert as a string.
func (a *Alert) ToString() string {
	alert := fmt.Sprintf(`"%s" posted in /r/%s by %s. %s`, a.Post.Title, a.Post.SubredditName, a.Post.Author, a.Post.URL)

	if a.Post.NSFW {
		alert = alert + " \x0304NSFW"
	}

	return alert
}
