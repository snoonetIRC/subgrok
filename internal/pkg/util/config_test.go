// util contains general-purpose application utilities.
//
// This file handles tests for the configuration loader.
package util

import (
	"testing"

	"github.com/shibukawa/configdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type loaderMock struct {
	mock.Mock
}

func (c *loaderMock) Open() {
	panic("oh no")
}

func (c *loaderMock) Retrieve() ([]byte, error) {
	panic("oh no")
	return []byte{}, nil
}

func (c *loaderMock) Directory() *configdir.Config {
	panic("oh no")
	return &configdir.Config{}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name       string
		want       *Config
		wantErr    bool
		wantErrMsg string
		loader     ConfigLoader
	}{
		{
			name:       "vendor configuration directory does not contain file",
			want:       &Config{},
			wantErr:    true,
			wantErrMsg: "afilenamethatdoesnotexist does not exist in the configuration directory",
			loader:     NewConfigLoader("avendorthatdoesnotexist", "anapplicationthatdoesnotexist", "afilenamethatdoesnotexist"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.loader)

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && err.Error() != tt.wantErrMsg {
				t.Errorf("Load() error message = %v, wantErr %v", err.Error(), tt.wantErrMsg)
				return
			}

			assert.Equal(t, got, tt.want, "returned config values should match")
		})
	}
}
