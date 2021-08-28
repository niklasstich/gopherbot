package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/gopherbot/userdata"
)

var BalanceCommand = discordgo.ApplicationCommand{
	ID:            "",
	ApplicationID: "",
	Type:          discordgo.ChatApplicationCommand,
	Name:          "balance",
	Description:   "Gets your own point balance or balance of a specific user",
	Version:       "1",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user who's balance you are trying to get.",
			Required:    false,
		},
	},
}

func BalanceHandler(s *discordgo.Session, interact *discordgo.InteractionCreate) {
	dcUser := interact.Member.User
	if len(interact.ApplicationCommandData().Options) > 0 && interact.ApplicationCommandData().Options[0] != nil {
		dcUser = interact.ApplicationCommandData().Options[0].UserValue(s)
	}

	if dcUser.ID == s.State.User.ID {
		sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
			Content: "Yeah, I bet you'd want to know that.",
		})
		return
	}

	user := userdata.GetUser(dcUser.ID)
	if user == nil {
		var err error
		user, err = userdata.CreateUser(dcUser.ID)
		if err != nil {
			sendDbErrorResponse(s, interact)
			return
		}
	}

	sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
		Content: fmt.Sprintf(":coin: User %s has %d points!", dcUser.String(), user.Points),
	})
}
