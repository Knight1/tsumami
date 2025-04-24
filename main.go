package main

import (
	"flag"
	"log"
	"strconv"
)

var (
	max_tries_per_download       = 3
	max_nr_of_magazines_per_year = 13 //number of issues per year, e.g. ct=27, ix=13 (due to special editions)
	LOGIN_URL                    = "https://www.heise.de/sso/login/login"
	verbose                      bool
)

// Array all magazines
func main() {

	email := flag.String("email", "", "email address")
	password := flag.String("password", "", "password")
	flag.Bool("help", false, "help")
	flag.Bool("version", false, "version")
	verbose = *flag.Bool("verbose", false, "Verbose Logging")
	flag.String("magazines", "", "magazines")
	//flag.String("beginYear")
	//flag.String("endYear")
	flag.Parse()

	sessionCookie, err := loginSSO(*email, *password)
	if err != nil {
		log.Fatal(err)
	}
	getAccessRights(sessionCookie)

	for i := 1; i <= max_nr_of_magazines_per_year; i++ {
		downloadMagazine(sessionCookie, "ct", "2025", strconv.Itoa(i))
	}
}
