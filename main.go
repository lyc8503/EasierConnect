package main

import (
	"flag"
	"log"
)

func main() {
	server := ""
	flag.StringVar(&server, "server", "", "EasyConnect server address (e.g. vpn.nju.edu.cn)")
	username := ""
	flag.StringVar(&username, "username", "", "Your username")
	password := ""
	flag.StringVar(&password, "password", "", "Your password")
	flag.Parse()

	if server == "" || username == "" || password == "" {
		log.Fatal("Missing required cli args, refer to `EasierConnect --help`.")
	}

	twfId := WebLogin(server, username, password)

	for {
		err := Connect(server, twfId)
		log.Printf("ERROR OCCURRED connecting to the server, retrying: %s", err)
	}
}
