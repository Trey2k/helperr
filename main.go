package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Trey2k/helperr/common"
	"github.com/gorilla/mux"
)

func main() {

	helperr, err := newHelperr()
	if err != nil {
		common.ErrorLogger.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/invite/{id}", InviteHandler)
	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", common.Config.IP, common.Config.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go common.ErrorLogger.Fatal(srv.ListenAndServe())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	common.InfoLogger.Println("Helperr is online!")
	<-stop
	common.InfoLogger.Println("Helperr is shutting down!")
	helperr.Destroy()

}
