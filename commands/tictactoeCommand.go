package commands

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var TictactoeCommand = discordgo.ApplicationCommand{
	ID:            "",
	ApplicationID: "",
	Type:          discordgo.ChatApplicationCommand,
	Name:          "tictactoe",
	Description:   "Play Tic Tac Toe with a friend!",
	Version:       "1",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The opponent who you want to offer a game to!",
			Required:    true,
		},
		{
			Type: discordgo.ApplicationCommandOptionUser,
		},
	},
}

func HellodiscordHandler(s *discordgo.Session, interact *discordgo.InteractionCreate) {
	//TODO: game logic and change original message here
	sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
		Content: "Hello Discord!",
	})
}

func TictactoeHandler(s *discordgo.Session, interact *discordgo.InteractionCreate) {
	err := s.InteractionRespond(interact.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Foo bar? Bar foo!",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Style: discordgo.PrimaryButton,
							Emoji: discordgo.ComponentEmoji{
								Name:     "⬜",
								Animated: false,
							},
							Disabled: false,
							CustomID: "hellodiscord",
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Error("failed: ", err.Error())
	}
}
