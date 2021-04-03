package config

import (
	"errors"
	"testing"

	"github.com/shibukawa/configdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type processorMock struct{ mock.Mock }

func (c *processorMock) Open(l *Loader) *configdir.ConfigDir {
	return &configdir.ConfigDir{}
}

func (c *processorMock) Retrieve(l *Loader) ([]byte, error) {
	return []byte{'f', 'o', 'o'}, nil
}

func (c *processorMock) Directory(l *Loader, f *configdir.ConfigDir) *configdir.Config {
	return &configdir.Config{}
}

type processorMockErrorRaised struct{ mock.Mock }

func (c *processorMockErrorRaised) Open(l *Loader) *configdir.ConfigDir {
	return &configdir.ConfigDir{}
}

func (c *processorMockErrorRaised) Retrieve(l *Loader) ([]byte, error) {
	return []byte{}, errors.New("error message here")
}

func (c *processorMockErrorRaised) Directory(l *Loader, f *configdir.ConfigDir) *configdir.Config {
	return nil
}

func TestProcessFile(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		want          []byte
		wantErr       bool
		wantErrMsg    string
		processorMock ConfigProcessor
	}{
		{
			name:          "No error processing file returns no error and config bytes",
			args:          []string{"some vendor", "some application", "some filename"},
			want:          []byte{'f', 'o', 'o'},
			wantErr:       false,
			processorMock: &processorMock{},
		},
		{
			name:          "When an error is raised by Retrieve(), config is empty and error is returned",
			args:          []string{"some vendor", "some application", "some filename"},
			want:          []byte{},
			wantErr:       true,
			wantErrMsg:    "error message here",
			processorMock: &processorMockErrorRaised{},
		},
	}

	for _, tt := range tests {
		processor = tt.processorMock

		got, err := ProcessFile(tt.args[0], tt.args[1], tt.args[2])

		if (err != nil) != tt.wantErr {
			t.Errorf("ProcessFile() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		if (err != nil) && err.Error() != tt.wantErrMsg {
			t.Errorf("ProcessFile() error message = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
			return
		}

		assert.Equal(t, got, tt.want, "returned file content should match")
	}
}
