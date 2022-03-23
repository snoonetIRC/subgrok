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

// Poller watches subreddits for new posts, pushing messages via its Bot member
type Poller struct {
	API           *reddit.Client
	Bot           *subgrok.Bot // The IRC bot which will receive messages on post creation
	Config        *config.Config
	Subscriptions *subscription.Subscriptions
	LastPoll      *time.Time
	TooRecent     map[string]bool
}

// Load builds a new Poller
func Load(config *config.Config, bot *subgrok.Bot) *Poller {
	client, err := reddit.NewClient(*config.Reddit.Credentials())

	if err != nil {
		panic(err)
	}

	poller := &Poller{
		API:    client,
		Bot:    bot,
		Config: config,

		Subscriptions: &subscription.Subscriptions{},
	}

	poller.TooRecent = make(map[string]bool)

	poller.Subscriptions.ReloadFromDatabase(bot.Database)

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

			p.Subscriptions.ReloadFromDatabase(p.Bot.Database)

			for _, alert := range alerts {
				for channel := range alert.Channels {
					p.Bot.AlertChannel(channel, alert)
					time.Sleep(5 * time.Second)
				}
			}

			if len(errs) > 0 {
				var errorStrings []string

				for _, err := range errs {
					// Strip console-spamming "badger badger badger" message from
					// reddit's rate-limiting error page
					e := err.Error()
					e = strings.ReplaceAll(e, "badger ", "")

					errorStrings = append(errorStrings, e)
				}

				p.Bot.Connection.Log.Printf("Received errors from reddit: %s", strings.Join(errorStrings, ", "))
			}

			time.Sleep(p.Config.Reddit.PollWaitDuration)
		}
	}()
}

func (p *Poller) checkSubscriptions() ([]*alert.Alert, []error) {
	var (
		redditErrors []error
		alerts       []*alert.Alert
	)

	posts, client, err := p.API.Subreddit.NewPosts(context.Background(), p.Subscriptions.ToSubredditString(), &reddit.ListOptions{
		Limit: 10 * len(p.Subscriptions.Subreddits),
	})

	if err != nil {
		redditErrors = append(redditErrors, err)
	} else {
		p.Bot.Connection.Log.Printf("API requests: used %d, remaining %d (resets %s)",
			client.Rate.Used, client.Rate.Remaining, client.Rate.Reset.String())
	}

	tooRecent := make(map[string]bool)

	for _, post := range posts {
		// Skip posts which were created before the last poll time unless they
		// were deemed to be too new to display
		if p.LastPoll == nil || (!p.TooRecent[post.ID] && !post.Created.After(*p.LastPoll)) {
			continue
		}

		if !p.postIsOldEnough(post) {
			tooRecent[post.ID] = true

			continue // Skip posts that aren't yet old enough
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

	// Clean up post delay cache - some of the old items could be removed by
	// moderators and end up staying in there forever
	p.TooRecent = tooRecent

	p.setLastPollTime()

	return alerts, redditErrors
}

// setLastPollTime sets the last time the poller ran. Posts retrieved by
// go-reddit have a UTC Created time, so the poller also uses UTC.
func (p *Poller) setLastPollTime() {
	lastPoll := time.Now().UTC()
	p.LastPoll = &lastPoll
}

func (p *Poller) postIsOldEnough(post *reddit.Post) bool {
	var minimumPostAgeSeconds = p.Config.Reddit.MinimumPostAge

	if minimumPostAgeSeconds == 0 {
		return true
	}

	oldEnough := time.Since(post.Created.Time) > p.Config.Reddit.MinimumPostAge

	if !oldEnough && p.Config.IRC.Debug {
		p.Bot.Connection.Log.Printf("Post isn't old enough: %s (%s) (Minimum %.4f seconds)\n",
			post.Title, post.Created.String(), p.Config.Reddit.MinimumPostAge.Seconds())
	}

	return oldEnough
}
