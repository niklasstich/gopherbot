package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type config struct {
	Discord     discordConfig
	Application applicationConfig
}

type discordConfig struct {
	BotToken                 string
	GuildId                  string
	ClearSlashCommandsOnQuit bool
}

type applicationConfig struct {
	PingpongMessage string
}

var Conf config

func init() {
	//ignore metadata
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		log.Fatal("Failed to parse config file: ", err.Error())
	}
}
