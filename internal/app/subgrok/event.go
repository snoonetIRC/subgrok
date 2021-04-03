package subgrok

import irc "github.com/thoj/go-ircevent"

// events returns a map of IRC events to their handling functions
func events(b *Bot) map[string]func(e *irc.Event) {
	return map[string]func(e *irc.Event){
		"001":     b.callback001,
		"PRIVMSG": b.callbackPrivmsg,
	}
}

// callback001 runs when the bot connects to an IRC server
func (b *Bot) callback001(e *irc.Event) {
	if b.Config.IRC.Modes != "" {
		b.Connection.Mode(b.Connection.GetNick(), b.Config.IRC.Modes)
	}

	if b.Config.IRC.NickservAccount != "" && b.Config.IRC.NickservPassword != "" {
		b.Connection.Privmsgf("nickserv", "identify %s %s", b.Config.IRC.NickservAccount, b.Config.IRC.NickservPassword)
	}

	// TODO: join stored channels as well as administrative ones
	b.joinChannels(b.Config.IRC.AdminChannels)
}

// callbackPrivmsg runs when the bot receives a message, either in a channel or
// directly from another user
func (b *Bot) callbackPrivmsg(e *irc.Event) {
	// TODO: this will be the entrypoint for bot commands
}
