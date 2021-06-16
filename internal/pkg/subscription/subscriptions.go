package subscription

import (
	"sort"
	"strings"
)

type Subscriptions struct {
	ChannelToSubreddits map[string][]string
	SubredditToChannels map[string][]string
	Subreddits          []string
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

// ToSubredditString creates a +-separated string of all subreddit names.
func (s *Subscriptions) ToSubredditString() string {
	return strings.Join(s.Subreddits, "+")
}
