package discord

import (
	"fmt"
	"runtime"

	"github.com/Trey2k/helperr/common"
	"github.com/bwmarrin/discordgo"
)

// SBot discord bot instence
type SBot struct {
	Session *discordgo.Session

	Commands map[string]*SCommand
}

// Create a new discord bot instence
// this may be pointless as config only allowes one bot at the moment
// better to be flexible though
func NewBot() (*SBot, error) {
	bot := &SBot{}

	runtime.SetFinalizer(bot, finalizeBot)

	var err error
	bot.Session, err = discordgo.New(fmt.Sprintf("Bot %s", common.Config.Discord.BotToken))
	if err != nil {
		return nil, err
	}

	bot.setupCommands()

	err = bot.Session.Open()
	if err != nil {
		return nil, err
	}

	return bot, err
}

func (bot *SBot) MessageUser(user *discordgo.User, content string) error {
	channel, err := bot.Session.UserChannelCreate(user.ID)
	if err != nil {
		return err
	}

	_, err = bot.Session.ChannelMessageSend(channel.ID, content)
	return err
}

// If program ends normally GC wont be called
func (bot *SBot) Destroy() {
	finalizeBot(bot)
}

// GC function for bot
func finalizeBot(bot *SBot) {
	for id, cmd := range bot.Commands {
		err := bot.Session.ApplicationCommandDelete(common.Config.Discord.AppID, common.Config.Discord.GuildID, id)
		if err != nil {
			common.ErrorLogger.Printf("cannot delete slash command %q: %v", cmd.AppCmd.Name, err)
		}
	}

	err := bot.Session.Close()
	common.ErrorLogger.Printf("cannot close discord session: %v", err)
}
