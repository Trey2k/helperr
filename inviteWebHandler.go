package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/Trey2k/helperr/common"
	"github.com/gorilla/mux"
)

func (helperr *sHelperr) InviteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	invite, ok := helperr.Invites[id]
	if !ok {
		common.WarningLogger.Printf("Invalid invite id '%s'", id)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if invite.Expires.Before(time.Now()) {
		common.WarningLogger.Printf("Expired invite id '%s'", id)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	bytes, err := os.ReadFile("templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		common.ErrorLogger.Println(err)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		common.ErrorLogger.Println(err)
		return
	}
}

func (helperr *sHelperr) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	invite, ok := helperr.Invites[id]
	if !ok {
		common.WarningLogger.Printf("Invalid invite id '%s'", id)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if invite.Expires.Before(time.Now()) {
		common.WarningLogger.Printf("Expired invite id '%s'", id)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	users, err := helperr.JF.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		common.ErrorLogger.Println(err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		common.ErrorLogger.Println(err)
		return
	}

	if len(r.Form["username"]) < 1 || len(r.Form["password"]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		common.WarningLogger.Println("Invalid signup request. Missing info.")
		return
	}

	username := r.Form["username"][0]
	password := r.Form["password"][0]

	if len(username) < 3 {
		fmt.Fprint(w, "Username too short!")
		common.WarningLogger.Printf("username provided already taken '%s'", username)
		return
	}

	if len(password) < 8 || !verifyPassword(password) {
		fmt.Fprint(w, "Password too weak! Must contain 1 Upper case letter, 1 Lower case letter, 1 special charecter, 1 number and be atleast 8 charecters long.")
		common.WarningLogger.Printf("username provided already taken '%s'", username)
		return
	}

	for i := 0; i < len(users); i++ {
		if users[i].Name == username {
			fmt.Fprint(w, "Username taken!")
			common.WarningLogger.Printf("username provided already taken '%s'", username)
			return
		}
	}

	userInfo, err := helperr.JF.NewUser(username, password)
	if err != nil || userInfo == nil {
		common.ErrorLogger.Printf("failed to create user '%s'", username)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = helperr.DeleteInvite(id)
	if err != nil {
		common.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sql := `
	INSERT INTO users (discordid, jellyfinid)
	VALUES ($1, $2)
	`

	_, err = helperr.DB.Exec(sql, helperr.Invites[id].InvitedUserID, userInfo.Id)
	if err != nil {
		common.ErrorLogger.Printf("failed to create user db relation link for '%s'", username)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	common.InfoLogger.Printf("User account created for %s", username)

	helperr.Invites[id] = nil

	http.Redirect(w, r, common.Config.Jellyfin.PublicURI, http.StatusTemporaryRedirect)
}

func verifyPassword(s string) bool {
	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case c == '#' || c == '|':
			return false
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasNumber && hasUpperCase && hasLowercase && hasSpecial
}
