package main

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

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

			parts := strings.Split(addr, ":")
			target, err := net.ResolveIPAddr("ip", parts[0])
			if err != nil {
				return nil, errors.New("resolve ip addr failed: " + parts[0])
			}

			port, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, errors.New("invalid port: " + parts[1])
			}

			addrTarget := tcpip.FullAddress{
				NIC:  defaultNIC,
				Port: uint16(port),
				Addr: tcpip.Address(target.IP),
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
