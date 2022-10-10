package main

import (
	"fmt"
	"time"

	"github.com/Trey2k/helperr/common"
	"github.com/bwmarrin/discordgo"
	uuid "github.com/satori/go.uuid"
)

type sInvite struct {
	InvitedUserID string
	Expires       time.Time
}

func (helperr *sHelperr) initInvites() error {
	helperr.Invites = make(map[string]*sInvite)

	sql := `
	SELECT id, discordid, expires FROM invites
	`

	rows, err := helperr.DB.Query(sql)
	if err != nil {
		return err
	}

	for rows.Next() {
		invite := &sInvite{}
		var inviteID string
		err = rows.Scan(&inviteID, &invite.InvitedUserID, &invite.Expires)
		if err != nil {
			return err
		}

		helperr.Invites[inviteID] = invite

		if invite.Expires.Before(time.Now()) {
			helperr.DeleteInvite(inviteID)
		}
	}
	return nil
}

func (helperr *sHelperr) newInvite(user *discordgo.User) (string, error) {
	InviteID := uuid.NewV4()
	helperr.Invites[InviteID.String()] = &sInvite{
		InvitedUserID: user.ID,
		Expires:       time.Now().Add(time.Hour * 24),
	}

	sql := `
	INSERT INTO invites (id, discordid, expires)
	VALUES ($1, $2, $3)
	`

	_, err := helperr.DB.Exec(sql, InviteID, user.ID, helperr.Invites[InviteID.String()].Expires)
	if err != nil {
		return "", err
	}

	err = helperr.DiscordBot.MessageUser(user, fmt.Sprintf("Congradulations! You have been invited to the **BigNass** Jellyfin server! :partying_face:\nCreate your account here: %s/invite/%s\nHurry up becuase it expires on: `%s`", common.Config.PublicURL, InviteID,
		helperr.Invites[InviteID.String()].Expires.Format(time.RFC822)))
	return InviteID.String(), err
}

func (helperr *sHelperr) DeleteInvite(inviteID string) error {
	sql := `
	DELETE FROM public.invites WHERE id=$1
	`

	_, err := helperr.DB.Exec(sql, inviteID)
	if err != nil {
		return fmt.Errorf("failed delete invite '%s'", inviteID)
	}

	helperr.Invites[inviteID] = nil
	return nil
}
