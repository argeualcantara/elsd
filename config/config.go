package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env"
)

// Config exposes the properties that the application uses during runtime
type Config struct {
	Address string `env:"ELS_ADDRESS" envDefault:"localhost"`
	Port    int    `env:"ELS_PORT" envDefault:"7300"`
	IsDebug bool   `env:"ELS_DEBUG" envDefault:"false"`
}

var (
	configInstance *Config
)

// Load returns a Config structure populated from environment variables
func Load() *Config {
	if configInstance != nil {
		return configInstance
	}

	// Load environment variables
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	var mutex = &sync.Mutex{}
	mutex.Lock()
	configInstance = &cfg
	mutex.Unlock()

	return configInstance
}
