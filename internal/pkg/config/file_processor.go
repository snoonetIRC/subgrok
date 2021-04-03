// The config package loads and parses the application's YAML configuration.
package config

import (
	"errors"
	"fmt"

	"github.com/shibukawa/configdir"
)

const errNoSuchFileFmt = "%s does not exist in the configuration directory"

type Loader struct {
	Application string
	Filename    string
	Vendor      string
}

// ConfigProcessor is a wrapper for the parts of "configdir" we use.
type ConfigProcessor interface {
	Open(*Loader) *configdir.ConfigDir
	Directory(*Loader, *configdir.ConfigDir) *configdir.Config
	Retrieve(*Loader) ([]byte, error)
}

type Processor struct{}

var processor ConfigProcessor

// Open opens the system's configuration directories.
func (p *Processor) Open(l *Loader) *configdir.ConfigDir {
	cd := configdir.New(l.Vendor, l.Application)

	return &cd
}

// Directory finds the correct (if any) directory for the application's
// configuration. This checks various locations on the filesystem (e.g. local
// and system level AppData).
func (p *Processor) Directory(l *Loader, f *configdir.ConfigDir) *configdir.Config {
	return f.QueryFolderContainsFile(l.Filename)
}

// Retrieve attempts to retrieve data from the config file if one exists.
func (p *Processor) Retrieve(l *Loader) ([]byte, error) {
	var (
		err  error
		data []byte
	)

	if dir := p.Directory(l, p.Open(l)); dir != nil {
		data, err = dir.ReadFile(l.Filename)
	} else {
		err = errors.New(fmt.Sprintf(errNoSuchFileFmt, l.Filename))
	}

	return data, err
}

// init sets the Processor used by ProcessFile. This can be overridden in tests.
func init() {
	processor = &Processor{}
}

// ProcessFile attempts to load data from config files.
func ProcessFile(vendor string, application string, filename string) ([]byte, error) {
	loader := &Loader{
		Application: application,
		Filename:    filename,
		Vendor:      vendor,
	}

	return processor.Retrieve(loader)
}
