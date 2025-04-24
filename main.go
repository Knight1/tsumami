package main

import (
	"flag"
)

var (
	max_tries_per_download       = 3
	max_nr_of_magazines_per_year = 13 //number of issues per year, e.g. ct=27, ix=13 (due to special editions)
	LOGIN_URL                    = "https://www.heise.de/sso/login/login"
)

// Array all magazines
func main() {

	email := flag.String("email", "", "email address")
	password := flag.String("password", "", "password")
	flag.Bool("help", false, "help")
	flag.Bool("version", false, "version")
	flag.Bool("verbose", false, "Verbose Logging")
	flag.String("magazines", "", "magazines")
	//flag.String("beginYear")
	//flag.String("endYear")
	flag.Parse()

	loginSSO(*email, *password)

}
