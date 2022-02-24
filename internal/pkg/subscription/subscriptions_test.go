package subscription

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name          string
		want          *Subscriptions
		subscriptions *Subscriptions
	}{
		{
			name: "one channel, one subreddit",
			want: &Subscriptions{
				ChannelToSubreddits: map[string]map[string]bool{
					"##channel": {"first": true},
				},
				SubredditToChannels: map[string]map[string]bool{
					"first": {"##channel": true},
				},
				Subreddits: []string{"first"},
			},
			subscriptions: &Subscriptions{
				ChannelToSubreddits: map[string]map[string]bool{
					"##channel": {"first": true},
				},
			},
		},
		{
			name: "one channel, two subreddits",
			want: &Subscriptions{
				ChannelToSubreddits: map[string]map[string]bool{
					"##channel": {"first": true, "second": true},
				},
				SubredditToChannels: map[string]map[string]bool{
					"first":  {"##channel": true},
					"second": {"##channel": true},
				},
				Subreddits: []string{"first", "second"},
			},
			subscriptions: &Subscriptions{
				ChannelToSubreddits: map[string]map[string]bool{
					"##channel": {"first": true, "second": true},
				},
			},
		},
		{
			name: "many channels, many subreddits",
			want: &Subscriptions{
				ChannelToSubreddits: map[string]map[string]bool{
					"##channel": {"first": true, "second": true},
					"##other":   {"second": true, "third": true},
				},
				SubredditToChannels: map[string]map[string]bool{
					"first":  {"##channel": true},
					"second": {"##channel": true, "##other": true},
					"third":  {"##other": true},
				},
				Subreddits: []string{"first", "second", "third"},
			},
			subscriptions: &Subscriptions{
				ChannelToSubreddits: map[string]map[string]bool{
					"##channel": {"first": true, "second": true},
					"##other":   {"second": true, "third": true},
				},
			},
		},
	}

	for _, tt := range tests {
		tt.subscriptions.Update()

		if !cmp.Equal(tt.subscriptions, tt.want) {
			t.Logf("got = %v want = %v", tt.subscriptions, tt.want)
			t.Fail()
		}
	}
}

func TestToSubredditString(t *testing.T) {
	tests := []struct {
		name          string
		want          string
		subscriptions *Subscriptions
	}{
		{
			name:          "No subreddits",
			want:          "",
			subscriptions: &Subscriptions{},
		},
		{
			name: "One subreddit",
			want: "subreddit",
			subscriptions: &Subscriptions{
				Subreddits: []string{"subreddit"},
			},
		},
		{
			name: "Many subreddits",
			want: "first+second+third",
			subscriptions: &Subscriptions{
				Subreddits: []string{"first", "second", "third"},
			},
		},
	}

	for _, tt := range tests {
		got := tt.subscriptions.ToSubredditString()

		assert.Equal(t, got, tt.want)
	}
}
