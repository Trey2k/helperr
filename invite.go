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

	// Delete get saved invites
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

	// Delete invites for users who already have accounts
	sql = `
	SELECT discordid, jellyfinid FROM users
	`

	rows, err = helperr.DB.Query(sql)
	if err != nil {
		return err
	}

	for rows.Next() {
		var discordid, jellyfinid string
		err = rows.Scan(&discordid, &jellyfinid)
		if err != nil {
			return err
		}

		for k, v := range helperr.Invites {
			if v.InvitedUserID == discordid {
				err := helperr.DeleteInvite(k)
				if err != nil {
					return err
				}

				helperr.Invites[k] = nil
			}
		}
	}

	return nil
}

func (helperr *sHelperr) newInvite(user *discordgo.User) (string, error) {
	hasAccount, err := helperr.userHasAccount(user.ID)
	if err != nil {
		return "", err
	}

	if hasAccount {
		return "", fmt.Errorf("user already has account")
	}

	InviteID := uuid.NewV4()
	helperr.Invites[InviteID.String()] = &sInvite{
		InvitedUserID: user.ID,
		Expires:       time.Now().Add(time.Hour * 24),
	}

	sql := `
	INSERT INTO invites (id, discordid, expires)
	VALUES ($1, $2, $3)
	`

	_, err = helperr.DB.Exec(sql, InviteID, user.ID, helperr.Invites[InviteID.String()].Expires)
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

func (helperr *sHelperr) userHasAccount(discordid string) (bool, error) {
	sql := `
	SELECT discordid FROM users
	`

	rows, err := helperr.DB.Query(sql)
	if err != nil {
		return true, err
	}

	for rows.Next() {
		var id string
		err = rows.Scan(&discordid)
		if err != nil {
			return true, err
		}

		if discordid == id {
			return true, nil
		}
	}

	return false, nil
}
