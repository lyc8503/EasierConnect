package main

import (
	"context"
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
	_, err = rand.Read(random) // Ignore the err
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

func MustQueryIp(server string, token *[48]byte) []byte {
	conn := MustTLSConn(server)
	defer conn.Close()

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

	reply := make([]byte, 0x40)
	n, err = conn.Read(reply)
	log.Printf("query ip: read %d bytes", n)
	DumpHex(reply[:n])

	if reply[0] != 0x00 {
		panic("unexpected query ip reply.")
	}

	return reply[4:8]
}

func StartRXStream(ctx context.Context, server string, token *[48]byte, ipRev *[4]byte, inbound chan []byte) {
	conn, err := TLSConn(server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// RECV STREAM START
	message := []byte{0x06, 0x00, 0x00, 0x00}
	message = append(message, token[:]...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
	message = append(message, ip[:]...)

	n, err := conn.Write(message)

	if err != nil {
		panic(err)
	}
	log.Printf("recv handshake: wrote %d bytes", n)
	DumpHex(message[:n])

	reply := make([]byte, 1500)
	n, err = conn.Read(reply)
	log.Printf("recv handshake: read %d bytes", n)
	DumpHex(reply[:n])

	if reply[0] != 0x01 {
		return errors.New("unexpected recv handshake reply.")
	}

	for true {
		n, err = conn.Read(reply)

		if err != nil {
			panic(err)
		}

		log.Printf("recv: read %d bytes", n)
		DumpHex(reply[:n])

		n, err = targetDev.Write(reply[:n])

		if err != nil {
			panic(err)
		}

		// if strings.Contains(string(reply[:n]), "abcdefghijklmnopqrstuvwabcdefghi") {
		// 	panic(">>> PING REPLY RECEIVED   TEST PASSED <<<")
		// }

		// time.Sleep(time.Second)
	}

	return nil
}

func StartTXStream(ctx context.Context, server string, token *[48]byte, ipRev *[4]byte, outbound chan []byte) {
	conn, err := TLSConn(server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// SEND STREAM START
	message := []byte{0x05, 0x00, 0x00, 0x00}
	message = append(message, token[:]...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
	message = append(message, ip[:]...)

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
		return errors.New("unexpected send handshake reply.")
	}

	message := make([]byte, 2000)

	n, err = conn.Write(message[:n])
	if err != nil {
		panic(err)
	}

	log.Printf("send: wrote %d bytes", n)
	DumpHex([]byte(message[:n]))

	return nil
}

func StartProtocol(ctx context.Context, server string, token *[48]byte, ipRev *[4]byte) {

	go StartRXStream()
	go StartTXStream()



}
