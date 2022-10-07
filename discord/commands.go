package discord

import (
	"fmt"

	"github.com/Trey2k/helperr/common"
	"github.com/bwmarrin/discordgo"
)

type SCommand struct {
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
	AppCmd  *discordgo.ApplicationCommand
}

// Register a new discord command
func (bot *SBot) RegisterCommand(appCmd *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) error {

	cmd := &SCommand{
		AppCmd:  appCmd,
		Handler: handler,
	}

	rcmd, err := bot.Session.ApplicationCommandCreate(common.Config.Discord.AppID, common.Config.Discord.GuildID, cmd.AppCmd)
	if err != nil {
		return fmt.Errorf("cannot create slash command %q: %v", cmd.AppCmd.Name, err)
	}

	bot.Commands[rcmd.ID] = cmd

	return nil

}

// Initalize the slash commands
func (bot *SBot) setupCommands() {
	bot.Session.AddHandler(bot.handleCommand)

	bot.Commands = make(map[string]*SCommand)
}

// Handle the slash commands
func (bot *SBot) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	bot.Commands[i.ApplicationCommandData().ID].Handler(s, i)
}
