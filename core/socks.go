package core

import (
	"bytes"
	"context"
	"errors"
	"log"
	"math"
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

			var target *net.IPAddr

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

			target, err = net.ResolveIPAddr("ip", ip)
			if err != nil {
				return nil, errors.New("resolve ip addr failed: " + ip)
			}

			if !useL3transport && config.IsDomainRuleAvailable() {
				log.Printf("final ip: %s", target.IP.String())
				allowedPorts, useL3transport = config.GetSingleDomainRule(target.IP.String())
			}

			if !useL3transport && config.IsIpv4RuleAvailable() {
				if DebugDump {
					log.Printf("Ipv4Rule is available ")
				}
				for _, rule := range *config.GetIpv4Rules() {
					if rule.CIDR {
						_, cidr, _ := net.ParseCIDR(rule.Rule)
						if DebugDump {
							log.Printf("Cidr test: %s %s %v", target.IP, rule.Rule, cidr.Contains(target.IP))
						}

						if cidr.Contains(target.IP) {
							if DebugDump {
								log.Printf("Cidr matched: %s %s", target.IP, rule.Rule)
							}

							useL3transport = true
							allowedPorts = rule.Ports
						}
					} else {
						if DebugDump {
							log.Printf("raw match test: %s %s", target.IP, rule.Rule)
						}

						ip1 := net.ParseIP(strings.Split(rule.Rule, "~")[0])
						ip2 := net.ParseIP(strings.Split(rule.Rule, "~")[1])

						if bytes.Compare(target.IP, ip1) >= 0 && bytes.Compare(target.IP, ip2) <= 0 {
							if DebugDump {
								log.Printf("raw matched: %s %s", ip1, ip2)
							}

							useL3transport = true
							allowedPorts = rule.Ports
						}
					}
				}
			}

			if config.IsDomainRuleAvailable() {
				allowAllWebSitesPorts, allowAllWebSites := config.GetSingleDomainRule("*")

				if allowAllWebSites {
					if allowAllWebSitesPorts[0] > 0 && allowAllWebSitesPorts[1] > 0 {
						allowedPorts[0] = int(math.Min(float64(allowedPorts[0]), float64(allowAllWebSitesPorts[0])))
						allowedPorts[1] = int(math.Max(float64(allowedPorts[1]), float64(allowAllWebSitesPorts[1])))

						useL3transport = true
					}
				}
			}

			log.Printf("Addr: %s, AllowedPorts: %v, useL3transport: %v, useCustomDns: %v, ResolvedIp: %s", addr, allowedPorts, useL3transport, hasDnsRule, ip)

			if (!useL3transport && hasDnsRule) || (useL3transport && port >= allowedPorts[0] && port <= allowedPorts[1]) {
				if network != "tcp" {
					return nil, errors.New("only support tcp")
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
