// subgrok is an IRC bot which monitors new posts to subreddits. It receives
// messages from subpoll when a new post is created, and sends alerts to the
// channels which are subscribed to the subreddit in question.
//
// This file contains setup for the IRC bot's connections to IRC networks.
package subgrok

// Bot is a SubGrok instance
type Bot struct{}

// Load configures the IRC bot when it is first launched
func Load() *Bot {}

// Connect connects the bot to an IRC network
func (b *Bot) Connect() {}
