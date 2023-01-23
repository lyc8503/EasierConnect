package component

import (
	"EasierConnect/core"
	"fmt"
	"log"
	"os"
)

func Process(host, port, username, password, socksBind string) {
	// CLI args
	//host, port, username, password, socksBind := "", 0, "", "", ""

	_, debugDump := os.LookupEnv("DEBUG_DUMP")

	if host == "" || username == "" || password == "" {
		log.Fatal("Missing required cli args, refer to `EasierConnect --help`.")
	}
	server := fmt.Sprintf("%s:%s", host, port)

	// Web login part (Get TWFID & ECAgent Token => Final token used in binary stream)
	twfId := core.WebLogin(server, username, password)
	agentToken := core.ECAgentToken(server, twfId)
	token := (*[48]byte)([]byte(agentToken + twfId))

	// Query IP (keep the connection used, so it's not closed too early, otherwise i/o stream will be closed)
	ip, conn := core.MustQueryIp(server, token)
	defer conn.Close()
	log.Printf("IP: %d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])

	// Link-level endpoint used in gvisor netstack
	endpoint := &core.EasyConnectEndpoint{}
	ipStack := core.SetupStack(ip, endpoint)

	// Sangfor Easyconnect protocol
	core.StartProtocol(endpoint, server, token, &[4]byte{ip[3], ip[2], ip[1], ip[0]}, debugDump)

	// Socks5 server
	core.ServeSocks5(ipStack, ip, socksBind)
}
