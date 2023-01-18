package main

import (
	"context"
	"log"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"gvisor.dev/gvisor/pkg/tcpip/transport/tcp"
	"gvisor.dev/gvisor/pkg/tcpip/transport/udp"
)

const defaultNIC tcpip.NICID = 1
const defaultMTU uint32 = 1400

// implements LinkEndpoint
type EasyConnectEndpoint struct {
	dispatcher stack.NetworkDispatcher
	inbound    chan []byte
	outbound   chan []byte
}

func (ep *EasyConnectEndpoint) MTU() uint32 {
	return defaultMTU
}

func (ep *EasyConnectEndpoint) MaxHeaderLength() uint16 {
	return 0
}

func (ep *EasyConnectEndpoint) LinkAddress() tcpip.LinkAddress {
	return ""
}

func (ep *EasyConnectEndpoint) Capabilities() stack.LinkEndpointCapabilities {
	return stack.CapabilityNone
}

func (ep *EasyConnectEndpoint) Attach(dispatcher stack.NetworkDispatcher) {
	ep.dispatcher = dispatcher
}

func (ep *EasyConnectEndpoint) IsAttached() bool {
	return ep.dispatcher != nil
}

func (ep *EasyConnectEndpoint) Wait() {}

func (ep *EasyConnectEndpoint) ARPHardwareType() header.ARPHardwareType {
	return header.ARPHardwareNone
}

func (ep *EasyConnectEndpoint) AddHeader(buffer *stack.PacketBuffer) {}

func (ep *EasyConnectEndpoint) WritePackets(list stack.PacketBufferList) (int, tcpip.Error) {
	// ep.dispatcher.DeliverNetworkPacket()
	for _, packetBuffer := range list.AsSlice() {
		packetBuffer.IncRef()
		// select {
		// 	case <-ep.done:
		// 		return 0, &tcpip.ErrClosedForSend{}
		// 	case ep.outbound <- packetBuffer:
		// }

		log.Print(packetBuffer.AsSlices())
	}
	return list.Len(), nil
}

func SetupStack(ip []byte, inbound chan []byte, outbound chan []byte) {

	// init IP stack
	ipStack := stack.New(stack.Options{
		NetworkProtocols:   []stack.NetworkProtocolFactory{ipv4.NewProtocol},
		TransportProtocols: []stack.TransportProtocolFactory{tcp.NewProtocol, udp.NewProtocol},
		HandleLocal:        true,
	})

	// custom link endpoint & nic
	endpoint := EasyConnectEndpoint{}
	err := ipStack.CreateNIC(defaultNIC, &endpoint)
	if err != nil {
		panic(err)
	}

	// assign ip
	addr := tcpip.Address(ip)
	protoAddr := tcpip.ProtocolAddress{
		AddressWithPrefix: tcpip.AddressWithPrefix{
			Address:   addr,
			PrefixLen: 32,
		},
		Protocol: ipv4.ProtocolNumber,
	}

	err = ipStack.AddProtocolAddress(defaultNIC, protoAddr, stack.AddressProperties{})
	if err != nil {
		panic(err)
		// return nil, errors.New("parse local address ", protoAddr.AddressWithPrefix, ": ", err.String())
	}

	// other settings
	sOpt := tcpip.TCPSACKEnabled(true)
	ipStack.SetTransportProtocolOption(tcp.ProtocolNumber, &sOpt)
	cOpt := tcpip.CongestionControlOption("cubic")
	ipStack.SetTransportProtocolOption(tcp.ProtocolNumber, &cOpt)
	ipStack.AddRoute(tcpip.Route{Destination: header.IPv4EmptySubnet, NIC: defaultNIC})

	// Now the stack is available
	addrTarget := tcpip.FullAddress{
		NIC:  defaultNIC,
		Port: 2333,
		Addr: tcpip.Address([]byte{1, 1, 1, 1}),
	}

	bind := tcpip.FullAddress{
		NIC:  defaultNIC,
		Addr: tcpip.Address([]byte{1, 2, 3, 4}),
	}

	tcpConn, tcpErr := gonet.DialTCPWithBind(context.Background(), ipStack, bind, addrTarget, header.IPv4ProtocolNumber)

	log.Print(tcpErr)
	tcpConn.Write([]byte("GET /"))

	buf := make([]byte, 1024)
	tcpConn.Read(buf)
}
