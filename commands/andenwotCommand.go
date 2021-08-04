package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/gopherbot/userdata"
)

const andenWotCost = 50

var AndenwotCommand = discordgo.ApplicationCommand{
	ID:            "",
	ApplicationID: "",
	Name:          "andenwot",
	Description:   "Post andenwot.mp4. Costs 50 points to mention a user, otherwise free.",
	Version:       "1",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "mention",
			Description: "Mentioning a user costs 50 points.",
			Required:    false,
			Options:     nil,
		},
	},
}

func AndenwotHandler(s *discordgo.Session, interact *discordgo.InteractionCreate) {
	var pingee *discordgo.User
	var user *userdata.DBUser

	msg := "https://cdn.discordapp.com/attachments/609487802165493805/872101181181276190/andenwot.mp4"

	if len(interact.ApplicationCommandData().Options) > 0 && interact.ApplicationCommandData().Options[0] != nil {
		pingee = interact.ApplicationCommandData().Options[0].UserValue(s)
		if pingee != nil {
			user = userdata.GetUser(interact.Member.User.ID)
			if user == nil {
				var err error
				user, err = userdata.CreateUser(interact.Member.User.ID)
				if err != nil {
					sendDbErrorResponse(s, interact)
					return
				}
			}
			if !user.CanAfford(andenWotCost) {
				sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
					Content: ":x: Not enough points to ping. Try again without user argument.",
				})
				return
			}
			//cannot error at this point
			_ = user.RemovePoints(50)
			err := user.WriteToDB()
			if err != nil {
				sendDbErrorResponse(s, interact)
				return
			}
			msg = pingee.Mention() + " " + msg
		}
	}

	//finally, send the message
	sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
		Content: msg,
	})
}
