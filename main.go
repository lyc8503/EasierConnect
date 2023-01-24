package main

import (
	"EasierConnect/core"
	"flag"
	"fmt"
	"log"
	"runtime"
)

func main() {
	// CLI args
	host, port, username, password, socksBind, twfId := "", 0, "", "", "", ""
	flag.StringVar(&host, "server", "", "EasyConnect server address (e.g. vpn.nju.edu.cn)")
	flag.StringVar(&username, "username", "", "Your username")
	flag.StringVar(&password, "password", "", "Your password")
	flag.StringVar(&socksBind, "socks-bind", ":1080", "The addr socks5 server listens on (e.g. 0.0.0.0:1080)")
	flag.StringVar(&twfId, "twf-id", "", "Login using twfID captured (mostly for debug usage)")
	flag.IntVar(&port, "port", 443, "EasyConnect port address (e.g. 443)")
	debugDump := false
	flag.BoolVar(&debugDump, "debug-dump", false, "Enable traffic debug dump (only for debug usage)")
	flag.Parse()

	if host == "" || ((username == "" || password == "") && twfId == "") {
		log.Fatal("Missing required cli args, refer to `EasierConnect --help`.")
	}
	server := fmt.Sprintf("%s:%d", host, port)

	client := core.NewEasyConnectClient(server)

	var ip []byte
	var err error
	if twfId != "" {
		if len(twfId) != 16 {
			panic("len(twfid) should be 16!")
		}
		ip, err = client.LoginByTwfId(twfId)
	} else {
		ip, err = client.Login(username, password)
		if err == core.ERR_NEXT_AUTH_SMS {
			fmt.Print(">>>Please enter your sms code<<<:")
			smsCode := ""
			fmt.Scan(&smsCode)

			ip, err = client.AuthSMSCode(smsCode)
		}
	}

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Login success, your IP: %d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])

	client.ServeSocks5(socksBind, debugDump)

	runtime.KeepAlive(client)
}
