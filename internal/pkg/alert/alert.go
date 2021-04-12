package alert

import "fmt"

// Alerts are pushed to the IRC bot when a new post is made
type Alert struct {
	Author    string
	Channels  []string
	PostTitle string
	PostURL   string
	SubReddit string
	NSFW      bool
}

// ToString formats the alert as a string.
func (a *Alert) ToString() string {
	alert := fmt.Sprintf(`"%s" posted in /r/%s by %s. %s`, a.PostTitle, a.SubReddit, a.Author, a.PostURL)

	if a.NSFW {
		alert = alert + " \x0304NSFW"
	}

	return alert
}
