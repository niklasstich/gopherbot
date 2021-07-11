package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/gopherbot/config"
	log "github.com/sirupsen/logrus"
)

var PingCommand = discordgo.ApplicationCommand{
	ID:            "",
	ApplicationID: "",
	Name:          "pingpong",
	Description:   "Ping? Pong!",
	Version:       "1",
	Options:       nil,
}

func PingHandler(s *discordgo.Session, interact *discordgo.InteractionCreate) {
	err := s.InteractionRespond(interact.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: config.Conf.Application.PingpongMessage,
		},
	})
	if err != nil {
		log.Error("Failed to pong the ping :( ", err.Error())
	}
}
