// subpoll watches subreddits for new posts, pushing new messages to subgrok
// for distribution via IRC.
package subpoll

import (
	"context"
	"strings"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/n7st/subgrok/internal/app/subgrok"
	"github.com/n7st/subgrok/internal/pkg/config"
)

import "github.com/davecgh/go-spew/spew"

type Subscriptions struct {
	ChannelToSubreddits map[string][]string
	SubredditToChannels map[string][]string
	Subreddits          []string
}

// Poller watches subreddits for new posts, pushing messages via its Bot member
type Poller struct {
	API           *reddit.Client
	Bot           *subgrok.Bot // The IRC bot which will receive messages on post creation
	Config        *config.Config
	Subscriptions *Subscriptions
	LastPoll      *time.Time
}

// Alerts are pushed to the IRC bot when a new post is made
type Alert struct {
	Channels  []string
	SubReddit string
	PostTitle string
	PostURL   string
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

		Subscriptions: &Subscriptions{
			ChannelToSubreddits: map[string][]string{
				"##Mike": {"metal", "homelab", "funny", "pics"},
			},
		},
	}

	poller.Subscriptions.Update()

	return poller
}

// Update reorders subscriptions so duplicate API calls are not made
func (s *Subscriptions) Update() {
	s.invert()
	s.createList()
}

// createList formats the subscribed subreddits from every channel as a string
// slice
func (s *Subscriptions) createList() {
	var subreddits []string

	for subreddit := range s.SubredditToChannels {
		subreddits = append(subreddits, subreddit)
	}

	s.Subreddits = subreddits
}

// invert reorders the poller's subscription list to be one subreddit to many
// channels
func (s *Subscriptions) invert() {
	inverted := make(map[string][]string)

	for channel, subreddits := range s.ChannelToSubreddits {
		for _, subreddit := range subreddits {
			if inverted[subreddit] == nil {
				inverted[subreddit] = []string{}
			}

			inverted[subreddit] = append(inverted[subreddit], channel)
		}
	}

	s.SubredditToChannels = inverted
}

func (s *Subscriptions) ToSubredditString() string {
	return strings.Join(s.Subreddits, "+")
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

			for _, alert := range alerts {
				for _, channel := range alert.Channels {
					p.Bot.Connection.Privmsgf(channel, "%s %s %s", alert.SubReddit, alert.PostTitle, alert.PostURL)
				}
			}

			if len(errs) > 0 {
				spew.Dump(errs)
			}

			spew.Dump(alerts)

			time.Sleep(p.Config.Reddit.PollWaitDuration)
		}
	}()
}

func (p *Poller) checkSubscriptions() ([]*Alert, []error) {
	var (
		errors []error
		alerts []*Alert
	)

	posts, _, err := p.API.Subreddit.NewPosts(context.Background(), p.Subscriptions.ToSubredditString(), &reddit.ListOptions{
		Limit: 5 * len(p.Subscriptions.Subreddits),
	})

	if err != nil {
		errors = append(errors, err)
	}

	for _, post := range posts {
		if p.LastPoll == nil || !post.Created.After(*p.LastPoll) {
			continue // Skip posts which were created before the last poll time
		}

		if err != nil {
			errors = append(errors, err)
			continue
		}

		alerts = append(alerts, &Alert{
			Channels:  p.Subscriptions.SubredditToChannels[post.SubredditName],
			SubReddit: post.SubredditName,
			PostTitle: post.Title,
			PostURL:   post.URL,
		})
	}

	p.setLastPollTime()

	return alerts, errors
}

// setLastPollTime sets the last time the poller ran. Posts retrieved by
// go-reddit have a UTC Created time, so the poller also uses UTC.
func (p *Poller) setLastPollTime() {
	lastPoll := time.Now().UTC()
	p.LastPoll = &lastPoll
}
