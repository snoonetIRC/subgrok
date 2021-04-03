package main

import (
	"github.com/n7st/subgrok/internal/app/subgrok"
	"github.com/n7st/subgrok/internal/app/subpoll"
	"github.com/n7st/subgrok/internal/pkg/config"
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
