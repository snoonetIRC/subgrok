package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		alert *Alert
	}{
		{
			name: "Non-NSFW",
			want: `"Title" posted in /r/SubReddit by Author. URL`,
			alert: &Alert{
				Author:    "Author",
				PostTitle: "Title",
				PostURL:   "URL",
				SubReddit: "SubReddit",
				NSFW:      false,
			},
		},
		{
			name: "NSFW",
			want: `"Title" posted in /r/SubReddit by Author. URL ` + "\x0304NSFW",
			alert: &Alert{
				Author:    "Author",
				PostTitle: "Title",
				PostURL:   "URL",
				SubReddit: "SubReddit",
				NSFW:      true,
			},
		},
	}

	for _, tt := range tests {
		got := tt.alert.ToString()

		assert.Equal(t, got, tt.want)
	}
}
