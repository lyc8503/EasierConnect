package component

import (
	"EasierConnect/core"
	"fmt"
	"log"
)

func Process(host, port, username, password, socksBind string) {
	server := fmt.Sprintf("%s:%s", host, port)
	client := core.NewEasyConnectClient(server)

	var ip []byte
	var err error
	ip, err = client.Login(username, password)
	if err == core.ERR_NEXT_AUTH_SMS {
		// TODO: input sms code via gui
		fmt.Print(">>>Please enter your sms code<<<:")
		smsCode := ""
		fmt.Scan(&smsCode)

		ip, err = client.AuthSMSCode(smsCode)
	}

	if err != nil {
		// TODO: show error in gui
		panic(err)
	}

	log.Printf("Login success, your IP: %d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])

	client.ServeSocks5(socksBind, false)
}
