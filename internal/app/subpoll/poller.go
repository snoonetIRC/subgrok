// subpoll watches subreddits for new posts, pushing new messages to subgrok
// for distribution via IRC.
package subpoll

import (
	"fmt"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/n7st/subgrok/internal/app/subgrok"
	"github.com/n7st/subgrok/internal/pkg/config"
)

// Poller watches subreddits for new posts, pushing messages via its Bot member
type Poller struct {
	API *reddit.Client
	Bot *subgrok.Bot // The IRC bot which will receive messages on post creation
}

// Alerts are pushed to the IRC bot when a new post is made
type Alert struct {
	SubReddit string
	PostTitle string
	PostURL   string
}

// Load builds a new Poller
func Load(config *config.Config, bot *subgrok.Bot) *Poller {
	//client, err := reddit.NewReadonlyClient(*config.Reddit.Credentials())
	client, err := reddit.NewReadonlyClient()

	if err != nil {
		panic(err)
	}

	return &Poller{
		API: client,
		Bot: bot,
	}
}

// Poll looks at subreddits for new posts. If a new post is found, a message is
// pushed to the Bot.
func (p *Poller) Poll() {
	for {
		if p.Bot.Connection.Connected() {
			time.Sleep(1 * time.Second)
			break
		} else {
			time.Sleep(10 * time.Second)
		}
	}

	go func() {
		for {
			fmt.Println("Tick")
			time.Sleep(10 * time.Second)
		}
	}()
}
