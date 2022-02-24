// subgrok is an IRC bot which monitors new posts to subreddits. It receives
// messages from subpoll when a new post is created, and sends alerts to the
// channels which are subscribed to the subreddit in question.
//
// This file contains setup for the IRC bot's connections to IRC networks.
package subgrok

import (
	"crypto/tls"

	irc "github.com/thoj/go-ircevent"

	"github.com/snoonetIRC/subgrok/internal/pkg/alert"
	"github.com/snoonetIRC/subgrok/internal/pkg/config"
	"github.com/snoonetIRC/subgrok/internal/pkg/store"
)

// Bot is a SubGrok instance
type Bot struct {
	Config     *config.Config
	Connection *irc.Connection
	Database   *store.FileDB

	// Permissions contain users who are an operator for the channels the bot is
	// in ("##mike"["Mike"]true)
	Permissions map[string]map[string]bool
}

// Load configures the IRC bot when it is first launched
func Load(config *config.Config) *Bot {
	connection := irc.IRC(config.IRC.Nickname, config.IRC.Ident)

	applyConfigToConnection(connection, config)

	bot := &Bot{Connection: connection, Config: config}
	bot.Permissions = make(map[string]map[string]bool)

	for name, fn := range events(bot) {
		bot.Connection.AddCallback(name, fn)
	}

	return bot
}

// applyConfigToConnection sets the bot's configuration against its connection
// to the IRC network
func applyConfigToConnection(connection *irc.Connection, config *config.Config) {
	connection.Debug = config.IRC.Debug
	connection.VerboseCallbackHandler = config.IRC.Debug
	connection.RealName = config.IRC.RealName

	if config.IRC.UseTLS {
		connection.UseTLS = true
		connection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

// Connect connects the bot to an IRC network
func (b *Bot) Connect() {
	if err := b.Connection.Connect(b.Config.IRC.Hostname()); err != nil {
		panic(err)
	}

	// TODO: start healthcheck

	b.Connection.Loop()
}

// joinChannels attempts to join a provided list of IRC channels
func (b *Bot) joinChannels(channels []string) {
	for _, channel := range channels {
		b.Connection.Join(channel)
	}
}

// AlertChannel sends a formatted "new post" announcement to a given channel.
func (b *Bot) AlertChannel(channel string, alert *alert.Alert) {
	b.Connection.Privmsg(channel, alert.ToString())
}
