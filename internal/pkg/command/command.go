package command

import (
	"fmt"
	"strings"

	"github.com/snoonetIRC/subgrok/internal/pkg/store"
)

const (
	errNoArguments = "A subreddit name is required (e.g. 'irc')"

	subredditPrefix = "/r/"
)

func ToggleSubscription(channel string, messageArguments []string, subscribed bool, s *store.FileDB) string {
	subreddit := normaliseSubreddit(messageArguments)

	if subreddit == "" {
		return errNoArguments
	}

	err := s.ToggleSubscription(channel, subreddit, subscribed)

	if err != nil {
		failureMessage := "Failed to subscribe %s to %s"

		if !subscribed {
			failureMessage = "Failed to unsubscribe %s from %s"
		}

		return fmt.Sprintf(failureMessage, channel, subreddit)
	}

	successMessage := "Subscribed %s to %s"

	if !subscribed {
		successMessage = "Unsubscribed %s from %s"
	}

	return fmt.Sprintf(successMessage, channel, subreddit)
}

// normaliseSubreddit converts a subreddit name to the internal standard, which
// doesn't include the "/r/" prefix
func normaliseSubreddit(arguments []string) string {
	if len(arguments) == 0 {
		return ""
	}

	return strings.TrimPrefix(arguments[0], subredditPrefix)
}
