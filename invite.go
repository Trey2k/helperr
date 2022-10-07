package main

import (
	"fmt"
	"time"

	"github.com/Trey2k/helperr/common"
	"github.com/bwmarrin/discordgo"
	uuid "github.com/satori/go.uuid"
)

type sInvite struct {
	InvitedUser *discordgo.User
	Expires     time.Time
}

func (helperr *sHelperr) newInvite(user *discordgo.User) (string, error) {
	InviteID := uuid.NewV4()
	helperr.Invites[InviteID.String()] = &sInvite{
		InvitedUser: user,
		Expires:     time.Now().Add(time.Hour * 24),
	}
	err := helperr.DiscordBot.MessageUser(user, fmt.Sprintf("Congradualtions! You have been invited to the BigNass Jellyfin server! :partying_face:\nCreate your account here: %s/invite/%s\nHurry up becuase it expires on: %s", common.Config.PublicURL, InviteID,
		helperr.Invites[InviteID.String()].Expires.Format(time.RFC822)))
	return InviteID.String(), err
}
