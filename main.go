package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/gopherbot/commands"
	"github.com/niklasstich/gopherbot/config"
	"github.com/niklasstich/gopherbot/userdata"

	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var (
	commandList = []*discordgo.ApplicationCommand{
		&commands.PingCommand,
		&commands.DailyCurrencyCommand,
		&commands.GiftCurrencyCommand,
		&commands.StatusCommand,
		&commands.AndenwotCommand,
		&commands.BalanceCommand,
	}
	handlers = map[string]func(s *discordgo.Session, interact *discordgo.InteractionCreate){
		commands.PingCommand.Name:          commands.PingHandler,
		commands.DailyCurrencyCommand.Name: commands.DailyCurrencyHandler,
		commands.GiftCurrencyCommand.Name:  commands.GiftCurrencyHandler,
		commands.StatusCommand.Name:        commands.StatusHandler,
		commands.AndenwotCommand.Name:      commands.AndenwotHandler,
		commands.BalanceCommand.Name:       commands.BalanceHandler,
	}
	commandIds = make([]string, 0)
)

var discordSess *discordgo.Session

func init() {
	verbosePtr := flag.Bool("verbose", false, "Enable verbosed log output")
	flag.Parse()

	if *verbosePtr {
		log.SetLevel(log.TraceLevel)
	}
}

//goland:noinspection GoNilness
func main() {
	log.Println("Connecting to discord...")
	var err error
	discordSess, err = discordgo.New("Bot " + config.Conf.Discord.BotToken)
	if err != nil {
		FailFast(err.Error())
	}

	discordSess.AddHandler(func(_ *discordgo.Session, _ *discordgo.Ready) {
		log.Println("Bot connected and ready.")
	})

	//add intents - seems to be new
	discordSess.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageTyping | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildIntegrations | discordgo.IntentsAllWithoutPrivileged

	err = discordSess.Open()
	if err != nil {
		FailFast(err.Error())
	}

	defer discordSess.Close()
	defer userdata.EnsureDBClosed()

	//add first level handler for slash command interactions
	discordSess.AddHandler(interactionFLH)

	log.Info("Registering slash commands...")

	//register slash commands
	for _, cmd := range commandList {
		retval, err := discordSess.ApplicationCommandCreate(discordSess.State.User.ID, config.Conf.Discord.GuildId, cmd)
		if err != nil {
			log.Error("Failed registering slash command: ", err.Error())
			continue
		}
		log.Debugf("Registered %s command", cmd.Name)
		cmd.ID = retval.ID
		commandIds = append(commandIds, retval.ID)
	}
	log.Infof("Registered %d slash commands", len(commandList))

	//defer clearing slash commandList if set, so we clear on graceful quit
	if config.Conf.Discord.ClearSlashCommandsOnQuit {
		defer ClearSlashCommands()
	}

	log.Println("CTRL-C to stop bot.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sig
	log.Println("Got termination signal, shutting down gracefully.")
}

//interactionFLH is the first level handler for all interactions. If we know a command with the given name, we call that handler.
func interactionFLH(s *discordgo.Session, interact *discordgo.InteractionCreate) {
	if handler, ok := handlers[interact.ApplicationCommandData().Name]; ok {
		handler(s, interact)
	} else {
		log.Warningf("Interaction received for command name %s, but no handler was found.", interact.ApplicationCommandData().Name)
		commands.SendGenericErrorResponse(s, interact)
	}
}

// FailFast provides a method for quitting cleanly on unrecoverable errors
func FailFast(v ...interface{}) {
	if config.Conf.Discord.ClearSlashCommandsOnQuit {
		ClearSlashCommands()
	}
	discordSess.Close()
	userdata.EnsureDBClosed()
	log.Fatal(v)
}

// ClearSlashCommands clears all slash commandList from the Discord API
func ClearSlashCommands() {
	log.Info("Cleaning up slash commands...")
	for _, cmdId := range commandIds {
		err := discordSess.ApplicationCommandDelete(discordSess.State.User.ID, config.Conf.Discord.GuildId, cmdId)
		if err != nil {
			log.Error("Failed removing slash command: ", err.Error())
		}
		log.Debugf("Cleared %s command", cmdId)
	}
	log.Info("Cleaned up slash commandList.")
}
