package subscription

import (
	"sort"
	"strings"

	"github.com/snoonetIRC/subgrok/internal/pkg/store"
)

type Subscriptions struct {
	ChannelToSubreddits map[string]map[string]bool
	SubredditToChannels map[string]map[string]bool
	Subreddits          []string
}

func (s *Subscriptions) ReloadFromDatabase(db *store.FileDB) {
	subscriptions, err := db.GetSubscriptions()

	if err != nil {
		panic(err) // TODO
	}

	s.ChannelToSubreddits = subscriptions
	s.Update()
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

	sort.Strings(subreddits)

	s.Subreddits = subreddits
}

// invert reorders the poller's subscription list to be one subreddit to many
// channels
func (s *Subscriptions) invert() {
	inverted := make(map[string]map[string]bool)

	for channel, subreddits := range s.ChannelToSubreddits {
		for subreddit, subscribed := range subreddits {
			if !subscribed {
				// The IRC channel is not subscribed to receive notifications for
				// this subreddit, so do not add it to the list
				continue
			}

			if inverted[subreddit] == nil {
				inverted[subreddit] = make(map[string]bool)
			}

			inverted[subreddit][channel] = true
		}
	}

	s.SubredditToChannels = inverted
}

// ToSubredditString creates a +-separated string of all subreddit names.
func (s *Subscriptions) ToSubredditString() string {
	return strings.Join(s.Subreddits, "+")
}
