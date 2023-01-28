package core

import (
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"

	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

// local socks5 binding
var SocksBind string

type EasyConnectClient struct {
	queryConn net.Conn
	clientIp  []byte
	token     *[48]byte
	twfId     string

	endpoint *EasyConnectEndpoint
	ipStack  *stack.Stack

	server   string
	username string
	password string
}

func NewEasyConnectClient(server string) *EasyConnectClient {
	return &EasyConnectClient{
		server: server,
	}
}

func StartClient(host string, port int, username string, password string, twfId string, debugDump bool) {
	server := fmt.Sprintf("%s:%d", host, port)

	client := NewEasyConnectClient(server)

	var ip []byte
	var err error
	if twfId != "" {
		if len(twfId) != 16 {
			panic("len(twfid) should be 16!")
		}
		ip, err = client.LoginByTwfId(twfId)
	} else {
		ip, err = client.Login(username, password)
		if err == ERR_NEXT_AUTH_SMS {
			fmt.Print(">>>Please enter your sms code<<<:")
			smsCode := ""
			_, err := fmt.Scan(&smsCode)
			if err != nil {
				panic(err)
				return
			}

			ip, err = client.AuthSMSCode(smsCode)
		}
	}

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Login success, your IP: %d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])

	client.ServeSocks5(SocksBind, debugDump)

	runtime.KeepAlive(client)
}

func (client *EasyConnectClient) Login(username string, password string) ([]byte, error) {
	client.username = username
	client.password = password

	// Web login part (Get TWFID & ECAgent Token => Final token used in binary stream)
	twfId, err := WebLogin(client.server, client.username, client.password)

	// Store TWFID for AuthSMS
	client.twfId = twfId
	if err != nil {
		return nil, err
	}

	return client.LoginByTwfId(twfId)
}

func (client *EasyConnectClient) AuthSMSCode(code string) ([]byte, error) {
	if client.twfId == "" {
		return nil, errors.New("SMS Auth not required")
	}

	twfId, err := AuthSms(client.server, client.username, client.password, client.twfId, code)
	if err != nil {
		return nil, err
	}

	return client.LoginByTwfId(twfId)
}

func (client *EasyConnectClient) LoginByTwfId(twfId string) ([]byte, error) {
	agentToken, err := ECAgentToken(client.server, twfId)
	if err != nil {
		return nil, err
	}

	client.token = (*[48]byte)([]byte(agentToken + twfId))

	// Query IP (keep the connection used so it's not closed too early, otherwise i/o stream will be closed)
	client.clientIp, client.queryConn, err = QueryIp(client.server, client.token)
	if err != nil {
		return nil, err
	}

	return client.clientIp, nil
}

func (client *EasyConnectClient) ServeSocks5(socksBind string, debugDump bool) {
	// Link-level endpoint used in gvisor netstack
	client.endpoint = &EasyConnectEndpoint{}
	client.ipStack = SetupStack(client.clientIp, client.endpoint)

	// Sangfor Easyconnect protocol
	StartProtocol(client.endpoint, client.server, client.token,
		&[4]byte{client.clientIp[3], client.clientIp[2], client.clientIp[1], client.clientIp[0]}, debugDump)

	// Socks5 server
	ServeSocks5(client.ipStack, client.clientIp, socksBind)
}
