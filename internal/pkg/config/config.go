// The config package loads and parses the application's YAML configuration.
//
// The configuration can be loaded from:
// Linux:   ~/.config/snoonet/subgrok/config.yaml
// Windows: %APPDATA%\Roaming\snoonet\subgrok\config.yaml
// macOS:   ${HOME}/Library/Application Support/snoonet/subgrok/config.yaml
package config

import (
	"fmt"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
	"gopkg.in/yaml.v2"
)

// ircConfig contains the bot's IRC connection configuration.
type ircConfig struct {
	AdminChannels  []string `yaml:"admin_channels"`
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
	ID       string `yaml:"id"`
	Secret   string `yaml:"secret"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Config is the bot's configuration. It contains strictly IRC connection and
// reddit auth configuration, and should never contain state information such as
// channel subreddit subscription data.
type Config struct {
	IRC    *ircConfig    `yaml:"irc"`
	Reddit *redditConfig `yaml:"reddit"`
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

	return config, nil
}

func (ic *ircConfig) Hostname() string {
	return fmt.Sprintf("%s:%d", ic.Server, ic.Port)
}

func (rc *redditConfig) Credentials() *reddit.Credentials {
	return &reddit.Credentials{
		ID:       rc.ID,
		Secret:   rc.Secret,
		Username: rc.Username,
		Password: rc.Password,
	}
}
