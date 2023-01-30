package core

import (
	"gvisor.dev/gvisor/pkg/bufferv2"
	"gvisor.dev/gvisor/pkg/tcpip"
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
	OnRecv     func(buf []byte)
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

func (ep *EasyConnectEndpoint) AddHeader(stack.PacketBufferPtr) {}

func (ep *EasyConnectEndpoint) WritePackets(list stack.PacketBufferList) (int, tcpip.Error) {
	for _, packetBuffer := range list.AsSlice() {
		var buf []byte
		for _, t := range packetBuffer.AsSlices() {
			buf = append(buf, t...)
		}

		if ep.OnRecv != nil {
			ep.OnRecv(buf)
		}
	}
	return list.Len(), nil
}

func (ep *EasyConnectEndpoint) WriteTo(buf []byte) {
	if ep.IsAttached() {
		packetBuffer := stack.NewPacketBuffer(stack.PacketBufferOptions{
			Payload: bufferv2.MakeWithData(buf),
		})
		ep.dispatcher.DeliverNetworkPacket(header.IPv4ProtocolNumber, packetBuffer)
		packetBuffer.DecRef()
	}
}

func SetupStack(ip []byte, endpoint *EasyConnectEndpoint) *stack.Stack {

	// init IP stack
	ipStack := stack.New(stack.Options{
		NetworkProtocols:   []stack.NetworkProtocolFactory{ipv4.NewProtocol},
		TransportProtocols: []stack.TransportProtocolFactory{tcp.NewProtocol, udp.NewProtocol},
		HandleLocal:        true,
	})

	// create NIC associated to the endpoint
	err := ipStack.CreateNIC(defaultNIC, endpoint)
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
	}

	// other settings
	sOpt := tcpip.TCPSACKEnabled(true)
	ipStack.SetTransportProtocolOption(tcp.ProtocolNumber, &sOpt)
	cOpt := tcpip.CongestionControlOption("cubic")
	ipStack.SetTransportProtocolOption(tcp.ProtocolNumber, &cOpt)
	ipStack.AddRoute(tcpip.Route{Destination: header.IPv4EmptySubnet, NIC: defaultNIC})

	return ipStack
}
