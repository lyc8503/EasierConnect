package main

import (
	"EasierConnect/core"
	"EasierConnect/listener"
	"flag"
	"log"
)

func main() {
	// CLI args
	host, port, username, password, twfId := "", 0, "", "", ""
	flag.StringVar(&host, "server", "", "EasyConnect server address (e.g. vpn.nju.edu.cn)")
	flag.StringVar(&username, "username", "", "Your username")
	flag.StringVar(&password, "password", "", "Your password")
	flag.StringVar(&core.SocksBind, "socks-bind", ":1080", "The addr socks5 server listens on (e.g. 0.0.0.0:1080)")
	flag.StringVar(&twfId, "twf-id", "", "Login using twfID captured (mostly for debug usage)")
	flag.IntVar(&port, "port", 443, "EasyConnect port address (e.g. 443)")
	debugDump := false
	flag.BoolVar(&debugDump, "debug-dump", false, "Enable traffic debug dump (only for debug usage)")
	flag.Parse()

	if host == "" || ((username == "" || password == "") && twfId == "") {
		log.Printf("For more infomations: `EasierConnect --help`.\n")
		log.Printf("ECAgent is Listening on 54530. (EXPERIMENT)\n")
		listener.StartECAgent(debugDump)
	} else {
		core.StartClient(host, port, username, password, twfId, debugDump)
	}
}
