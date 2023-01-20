package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net"
	"os"

	tls "github.com/refraction-networking/utls"
)

func DumpHex(buf []byte) {
	stdoutDumper := hex.Dumper(os.Stdout)
	defer stdoutDumper.Close()
	stdoutDumper.Write(buf)
}

func TLSConn(server string) (*tls.UConn, error) {
	// dial vpn server
	dialConn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	log.Println("socket: connected to: ", dialConn.RemoteAddr())

	// using uTLS to construct a weird TLS Client Hello (required by Sangfor)
	// The VPN and HTTP Server share port 443, Sangfor uses a special SessionId to distinguish them. (which is very stupid...)
	conn := tls.UClient(dialConn, &tls.Config{InsecureSkipVerify: true}, tls.HelloCustom)

	random := make([]byte, 32)
	rand.Read(random) // Ignore the err
	conn.SetClientRandom(random)
	conn.SetTLSVers(tls.VersionTLS11, tls.VersionTLS11, []tls.TLSExtension{})
	conn.HandshakeState.Hello.Vers = tls.VersionTLS11
	conn.HandshakeState.Hello.CipherSuites = []uint16{tls.TLS_RSA_WITH_RC4_128_SHA, tls.FAKE_TLS_EMPTY_RENEGOTIATION_INFO_SCSV}
	conn.HandshakeState.Hello.CompressionMethods = []uint8{0}
	conn.HandshakeState.Hello.SessionId = []byte{'L', '3', 'I', 'P', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	log.Println("tls: connected to: ", conn.RemoteAddr())

	return conn, nil
}

func MustTLSConn(server string) *tls.UConn {
	conn, err := TLSConn(server)
	if err != nil {
		panic(err)
	}
	return conn
}

func MustQueryIp(server string, token *[48]byte) ([]byte, *tls.UConn) {
	conn := MustTLSConn(server)
	// defer conn.Close()
	// Query IP conn CAN NOT be closed, otherwise tx/rx handshake will fail

	// QUERY IP PACKET
	message := []byte{0x00, 0x00, 0x00, 0x00}
	message = append(message, token[:]...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}...)

	n, err := conn.Write(message)
	if err != nil {
		panic(err)
	}

	log.Printf("query ip: wrote %d bytes", n)
	DumpHex(message[:n])

	reply := make([]byte, 0x80)
	n, err = conn.Read(reply)
	if err != nil {
		panic(err)
	}

	log.Printf("query ip: read %d bytes", n)
	DumpHex(reply[:n])

	if reply[0] != 0x00 {
		panic("unexpected query ip reply.")
	}

	return reply[4:8], conn
}

func BlockRXStream(server string, token *[48]byte, ipRev *[4]byte, inbound chan []byte, debug bool) error {
	conn, err := TLSConn(server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// RECV STREAM START
	message := []byte{0x06, 0x00, 0x00, 0x00}
	message = append(message, token[:]...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
	message = append(message, ipRev[:]...)

	n, err := conn.Write(message)
	if err != nil {
		return err
	}
	log.Printf("recv handshake: wrote %d bytes", n)
	DumpHex(message[:n])

	reply := make([]byte, 1500)
	n, err = conn.Read(reply)
	if err != nil {
		return err
	}
	log.Printf("recv handshake: read %d bytes", n)
	DumpHex(reply[:n])

	if reply[0] != 0x01 {
		return errors.New("unexpected recv handshake reply")
	}

	for {
		n, err = conn.Read(reply)

		if err != nil {
			return err
		}

		inbound <- reply[:n]

		if debug {
			log.Printf("recv: read %d bytes", n)
			DumpHex(reply[:n])
		}
	}
}

func BlockTXStream(server string, token *[48]byte, ipRev *[4]byte, outbound chan []byte, debug bool) error {
	conn, err := TLSConn(server)
	if err != nil {
		return err
	}
	defer conn.Close()

	// SEND STREAM START
	message := []byte{0x05, 0x00, 0x00, 0x00}
	message = append(message, token[:]...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
	message = append(message, ipRev[:]...)

	n, err := conn.Write(message)
	if err != nil {
		return err
	}
	log.Printf("send handshake: wrote %d bytes", n)
	DumpHex(message[:n])

	reply := make([]byte, 1500)
	n, err = conn.Read(reply)
	if err != nil {
		return err
	}
	log.Printf("send handshake: read %d bytes", n)
	DumpHex(reply[:n])

	if reply[0] != 0x02 {
		return errors.New("unexpected send handshake reply")
	}

	for {
		message = <-outbound

		n, err = conn.Write(message)
		if err != nil {
			return err
		}

		if debug {
			log.Printf("send: wrote %d bytes", n)
			DumpHex([]byte(message[:n]))
		}
	}
}

func StartProtocol(inbound chan []byte, outbound chan []byte, server string, token *[48]byte, ipRev *[4]byte, debug bool) {
	RX := func() {
		counter := 0
		for counter < 3 {
			err := BlockRXStream(server, token, ipRev, inbound, debug)
			if err != nil {
				log.Print("Error occurred while recv, retrying: " + err.Error())
			}
			counter += 1
		}
		panic("recv retry limit exceeded.")
	}

	go RX()

	TX := func() {
		counter := 0
		for counter < 3 {
			err := BlockTXStream(server, token, ipRev, outbound, debug)
			if err != nil {
				log.Print("Error occurred while send, retrying: " + err.Error())
			}
			counter += 1
		}
		panic("send retry limit exceeded.")
	}

	go TX()
}
