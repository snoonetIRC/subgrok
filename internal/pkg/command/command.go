package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/snoonetIRC/subgrok/internal/pkg/store"
)

const (
	errNoArguments               = "A subreddit name is required (e.g. 'irc')"
	errSubredditTooLong          = "Subreddit name is too long (maximum 20 characters)";
	errSubredditValidationFailed = "Subreddit name failed validation"

	subredditPrefix = "/r/"
)

func ToggleSubscription(channel string, messageArguments []string, subscribed bool, s *store.FileDB) string {
	subreddit := normaliseSubreddit(messageArguments)

	if subreddit == "" {
		return errNoArguments
	}

	if len([]rune(subreddit)) > 20 {
		return errSubredditTooLong
	}

	matched, err := regexp.MatchString(`^[A-Za-z0-9-_]+$`, subreddit)

	if err != nil || !matched {
		return errSubredditValidationFailed
	}

	err = s.ToggleSubscription(channel, subreddit, subscribed)

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
