package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/gopherbot/config"
	"github.com/niklasstich/gopherbot/userdata"
	log "github.com/sirupsen/logrus"
	"time"
)

var DailyCurrencyCommand = discordgo.ApplicationCommand{
	ID:            "",
	ApplicationID: "",
	Name:          "claim",
	Description:   "Claim your daily currency, once every 24 hours!",
	Version:       "1",
	Options:       nil,
}

func DailyCurrencyHandler(s *discordgo.Session, interact *discordgo.InteractionCreate) {
	user := userdata.GetUser(interact.Member.User.ID)
	var err error
	if user == nil {
		user, err = userdata.CreateUser(interact.Member.User.ID)
		if err != nil {
			log.Errorf("Could not create user in db for id %s, bailing: %v", interact.User.ID, err.Error())
			sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
				Content: "❌ Failed to claim currency: Internal db error",
			})
			return
		}
	}

	//user cannot claim yet
	if time.Since(user.LastCurrencyClaimTime) < time.Hour*24 {
		log.Tracef("User id %s tried to claim currency early", user.ID)
		timeToWait := time.Until(user.LastCurrencyClaimTime.AddDate(0, 0, 1))
		timeToWaitMessage := fmt.Sprintf("❌ You must wait %s until you can claim your currency!",
			timeToWait.Truncate(time.Second).String())
		sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
			Content: timeToWaitMessage,
		})
		return
	}

	//user can claim points
	user.Points += config.Conf.Application.DailyCurrency
	user.LastCurrencyClaimTime = time.Now()
	err = user.WriteToDB()
	if err != nil {
		log.Errorf("Couldn't save user id %s back to db", user.ID)
		sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
			Content: "❌ Failed to claim currency: Internal db error",
		})
		return
	}
	//success
	successString := fmt.Sprintf(":white_check_mark: Successfully claimed %d points, you now have %d in total!",
		config.Conf.Application.DailyCurrency, user.Points)
	sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
		Content: successString,
	})
}
