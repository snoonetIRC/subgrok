// The config package loads and parses the application's YAML configuration.
//
// The configuration can be loaded from:
// Linux:   ~/.config/snoonet/subgrok/config.yaml
// Windows: %APPDATA%\Roaming\snoonet\subgrok\config.yaml
// macOS:   ${HOME}/Library/Application Support/snoonet/subgrok/config.yaml
package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
	"gopkg.in/yaml.v2"
)

const (
	defaultChannelMaximumSubscriptions = 20
	defaultCommandPrefix               = "!"
	defaultRedditPollWaitTime          = 60
)

type databaseConfig struct {
	Filepath string `yaml:"filepath"`
}

// ircConfig contains the bot's IRC connection configuration.
type ircConfig struct {
	AdminChannels  []string `yaml:"admin_channels"`
	CommandPrefix  string   `yaml:"command_prefix"`
	Debug          bool     `yaml:"debug"` // Debug also handles "verbose" mode
	Ident          string   `yaml:"ident"`
	Modes          string   `yaml:"modes"`
	Nickname       string   `yaml:"nickname"`
	Port           int      `yaml:"port"`
	RealName       string   `yaml:"real_name"`
	Server         string   `yaml:"server"`
	ServerPassword string   `yaml:"server_password"`
	UseTLS         bool     `yaml:"use_tls"`

	NickservAccount  string `yaml:"nickserv_account"`
	NickservPassword string `yaml:"nickserv_password"`

	MaxReconnect   int `yaml:"max_reconnect"`
	ReconnectDelay int `yaml:"reconnect_delay"`

	ReconnectDelayMinutes time.Duration
}

// redditConfig contains the bot's reddit authentication configuration.
type redditConfig struct {
	ID           string `yaml:"id"`
	Secret       string `yaml:"secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	PollWaitTime int    `yaml:"poll_wait_time"`

	PollWaitDuration time.Duration
}

type applicationConfig struct {
	ChannelMaximumSubscriptions int `yaml:"channel_maximum_subscriptions"`
}

// Config is the bot's configuration. It contains strictly IRC connection and
// reddit auth configuration, and should never contain state information such as
// channel subreddit subscription data.
type Config struct {
	Application *applicationConfig `yaml:"application"`
	Database    *databaseConfig    `yaml:"database"`
	IRC         *ircConfig         `yaml:"irc"`
	Reddit      *redditConfig      `yaml:"reddit"`
}

// Load processes the bot's configuration
func Load() (*Config, error) {
	config := &Config{}

	data, err := ProcessFile("snoonet", "subgrok", "config.yaml")

	if err != nil {
		return config, err
	}

	err = yaml.UnmarshalStrict(data, config)

	if err != nil {
		return config, err
	}

	if config.Database == nil || config.Database.Filepath == "" {
		return config, errors.New("a database filepath is required (config database.filepath)")
	}

	config.applyDefaults()

	return config, nil
}

func (c *Config) applyDefaults() {
	if c.Reddit == nil {
		c.Reddit = &redditConfig{}
	}

	c.Reddit.applyDefaults()

	if c.IRC == nil {
		c.IRC = &ircConfig{}
	}

	if c.Application == nil {
		c.Application = &applicationConfig{}
	}

	c.Application.applyDefaults()

	if c.IRC.CommandPrefix == "" {
		c.IRC.CommandPrefix = defaultCommandPrefix
	}
}

func (ic *ircConfig) Hostname() string {
	return fmt.Sprintf("%s:%d", ic.Server, ic.Port)
}

func (rc *redditConfig) applyDefaults() {
	if rc.PollWaitTime < defaultRedditPollWaitTime {
		rc.PollWaitTime = defaultRedditPollWaitTime
	}

	rc.PollWaitDuration = time.Duration(rc.PollWaitTime) * time.Second
}

func (rc *redditConfig) Credentials() *reddit.Credentials {
	return &reddit.Credentials{
		ID:       rc.ID,
		Secret:   rc.Secret,
		Username: rc.Username,
		Password: rc.Password,
	}
}

func (ac *applicationConfig) applyDefaults() {
	if ac.ChannelMaximumSubscriptions == 0 {
		ac.ChannelMaximumSubscriptions = defaultChannelMaximumSubscriptions
	}
}
