package main

import (
	"flag"
	"log"
)

func main() {

	// CLI args
	server, username, password := "", "", ""
	flag.StringVar(&server, "server", "", "EasyConnect server address (e.g. vpn.nju.edu.cn)")
	flag.StringVar(&username, "username", "", "Your username")
	flag.StringVar(&password, "password", "", "Your password")
	flag.Parse()

	if server == "" || username == "" || password == "" {
		log.Fatal("Missing required cli args, refer to `EasierConnect --help`.")
	}

	// Web login part (Get TWFID & ECAgent Token => Final token used in binary stream)
	twfId := WebLogin(server, username, password)
	agentToken := ECAgentToken(server, twfId)
	token := (*[48]byte)([]byte(agentToken + twfId))

	// Query IP
	ip := MustQueryIp(server + ":443", token)
	log.Printf("IP: %d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])

	// channels for outbound & inbound (relative to local machine)
	outbound, inbound := make(chan []byte, 64), make(chan []byte, 64)
	SetupStack(ip, inbound, outbound)

}
