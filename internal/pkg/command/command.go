package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/snoonetIRC/subgrok/internal/pkg/config"
	"github.com/snoonetIRC/subgrok/internal/pkg/store"
)

const (
	errInternal                  = "An internal error occurred"
	errNoArguments               = "A subreddit name is required (e.g. 'irc')"
	errSubredditTooLong          = "Subreddit name is too long (maximum 20 characters)"
	errSubredditValidationFailed = "Subreddit name failed validation"
	errTooManySubscriptions      = "Channel has too many subscriptions (maximum %d)"

	subredditPrefix = "/r/"
)

type SubscribeToggleArguments struct {
	Channel          string
	MessageArguments []string
	Subscribed       bool
	Database         *store.FileDB
	Config           *config.Config
}

func ToggleSubscription(a *SubscribeToggleArguments) string {
	subreddit := normaliseSubreddit(a.MessageArguments)

	if subreddit == "" {
		return errNoArguments
	}

	if len([]rune(subreddit)) > 20 {
		return errSubredditTooLong
	}

	if a.Subscribed {
		maxSubscriptionCount := a.Config.Application.ChannelMaximumSubscriptions
		channelSubscriptionCount, err := a.Database.GetChannelSubscriptionCount(a.Channel)

		if err != nil {
			return errInternal
		}

		if channelSubscriptionCount >= maxSubscriptionCount {
			return fmt.Sprintf(errTooManySubscriptions, maxSubscriptionCount)
		}
	}

	matched, err := regexp.MatchString(`^[A-Za-z0-9-_]+$`, subreddit)

	if err != nil || !matched {
		return errSubredditValidationFailed
	}

	err = a.Database.ToggleSubscription(a.Channel, subreddit, a.Subscribed)

	if err != nil {
		failureMessage := "Failed to subscribe %s to %s"

		if !a.Subscribed {
			failureMessage = "Failed to unsubscribe %s from %s"
		}

		return fmt.Sprintf(failureMessage, a.Channel, subreddit)
	}

	successMessage := "Subscribed %s to %s"

	if !a.Subscribed {
		successMessage = "Unsubscribed %s from %s"
	}

	return fmt.Sprintf(successMessage, a.Channel, subreddit)
}

func ListSubscriptions(a *SubscribeToggleArguments) string {
	subscriptions, err := a.Database.GetChannelSubscriptions(a.Channel)

	if err != nil {
		return errInternal
	}

	return fmt.Sprintf("%s subscribes to: %s", a.Channel, strings.Join(subscriptions, ", "))
}

// normaliseSubreddit converts a subreddit name to the internal standard, which
// doesn't include the "/r/" prefix
func normaliseSubreddit(arguments []string) string {
	if len(arguments) == 0 {
		return ""
	}

	return strings.TrimPrefix(arguments[0], subredditPrefix)
}
