package main

import (
	"fmt"
	"time"

	"github.com/Trey2k/helperr/common"
	"github.com/bwmarrin/discordgo"
)

var InviteCommandInfo = &discordgo.ApplicationCommand{
	Name:        "invite",
	Description: "Invite someone to jellyfin",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "user",
			Description: "The user to invite to jellyfin",
			Required:    true,
			Type:        discordgo.ApplicationCommandOptionUser,
		},
	},
}

func (helperr *sHelperr) InviteCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	if len(options) < 1 {
		common.ErrorLogger.Panicln("No option data provided")
		return
	}

	if options[0].Name != "user" {
		common.ErrorLogger.Panicln("Invalid option data provided")
		return
	}

	id, err := helperr.newInvite(options[0].UserValue(helperr.DiscordBot.Session))
	if err != nil {
		common.ErrorLogger.Println(err)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("Invite sent! Invite code will expire on %s", helperr.Invites[id].Expires.Format(time.RFC822)),
		},
	})
	if err != nil {
		common.ErrorLogger.Println(err)
		return
	}
}
