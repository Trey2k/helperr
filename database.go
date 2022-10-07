package main

import (
	"database/sql"
	"fmt"

	"github.com/Trey2k/helperr/common"
	_ "github.com/lib/pq"
)

func (helperr *sHelperr) initDB() error {
	conInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		common.Config.Database.Host, common.Config.Database.Port, common.Config.Database.User,
		common.Config.Database.Password, common.Config.Database.Database)

	var err error
	helperr.DB, err = sql.Open("postgres", conInfo)
	if err != nil {
		return err
	}

	err = helperr.DB.Ping()
	return err
}
