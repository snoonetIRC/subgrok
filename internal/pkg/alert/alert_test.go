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
			want: "\x033Link post:\x03" + ` "Title" posted in /r/Subreddit by Author. https://reddit.comURL`,
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					Permalink:     "URL",
					SubredditName: "Subreddit",
					NSFW:          false,
					IsSelfPost:    false,
				},
			},
		},
		{
			name: "NSFW link post",
			want: "\x033Link post:\x03" + ` "Title" posted in /r/Subreddit by Author. https://reddit.comURL ` + "\x02\x034NSFW\x03\x02",
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					Permalink:     "URL",
					SubredditName: "Subreddit",
					NSFW:          true,
					IsSelfPost:    false,
				},
			},
		},
		{
			name: "Non-NSFW self post",
			want: "\x033Self post:\x03" + ` "Title" posted in /r/Subreddit by Author. https://reddit.comURL`,
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					Permalink:     "URL",
					SubredditName: "Subreddit",
					NSFW:          false,
					IsSelfPost:    true,
				},
			},
		},
		{
			name: "NSFW self post",
			want: "\x033Self post:\x03" + ` "Title" posted in /r/Subreddit by Author. https://reddit.comURL ` + "\x02\x034NSFW\x03\x02",
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					Permalink:     "URL",
					SubredditName: "Subreddit",
					NSFW:          true,
					IsSelfPost:    true,
				},
			},
		},
		{
			name: "Post with spoiler",
			want: "\x033Link post:\x03" + ` "Title" posted in /r/Subreddit by Author. https://reddit.comURL ` + "\x02\x037Spoiler\x03\x02",
			alert: &Alert{
				Post: &reddit.Post{
					Author:        "Author",
					Title:         "Title",
					Permalink:     "URL",
					SubredditName: "Subreddit",
					Spoiler:       true,
					IsSelfPost:    false,
				},
			},
		},
	}

	for _, tt := range tests {
		got := tt.alert.ToString()

		assert.Equal(t, got, tt.want)
	}
}
