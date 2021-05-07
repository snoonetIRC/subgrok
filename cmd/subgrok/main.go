package main

import (
	"github.com/snoonetIRC/subgrok/internal/app/subgrok"
	"github.com/snoonetIRC/subgrok/internal/app/subpoll"
	"github.com/snoonetIRC/subgrok/internal/pkg/config"
)

const (
	vendor      = "snoonet"
	application = "subgrok"
	filename    = "config.yaml"
)

func main() {
	applicationConfig, err := config.Load()

	if err != nil {
		panic(err)
	}

	bot := subgrok.Load(applicationConfig)
	poller := subpoll.Load(applicationConfig, bot)

	poller.Poll()
	bot.Connect()
}
