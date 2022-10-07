package main

import (
	"database/sql"

	"github.com/Trey2k/helperr/discord"
	"github.com/Trey2k/helperr/jellyfin"
)

type sHelperr struct {
	DiscordBot *discord.SBot
	JF         *jellyfin.SJellyfin
	DB         *sql.DB

	Invites map[string]*sInvite
}

func newHelperr() (*sHelperr, error) {
	helperr := &sHelperr{}

	err := helperr.initDB()
	if err != nil {
		return nil, err
	}

	helperr.DiscordBot, err = discord.NewBot()
	if err != nil {
		return nil, err
	}

	helperr.JF, err = jellyfin.NewConnection()
	if err != nil {
		return nil, err
	}

	//err = helperr.DiscordBot.RegisterCommand(SearchCMDInfo, helperr.searchHandler)

	return helperr, err
}

func (helperr *sHelperr) Destroy() {
	helperr.DiscordBot.Destroy()
	helperr.DB.Close()
}
