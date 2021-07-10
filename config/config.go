package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type config struct {
	Discord DiscordConfig
	Application ApplicationConfig
}

type DiscordConfig struct {
	BotToken        string
}

type ApplicationConfig struct {
	PingpongMessage string
}

var Conf config

func init() {
	//ignore metadata
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		log.Fatal("Failed to parse config file: ", err.Error())
	}
}