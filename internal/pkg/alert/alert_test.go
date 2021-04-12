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
			name: "Non-NSFW link post",
			want: "\x0303Link post:\x03" + ` "Title" posted in /r/Subreddit by Author. URL`,
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					URL:           "URL",
					SubredditName: "Subreddit",
					NSFW:          false,
					IsSelfPost:    false,
				},
			},
		},
		{
			name: "NSFW link post",
			want: "\x0303Link post:\x03" + ` "Title" posted in /r/Subreddit by Author. URL ` + "\x0304NSFW",
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					URL:           "URL",
					SubredditName: "Subreddit",
					NSFW:          true,
					IsSelfPost:    false,
				},
			},
		},
		{
			name: "Non-NSFW self post",
			want: "\x0303Self post:\x03" + ` "Title" posted in /r/Subreddit by Author. URL`,
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					URL:           "URL",
					SubredditName: "Subreddit",
					NSFW:          false,
					IsSelfPost:    true,
				},
			},
		},
		{
			name: "NSFW self post",
			want: "\x0303Self post:\x03" + ` "Title" posted in /r/Subreddit by Author. URL ` + "\x0304NSFW",
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					URL:           "URL",
					SubredditName: "Subreddit",
					NSFW:          true,
					IsSelfPost:    true,
				},
			},
		},
	}

	for _, tt := range tests {
		got := tt.alert.ToString()

		assert.Equal(t, got, tt.want)
	}
}
