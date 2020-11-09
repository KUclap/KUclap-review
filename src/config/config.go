package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type application struct {
	Port int
	ORIGIN_ALLOWED string
}

type production struct {
	Kind string
	Server string
	Database string
}

type development struct {
	Kind string
	Server string
	Database string
}

type Config struct {
	Application	application
	Production	production
	Development	development
}

// Read and parse the configuration file
func (c *Config) Read() {
	if _, err := toml.DecodeFile("config/config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
