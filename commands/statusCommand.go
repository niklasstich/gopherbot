package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/gopherbot/userdata"
)

var StatusCommand = discordgo.ApplicationCommand{
	ID:            "",
	ApplicationID: "",
	Name:          "status",
	Description:   "Status of bot and debug information",
	Version:       "1",
	Options:       nil,
}

func StatusHandler(s *discordgo.Session, interact *discordgo.InteractionCreate) {
	//collect data
	var users int64
	guilds := len(s.State.Ready.Guilds)
	for _, guild := range s.State.Ready.Guilds {
		users += int64(guild.MemberCount)
	}

	statString := fmt.Sprintf("Currently running on %d servers, serving %d users.\nThis is shard %d; %d shard(s) total.\n"+
		"%d Users in DB.", guilds, users, s.ShardID, s.ShardCount, userdata.UserCount())
	sendInteractionResponse(s, interact, &discordgo.InteractionResponseData{
		Content: statString,
	})
}
