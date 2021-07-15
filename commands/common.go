package commands

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

//utility file that has a few functions that mutliple commands use

func sendInteractionResponse(s *discordgo.Session, interact *discordgo.InteractionCreate, rData *discordgo.InteractionResponseData) {
	err := s.InteractionRespond(interact.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: rData,
	})
	if err != nil {
		log.Error("Failed to send error response: ", err.Error())
	}
}

func sendDbErrorResponse(s *discordgo.Session, interact *discordgo.InteractionCreate) {
	sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
		Content: ":x: Internal database error",
	})
}
