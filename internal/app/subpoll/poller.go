// subpoll watches subreddits for new posts, pushing new messages to subgrok
// for distribution via IRC.
package subpoll

import "internal/app/subgrok"

// Poller watches subreddits for new posts, pushing messages via its Bot member
type Poller struct {
	Bot *subgrok.Bot // The IRC bot which will receive messages on post creation
}

// Alerts are pushed to the IRC bot when a new post is made
type Alert struct {
	SubReddit string
	PostTitle string
	PostURL   string
}

// Load builds a new Poller
func Load() *Poller {
	return &Poller{}
}

// Poll looks at subreddits for new posts. If a new post is found, a message is
// pushed to the Bot.
func (p *Poller) Poll() {}
