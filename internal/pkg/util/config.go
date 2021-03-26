// util contains general-purpose application utilities.
//
// This file handles configuration, which is stored in a YAML file in the
// system's default configuration directory.
//
// Linux:   ~/.config/snoonet/subgrok/config.yaml
// Windows: %APPDATA%\Roaming\snoonet\subgrok\config.yaml
// macOS:   ${HOME}/Library/Application Support/snoonet/subgrok/config.yaml
package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/shibukawa/configdir"
	"gopkg.in/yaml.v2"
)

// Config is the bot's configuration. It contains strictly IRC-level
// configuration, and should never contain state information such as
// channel subreddit subscription data.
type Config struct {
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

	Hostname              string
	ReconnectDelayMinutes time.Duration
}

type ConfigLoader interface {
	Open()
	Directory() *configdir.Config
	Retrieve() ([]byte, error)
}

type Loader struct {
	Application string
	Filename    string
	Vendor      string

	Processor configdir.ConfigDir
}

func (l *Loader) Open() {
	l.Processor = configdir.New(l.Vendor, l.Application)
}

func (l *Loader) Directory() *configdir.Config {
	dir := l.Processor.QueryFolderContainsFile(l.Filename)
	return dir
}

func (l *Loader) Retrieve() ([]byte, error) {
	var (
		err  error
		data []byte
	)

	if dir := l.Directory(); dir != nil {
		data, err = dir.ReadFile(l.Filename)
	} else {
		err = errors.New(fmt.Sprintf("%s does not exist in the configuration directory", l.Filename))
	}

	return data, err
}

func NewConfigLoader(vendor string, application string, filename string) *Loader {
	loader := &Loader{
		Application: application,
		Filename:    filename,
		Vendor:      vendor,
	}

	loader.Open()

	return loader
}

// Load processes the bot's configuration
func Load(l ConfigLoader) (*Config, error) {
	config := &Config{}

	data, err := l.Retrieve()

	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, config)

	if err != nil {
		return config, err
	}

	return config, nil
}
