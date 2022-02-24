package subgrok

import (
	"fmt"
	"strings"

	irc "github.com/thoj/go-ircevent"

	"github.com/snoonetIRC/subgrok/internal/pkg/command"
)

const (
	commandSubscribe   = "subscribe"
	commandUnsubscribe = "unsubscribe"

	prefixChannel      = "#"
	prefixHalfOperator = "%"
	prefixOperator     = "@"
)

// events returns a map of IRC events to their handling functions
func events(b *Bot) map[string]func(e *irc.Event) {
	return map[string]func(e *irc.Event){
		"001":     b.callback001,
		"353":     b.callback353,
		"MODE":    b.callbackMode,
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

// callback353 runs when a response to /NAMES is seen
func (b *Bot) callback353(e *irc.Event) {
	if len(e.Arguments) < 2 {
		return
	}

	channel := e.Arguments[2]

	for _, nickname := range strings.Split(e.Message(), " ") {
		if strings.HasPrefix(nickname, prefixOperator) || strings.HasPrefix(nickname, prefixHalfOperator) {
			b.toggleChannelOperator(channel, removeOperatorPrefixes(nickname), true)
		}
	}
}

func (b *Bot) callbackMode(e *irc.Event) {
	if len(e.Arguments) < 3 {
		return
	}

	argumentsContainChannel := func(args []string) bool {
		return strings.HasPrefix(args[0], prefixChannel)
	}

	// 0 = channel, 1 = mode, 2 = nickname
	if e.Arguments[1] == "+o" && argumentsContainChannel(e.Arguments) {
		b.toggleChannelOperator(e.Arguments[0], e.Arguments[2], true)
	}

	if e.Arguments[1] == "-o" && argumentsContainChannel(e.Arguments) {
		b.toggleChannelOperator(e.Arguments[0], e.Arguments[2], false)
	}
}

// callbackPrivmsg runs when the bot receives a message, either in a channel or
// directly from another user
func (b *Bot) callbackPrivmsg(e *irc.Event) {
	channel := e.Arguments[0]

	// Only channel (half) operators who aren't the bot are allowed to run commands
	if e.Nick == b.Connection.GetNick() || !b.isChannelOperator(channel, e.Nick) {
		return
	}

	var response string

	b.Connection.Log.Println(e.Message())

	// Keep it simple for now; we only have two commands (consider dispatch map
	// if the scope grows)
	if b.isCommand(e.Message(), commandSubscribe) {
		response = command.ToggleSubscription(channel, messageToArguments(e.Message()), true, b.Database)
	} else if b.isCommand(e.Message(), commandUnsubscribe) {
		response = command.ToggleSubscription(channel, messageToArguments(e.Message()), false, b.Database)
	}

	if response != "" {
		b.Connection.Privmsgf(e.Arguments[0], response)
	}
}

// isCommand checks the input message starts with the command prefix and a
// command
func (b *Bot) isCommand(message string, command string) bool {
	return strings.HasPrefix(message, fmt.Sprintf("%s%s", b.Config.IRC.CommandPrefix, command))
}

func (b *Bot) toggleChannelOperator(channel string, nickname string, operator bool) {
	if b.Permissions[channel] == nil {
		b.Permissions[channel] = make(map[string]bool)
	}

	b.Permissions[channel][nickname] = operator

	b.Connection.Log.Printf(`Toggled operator. Nick "%s" is operator in "%s": %t.`, nickname, channel, b.Permissions[channel][nickname])
}

func (b *Bot) isChannelOperator(channel string, nickname string) bool {
	if b.Permissions[channel] == nil {
		return false
	}

	return b.Permissions[channel][nickname]
}

// messageToArguments takes a message, removes the command prefix and returns a
// list of the rest of the words in the message (split on space)
func messageToArguments(message string) []string {
	return strings.Split(message, " ")[1:]
}

func removeOperatorPrefixes(nickname string) string {
	nickname = strings.TrimPrefix(nickname, prefixOperator)
	nickname = strings.TrimPrefix(nickname, prefixHalfOperator)

	return nickname
}
