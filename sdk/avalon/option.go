package avalon

import (
	"context"
	"time"
)

type Option func(config *Config)

type Downstream func(ctx context.Context, method string, req, resp interface{}, err error) error

type Config struct {
	Timeout    time.Duration
	HostPort   string
	Retry      int
	Wait       time.Duration
	Downstream Downstream

	MethodConfig map[string]*Config
}

// maybe reflect
func (config *Config) Merge(config2 *Config) {
	if config2 == nil {
		return
	}
	if config == nil {
		*config = *config2
	}
	if config2.Downstream != nil {
		config.Downstream = config2.Downstream
	}
	if config2.Timeout != 0 {
		config.Timeout = config2.Timeout
	}
	if config2.HostPort != "" {
		config.HostPort = config2.HostPort
	}
	if config2.Retry != 0 {
		config.Retry = config2.Retry
	}
}

func (config *Config) Get(method string) *Config {
	config.Merge(config.MethodConfig[method])
	return config
}

func defaultConfig(config *Config) {
	if config.Retry == 0 {
		config.Retry = 3
	}
	if config.Timeout == 0 {
		config.Timeout = 1 * time.Second
	}
	if config.Wait == 0 {
		config.Wait = 100 * time.Millisecond
	}
	if config.Downstream == nil {
		config.Downstream = func(ctx context.Context, method string, req, resp interface{}, err error) error {
			return err
		}
	}
}

var defaultOptions = []Option{defaultConfig}
