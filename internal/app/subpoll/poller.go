// subpoll watches subreddits for new posts, pushing new messages to subgrok
// for distribution via IRC.
package subpoll

import (
	"context"
	"strings"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/snoonetIRC/subgrok/internal/app/subgrok"
	"github.com/snoonetIRC/subgrok/internal/pkg/alert"
	"github.com/snoonetIRC/subgrok/internal/pkg/config"
	"github.com/snoonetIRC/subgrok/internal/pkg/subscription"
)

//import "github.com/davecgh/go-spew/spew"

// Poller watches subreddits for new posts, pushing messages via its Bot member
type Poller struct {
	API           *reddit.Client
	Bot           *subgrok.Bot // The IRC bot which will receive messages on post creation
	Config        *config.Config
	Subscriptions *subscription.Subscriptions
	LastPoll      *time.Time
}

// Load builds a new Poller
func Load(config *config.Config, bot *subgrok.Bot) *Poller {
	//client, err := reddit.NewClient(*config.Reddit.Credentials())
	client, err := reddit.NewReadonlyClient()

	if err != nil {
		panic(err)
	}

	poller := &Poller{
		API:    client,
		Bot:    bot,
		Config: config,

		Subscriptions: &subscription.Subscriptions{
			ChannelToSubreddits: map[string]map[string]bool{
				"##Mike": {"metal": true, "homelab": true, "funny": true, "pics": true},
			},
		},
	}

	poller.Subscriptions.Update()

	return poller
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
			alerts, errs := p.checkSubscriptions()

			p.updateSubscriptions()

			for _, alert := range alerts {
				for channel := range alert.Channels {
					p.Bot.AlertChannel(channel, alert)
				}
			}

			if len(errs) > 0 {
				var errorStrings []string

				for _, err := range errs {
					errorStrings = append(errorStrings, err.Error())
				}

				p.Bot.Connection.Log.Printf("Received errors from reddit: %s", strings.Join(errorStrings, ", "))
			}

			//spew.Dump(alerts)

			time.Sleep(p.Config.Reddit.PollWaitDuration)
		}
	}()
}

func (p *Poller) checkSubscriptions() ([]*alert.Alert, []error) {
	var (
		redditErrors []error
		alerts       []*alert.Alert
	)

	posts, _, err := p.API.Subreddit.NewPosts(context.Background(), p.Subscriptions.ToSubredditString(), &reddit.ListOptions{
		Limit: 5 * len(p.Subscriptions.Subreddits),
	})

	if err != nil {
		redditErrors = append(redditErrors, err)
	}

	for _, post := range posts {
		if p.LastPoll == nil || !post.Created.After(*p.LastPoll) {
			continue // Skip posts which were created before the last poll time
		}

		if err != nil {
			redditErrors = append(redditErrors, err)
			continue
		}

		alerts = append(alerts, &alert.Alert{
			Channels: p.Subscriptions.SubredditToChannels[post.SubredditName],
			Post:     post,
		})
	}

	p.setLastPollTime()

	return alerts, redditErrors
}

func (p *Poller) updateSubscriptions() {
	// TODO: retrieve subscription changes from bot cache, set them, clean up
	p.Bot.Database.GetSubscriptions()
}

// setLastPollTime sets the last time the poller ran. Posts retrieved by
// go-reddit have a UTC Created time, so the poller also uses UTC.
func (p *Poller) setLastPollTime() {
	lastPoll := time.Now().UTC()
	p.LastPoll = &lastPoll
}
