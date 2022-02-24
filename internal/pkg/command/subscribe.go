// The command package contains actions made by issuing a command to the bot.
// This file contains processing for the "subscribe" command.
package command

import (
	"fmt"

	"github.com/snoonetIRC/subgrok/internal/pkg/store"
)

func Subscribe(channel string, messageArguments []string, s *store.FileDB) string {
	subreddit := normaliseSubreddit(messageArguments)

	if subreddit == "" {
		return errNoArguments
	}

	err := s.ToggleSubscription(channel, subreddit, true)

	if err != nil {
		return fmt.Sprintf("Failed to subscribe %s to %s", channel, subreddit)
	}

	return fmt.Sprintf("Subscribed %s to %s", channel, subreddit)
}
