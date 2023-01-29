package core

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"EasierConnect/core/config"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"tailscale.com/net/socks5"
)

func ServeSocks5(ipStack *stack.Stack, selfIp []byte, bindAddr string) {
	server := socks5.Server{
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {

			log.Printf("socks dial: %s", addr)

			parts := strings.Split(addr, ":")

			ip := parts[0]
			port, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, errors.New("invalid port: " + parts[1])
			}

			var allowedPorts = []int{1, 65535} // [0] -> Min, [1] -> Max
			var useL3transport = true
			var hasDnsRule = false

			if config.IsDomainRuleAvailable() {
				allowedPorts, useL3transport = config.GetSingleDomainRule(ip)
			}

			if config.IsDnsRuleAvailable() {
				var dnsRules string
				dnsRules, hasDnsRule = config.GetSingleDnsRule(ip)

				if hasDnsRule {
					ip = dnsRules
				}
			}

			if !useL3transport && config.IsDomainRuleAvailable() {
				allowedPorts, useL3transport = config.GetSingleDomainRule(ip)
			}

			log.Printf("Addr: %s, AllowedPorts: %v, useL3transport: %v, useCustomDns: %v, ResolvedIp: %s", addr, allowedPorts, useL3transport, hasDnsRule, ip)

			if (!useL3transport && hasDnsRule) || (useL3transport && port >= allowedPorts[0] && port <= allowedPorts[1]) {
				if network != "tcp" {
					return nil, errors.New("only support tcp")
				}

				target, err := net.ResolveIPAddr("ip", ip)
				if err != nil {
					return nil, errors.New("resolve ip addr failed: " + ip)
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
			}
			goDialer := &net.Dialer{}
			goDial := goDialer.DialContext

			log.Printf("skip: %s", addr)

			return goDial(ctx, network, addr)
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
