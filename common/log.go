package common

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func initLogs() {
	dataDir := *Flags.DataDir

	// Make log dir if it does not exists
	if _, err := os.Stat(fmt.Sprintf("%s/logs", dataDir)); err != nil {
		if err := os.MkdirAll(fmt.Sprintf("%s/logs", dataDir), 0755); err != nil {
			panic(err)
		}
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/logs/log.txt", dataDir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(mw, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
