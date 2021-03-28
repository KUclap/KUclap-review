package config

import (
	"log"
	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Kind			string
	Port			string
	OriginAllowed	[]string
	Server			string
	Database		string
}

type Config struct {
	Development		Configuration
	Staging			Configuration
	PreProduction	Configuration
	Production		Configuration
	
}

func (c *Config) Read() {
	if _, err := toml.DecodeFile("config/config.toml", &c); err != nil {
		log.Fatal(err)
	}
}

func (c *Config) GetConfig() Configuration {

	switch KIND := Getenv("KIND", "development"); KIND {
		case "development":
			return	c.Development
		case "staging":
			return	c.Staging
		case "preproduction":
			return	c.PreProduction
		case "production":
			return	c.Production
		default:
			panic("required environment vairble !" + KIND)
	}

}
