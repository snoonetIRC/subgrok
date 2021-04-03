package main

import (
	"github.com/davecgh/go-spew/spew"
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

	spew.Dump(applicationConfig)
}
