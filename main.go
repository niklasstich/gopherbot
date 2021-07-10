package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/gopherbot/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Connecting to discord...")
	discordSess, err := discordgo.New("Bot " + config.Conf.Discord.BotToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	//add intents - seems to be new
	discordSess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged


	err = discordSess.Open()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer discordSess.Close()


	discordSess.AddHandler(pingpong)

	log.Println("Connected and running! CTRL-C to stop.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sig
	log.Println("Got termination signal, shutting down gracefully.")
}

//returns pingpong message
func pingpong(s *discordgo.Session, m *discordgo.MessageCreate) {
	//ignore self
	if m.Author.ID==s.State.User.ID {
		return
	}

	if m.Content=="ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, config.Conf.Application.PingpongMessage)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}