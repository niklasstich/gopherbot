package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/gopherbot/config"
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
	s.InteractionRespond(interact.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: config.Conf.Application.PingpongMessage,
		},
	})
}
