package config

import (
	"testing"
	"time"

	"github.com/shibukawa/configdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const exampleConfig = `
---
database:
  filepath: whatever

irc:
  admin_channels:
    - '##subgrok'
  debug: true
  ident: subgrokinstance
  modes: +B
  nickname: subgrokinstance
  port: 6697
  real_name: SubGrok
  server: irc.snoonet.org
  use_tls: true

  nickserv_account: someaccount
  nickserv_password: somepassword

reddit:
  id: 1234
  secret: secret
  username: username
  password: password
  poll_wait_duration: 60s
`

const exampleBadConfig = `
---
nope: it's bad
`

const exampleConfigNoDatabaseFilepath = ``

const exampleBadRedditWaitTime = `
---
database:
  filepath: whatever

irc:
  admin_channels:
    - '##subgrok'
  debug: true
  ident: subgrokinstance
  modes: +B
  nickname: subgrokinstance
  port: 6697
  real_name: SubGrok
  server: irc.snoonet.org
  use_tls: true

  nickserv_account: someaccount
  nickserv_password: somepassword

reddit:
  id: 1234
  secret: secret
  username: username
  password: password
  poll_wait_duration: 59s
`

var exampleValidConfig = &Config{
	Database: &databaseConfig{Filepath: "whatever"},
	IRC: &ircConfig{
		AdminChannels:    []string{"##subgrok"},
		CommandPrefix:    "!",
		Debug:            true,
		Ident:            "subgrokinstance",
		Modes:            "+B",
		Nickname:         "subgrokinstance",
		Port:             6697,
		RealName:         "SubGrok",
		Server:           "irc.snoonet.org",
		UseTLS:           true,
		NickservAccount:  "someaccount",
		NickservPassword: "somepassword",
	},
	Reddit: &redditConfig{
		ID:               "1234",
		Secret:           "secret",
		Username:         "username",
		Password:         "password",
		PollWaitDuration: time.Duration(60 * time.Second),
	},
	Application: &applicationConfig{ChannelMaximumSubscriptions: defaultChannelMaximumSubscriptions},
}

type processorMockValidContent struct{ mock.Mock }

func (c *processorMockValidContent) Open(l *Loader) *configdir.ConfigDir {
	return &configdir.ConfigDir{}
}

func (c *processorMockValidContent) Retrieve(l *Loader) ([]byte, error) {
	return []byte(exampleConfig), nil
}

func (c *processorMockValidContent) Directory(l *Loader, f *configdir.ConfigDir) *configdir.Config {
	return &configdir.Config{}
}

type processorMockInvalidContent struct{ mock.Mock }

func (c *processorMockInvalidContent) Open(l *Loader) *configdir.ConfigDir {
	return &configdir.ConfigDir{}
}

func (c *processorMockInvalidContent) Retrieve(l *Loader) ([]byte, error) {
	return []byte(exampleBadConfig), nil
}

func (c *processorMockInvalidContent) Directory(l *Loader, f *configdir.ConfigDir) *configdir.Config {
	return &configdir.Config{}
}

type processorMockInvalidDatabaseFilepath struct{ mock.Mock }

func (c *processorMockInvalidDatabaseFilepath) Open(l *Loader) *configdir.ConfigDir {
	return &configdir.ConfigDir{}
}

func (c *processorMockInvalidDatabaseFilepath) Retrieve(l *Loader) ([]byte, error) {
	return []byte(exampleConfigNoDatabaseFilepath), nil
}

func (c *processorMockInvalidDatabaseFilepath) Directory(l *Loader, f *configdir.ConfigDir) *configdir.Config {
	return &configdir.Config{}
}

type processorMockInvalidRedditWaitTime struct{ mock.Mock }

func (c *processorMockInvalidRedditWaitTime) Open(l *Loader) *configdir.ConfigDir {
	return &configdir.ConfigDir{}
}

func (c *processorMockInvalidRedditWaitTime) Retrieve(l *Loader) ([]byte, error) {
	return []byte(exampleBadRedditWaitTime), nil
}

func (c *processorMockInvalidRedditWaitTime) Directory(l *Loader, f *configdir.ConfigDir) *configdir.Config {
	return &configdir.Config{}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name          string
		want          *Config
		wantErr       bool
		wantErrMsg    string
		processorMock ConfigProcessor
	}{
		{
			name:          "With valid configuration",
			want:          exampleValidConfig,
			processorMock: &processorMockValidContent{},
		},
		{
			name:          "With invalid configuration",
			want:          &Config{},
			wantErr:       true,
			wantErrMsg:    "yaml: unmarshal errors:\n  line 3: field nope not found in type config.Config",
			processorMock: &processorMockInvalidContent{},
		},
		{
			name:          "File processing failure",
			want:          &Config{},
			wantErr:       true,
			wantErrMsg:    "error message here",
			processorMock: &processorMockErrorRaised{}, // from file_processor_test.go
		},
		{
			name:          "Missing database filepath",
			want:          &Config{},
			wantErr:       true,
			wantErrMsg:    "a database filepath is required (config database.filepath)",
			processorMock: &processorMockInvalidDatabaseFilepath{},
		},
		{
			name:          "reddit PollWaitDuration below minimum",
			want:          exampleValidConfig,
			processorMock: &processorMockInvalidRedditWaitTime{},
		},
	}

	for _, tt := range tests {
		processor = tt.processorMock

		got, err := Load()

		if (err != nil) != tt.wantErr {
			t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		if (err != nil) && err.Error() != tt.wantErrMsg {
			t.Errorf("Load() error message = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
			return
		}

		assert.Equal(t, got, tt.want, "returned config should match ("+tt.name+")")
	}
}
