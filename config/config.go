package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type config struct {
	Discord     discordConfig
	Application applicationConfig
}

type discordConfig struct {
	BotToken                 string
	GuildId                  string
	AppId                    string
	ClearSlashCommandsOnQuit bool
}

type applicationConfig struct {
	PingpongMessage string
	DailyCurrency   int64
}

var Conf config

func init() {
	//set default config as backup, and so we can write a default config file if needed
	viper.SetDefault("Discord", discordConfig{
		BotToken:                 "",
		GuildId:                  "",
		AppId:                    "",
		ClearSlashCommandsOnQuit: true,
	})
	viper.SetDefault("Application", applicationConfig{
		PingpongMessage: "Ping Pong!",
		DailyCurrency:   200,
	})

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$home/.config/gopherbot")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { //create default conf in default location
			wd, err := os.Getwd()
			if err != nil {
				log.Fatal("Failed to get working directory while creating default config file: ", err.Error())
			}
			confdir := wd + "/config"
			fp := confdir + "/config.yaml"
			log.Warnf("Couldn't find configuration file, creating new one at %s with default values", fp)
			err = os.MkdirAll(confdir, 0700)
			if err != nil {
				log.Fatalf("Failed to create directory %s: %v", confdir, err)
			}
			err = viper.SafeWriteConfigAs(fp)
			if err != nil {
				if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
					log.Fatalf("Failed to find config file, but config file in default location %s aready exists. Contradictory, bailing", fp)
				} else {
					log.Fatalf("Failed to create config file in default location %s: %v", fp, err.Error())
				}
			}
		} else {
			log.Fatal("Failed to read config file: ", err.Error())
		}
	}

	//unmarshal into struct
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalf("Failed to unmarshal config into Conf struct, %v", err)
	}

	//log configuration details
	if Conf.Discord.BotToken == "" {
		log.Fatal("No application token for the bot set. Please fill out the appropriate configuration file.")
	}

	if Conf.Discord.GuildId != "" {
		log.Infof("Running in guild mode for GuildID %s.", Conf.Discord.GuildId)
	}

	if Conf.Discord.AppId == "" {
		log.Warningf("Some components will not work if AppId isn't set in the configuration.")
	}
}
