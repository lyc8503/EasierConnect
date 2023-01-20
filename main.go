package main

import (
	"flag"
	"log"
)

func main() {

	// CLI args
	server, username, password, socksBind := "", "", "", ""
	flag.StringVar(&server, "server", "", "EasyConnect server address (e.g. vpn.nju.edu.cn)")
	flag.StringVar(&username, "username", "", "Your username")
	flag.StringVar(&password, "password", "", "Your password")
	flag.StringVar(&socksBind, "socks-bind", ":1080", "The addr socks5 server listens on (e.g. 0.0.0.0:1080)")
	debugDump := false
	flag.BoolVar(&debugDump, "debug-dump", false, "Enable traffic debug dump (only for debug usage)")
	flag.Parse()

	if server == "" || username == "" || password == "" {
		log.Fatal("Missing required cli args, refer to `EasierConnect --help`.")
	}

	// Web login part (Get TWFID & ECAgent Token => Final token used in binary stream)
	twfId := WebLogin(server, username, password)
	agentToken := ECAgentToken(server, twfId)
	token := (*[48]byte)([]byte(agentToken + twfId))

	// Query IP (keep the connection used so it's not closed too early, otherwise i/o stream will be closed)
	ip, conn := MustQueryIp(server+":443", token)
	defer conn.Close()
	log.Printf("IP: %d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])

	// channels for outbound & inbound (relative to local machine)
	outbound, inbound := make(chan []byte, 64), make(chan []byte, 64)
	ipStack := SetupStack(ip, inbound, outbound)

	// Sangfor Easyconnect protocol
	StartProtocol(inbound, outbound, server+":443", token, &[4]byte{ip[3], ip[2], ip[1], ip[0]}, debugDump)

	// Socks5 server
	ServeSocks5(ipStack, ip, socksBind)
}
