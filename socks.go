package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"tailscale.com/net/socks5"
)

func ServeSocks5(ipStack *stack.Stack, selfIp []byte, bindAddr string) {
	server := socks5.Server{
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {

			if network != "tcp" {
				return nil, errors.New("only support tcp")
			}

			var ip1, ip2, ip3, ip4 byte
			var port uint16
			n, err := fmt.Sscanf(addr, "%d.%d.%d.%d:%d", &ip1, &ip2, &ip3, &ip4, &port)

			if n != 5 || err != nil {
				return nil, errors.New("parse ipv4 addr failed")
			}

			addrTarget := tcpip.FullAddress{
				NIC:  defaultNIC,
				Port: port,
				Addr: tcpip.Address([]byte{ip1, ip2, ip3, ip4}),
			}

			bind := tcpip.FullAddress{
				NIC:  defaultNIC,
				Addr: tcpip.Address(selfIp),
			}

			return gonet.DialTCPWithBind(context.Background(), ipStack, bind, addrTarget, header.IPv4ProtocolNumber)
		},
	}

	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		panic("socks listen failed: " + err.Error())
	}

	log.Printf(">>>SOCKS5 SERVER listening on<<<: " + bindAddr)

	err = server.Serve(listener)
	panic(err)
}
