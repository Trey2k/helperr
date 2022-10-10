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
	http.HandleFunc("/", httpInterceptor(r))
	r.HandleFunc("/invite/{id}", helperr.InviteHandler).Methods("GET")
	r.HandleFunc("/invite/{id}", helperr.SignUpHandler).Methods("POST")
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", common.Config.IP, common.Config.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("www/Jellyfin-Sign-Up")))
	http.Handle("/static/", fileServer)

	go srv.ListenAndServe()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	common.InfoLogger.Println("Helperr is online!")
	<-stop
	common.InfoLogger.Println("Helperr is shutting down!")
	helperr.Destroy()

}

func httpInterceptor(router http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		router.ServeHTTP(rw, req)
	}
}
