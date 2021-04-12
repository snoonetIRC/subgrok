package subscription

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {}

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
