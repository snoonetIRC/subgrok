package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func TestToString(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		alert *Alert
	}{
		{
			name: "Non-NSFW",
			want: `"Title" posted in /r/Subreddit by Author. URL`,
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					URL:           "URL",
					SubredditName: "Subreddit",
					NSFW:          false,
				},
			},
		},
		{
			name: "NSFW",
			want: `"Title" posted in /r/Subreddit by Author. URL ` + "\x0304NSFW",
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					URL:           "URL",
					SubredditName: "Subreddit",
					NSFW:          true,
				},
			},
		},
	}

	for _, tt := range tests {
		got := tt.alert.ToString()

		assert.Equal(t, got, tt.want)
	}
}
