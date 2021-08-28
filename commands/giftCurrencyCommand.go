package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/gopherbot/userdata"
	log "github.com/sirupsen/logrus"
)

var GiftCurrencyCommand = discordgo.ApplicationCommand{
	ID:            "",
	ApplicationID: "",
	Type:          discordgo.ChatApplicationCommand,
	Name:          "gift",
	Description:   "Gift your friends some of your points!",
	Version:       "1",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user you wanna give your points to.",
			Required:    true,
		},

		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "amount",
			Description: "How many points you wanna give to the user.",
			Required:    true,
		},
	},
}

type giftCurrencyArgs struct {
	user   *discordgo.User
	amount int64
}

func GiftCurrencyHandler(s *discordgo.Session, interact *discordgo.InteractionCreate) {
	//cast arguments
	args := giftCurrencyArgs{
		user:   interact.ApplicationCommandData().Options[0].UserValue(s),
		amount: interact.ApplicationCommandData().Options[1].IntValue(),
	}

	//ignore requests to gift to us
	if args.user.ID == s.State.User.ID {
		sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
			Content: ":x: I am flattered, but you cannot send me any points!",
		})
		return
	}

	//ignore requests to gift yourself
	if args.user.ID == interact.User.ID {
		sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
			Content: ":x: You cannot gift yourself points.",
		})
		return
	}
	var err error
	//check if both users exist in db, if not create them
	//gifter is the one who gifts points, giftee is the receiver
	gifterUser := userdata.GetUser(interact.Member.User.ID)
	if gifterUser == nil {
		gifterUser, err = userdata.CreateUser(interact.Member.User.ID)
		if err != nil {
			sendDbErrorResponse(s, interact)
			return
		}
	}

	gifteeUser := userdata.GetUser(args.user.ID)
	if gifteeUser == nil {
		gifteeUser, err = userdata.CreateUser(args.user.ID)
		if err != nil {
			sendDbErrorResponse(s, interact)
			return
		}
	}

	//check if gifter has enough balance
	if !gifterUser.CanAfford(args.amount) {
		sErr := fmt.Sprintf(":x: You cannot afford that! You only have %d points.", gifterUser.Points)
		sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
			Content: sErr,
		})
		return
	}

	err = gifterUser.RemovePoints(args.amount)
	if err != nil {
		sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
			Content: ":x: internal error",
		})
		log.Error("Failed to remove points after making sure user could afford it, sync issue?")
		return
	}
	gifteeUser.AddPoints(args.amount)
	err = gifterUser.WriteToDB()
	if err != nil {
		sendDbErrorResponse(s, interact)
		return
	}
	err = gifteeUser.WriteToDB()
	if err != nil {
		sendDbErrorResponse(s, interact)
		gifterUser.AddPoints(args.amount)
		if err := gifterUser.WriteToDB(); err != nil {
			log.Error("Failed to restore gifter balance after failed write to giftee!")
		}
		return
	}

	sStr := fmt.Sprintf(":white_check_mark: %s gave %s %d points!", interact.Member.User.Mention(),
		args.user.Mention(), args.amount)
	sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
		Content: sStr,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Users: []string{
				interact.Member.User.ID,
				args.user.ID,
			},
		},
	})

}
