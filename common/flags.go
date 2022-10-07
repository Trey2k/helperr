package common

import "flag"

type SFlags struct {
	DataDir *string
}

var Flags SFlags

// We dont use init() se we can insure the order of execution.
func initFlags() {
	Flags = SFlags{}
	Flags.DataDir = flag.String("data", "/var/lib/helperr", "The path to the directory the data resides in")
	flag.Parse()
}
