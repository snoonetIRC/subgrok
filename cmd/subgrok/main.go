package main

import (
	"github.com/snoonetIRC/subgrok/internal/app/subgrok"
	"github.com/snoonetIRC/subgrok/internal/app/subpoll"
	"github.com/snoonetIRC/subgrok/internal/pkg/config"
	"github.com/snoonetIRC/subgrok/internal/pkg/store"
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

	database, err := store.NewStore(applicationConfig.GetBoltDB())

	if err != nil {
		panic(err)
	}

	bot := subgrok.Load(applicationConfig)

	bot.Database = database

	poller := subpoll.Load(applicationConfig, bot)

	poller.Poll()
	bot.Connect()
}
