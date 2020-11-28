package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type preproduction struct {
	Kind string
	Port string
	OriginAllowed string
	Server string
	Database string
}

type production struct {
	Kind string
	Port string
	OriginAllowed string
	Server string
	Database string
}

type development struct {
	Kind string
	Port string
	OriginAllowed string
	Server string
	Database string
}

// Config is struct for parse data from toml
type Config struct {
	PreProduction	preproduction
	Production		production
	Development		development
}

// Read and parse the configuration file
func (c *Config) Read() {
	if _, err := toml.DecodeFile("config/config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
