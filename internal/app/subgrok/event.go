package subgrok

import irc "github.com/thoj/go-ircevent"

// events returns a map of IRC events to their handling functions
func events(b *Bot) map[string]func(e *irc.Event) {
	return map[string]func(e *irc.Event){
		"001": b.callback001,
	}
}

// callback001 runs when the bot connects to an IRC server
func (b *Bot) callback001(e *irc.Event) {}
