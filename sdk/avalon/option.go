package avalon

import (
	"time"
)

type Config struct {
	Timeout     time.Duration
	HostPort    string
	ServiceName string
}

type ClientConfig struct {
	Retry int
	Wait  time.Duration
	Config
}

type ServerConfig struct {
	Config
}

var defaultClientConfig = &ClientConfig{Config: defaultConfig}
var defaultServerConfig = &ServerConfig{Config: defaultConfig}

var defaultConfig = Config{Timeout: 1 * time.Second, HostPort: "localhost:8888"}
