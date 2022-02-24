package command

import "strings"

const (
	errNoArguments = "A subreddit name is required (e.g. 'irc')"

	subredditPrefix = "/r/"
)

// normaliseSubreddit converts a subreddit name to the internal standard, which
// doesn't include the "/r/" prefix
func normaliseSubreddit(arguments []string) string {
	if len(arguments) == 0 {
		return ""
	}

	return strings.TrimPrefix(arguments[0], subredditPrefix)
}
