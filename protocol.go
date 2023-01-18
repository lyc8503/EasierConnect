package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	tls "github.com/refraction-networking/utls"
	"github.com/songgao/water"
)

func dumpHex(buf []byte) {
	stdoutDumper := hex.Dumper(os.Stdout)
	defer stdoutDumper.Close()
	stdoutDumper.Write(buf)
}

func tlsConn(server string) (*tls.UConn, error) {
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

func recvListen(conn *tls.UConn, token *[48]byte, ip *[4]byte, targetDev *water.Interface) error {
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
	dumpHex(message[:n])

	reply := make([]byte, 1500)
	n, err = conn.Read(reply)
	log.Printf("recv handshake: read %d bytes", n)
	dumpHex(reply[:n])

	if reply[0] != 0x01 {
		return errors.New("unexpected recv handshake reply.")
	}

	for true {
		n, err = conn.Read(reply)

		if err != nil {
			panic(err)
		}

		log.Printf("recv: read %d bytes", n)
		dumpHex(reply[:n])

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

func QueryIp(conn *tls.UConn, token *[48]byte) ([]byte, error) {
	// QUERY IP PACKET
	message := []byte{0x00, 0x00, 0x00, 0x00}
	message = append(message, token[:]...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}...)

	n, err := conn.Write(message)
	if err != nil {
		return nil, err
	}

	log.Printf("query ip: wrote %d bytes", n)
	dumpHex(message[:n])

	reply := make([]byte, 0x40)
	n, err = conn.Read(reply)
	log.Printf("query ip: read %d bytes", n)
	dumpHex(reply[:n])

	if reply[0] != 0x00 {
		panic("unexpected query ip reply.")
	}

	return reply[4:8], nil
}

func SendConnHandshake(conn *tls.UConn, token *[48]byte, ip *[4]byte) error {
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
	dumpHex(message[:n])

	reply := make([]byte, 1500)
	n, err = conn.Read(reply)
	if err != nil {
		return err
	}
	log.Printf("send handshake: read %d bytes", n)
	dumpHex(reply[:n])

	if reply[0] != 0x02 {
		return errors.New("unexpected send handshake reply.")
	}

	return nil
}

func Connect(server string, twfId string) error {
	agentToken := ECAgentToken(server, twfId)
	token := []byte(agentToken + twfId)
	server = server + ":443"

	// query IP
	conn, err := tlsConn(server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ip, err := QueryIp(conn, (*[48]byte)(token))
	if err != nil {
		panic(err)
	}

	log.Printf("IP: %d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
	ip[0], ip[1], ip[2], ip[3] = ip[3], ip[2], ip[1], ip[0] // reverse the ip slice for future use

	// TUN dev
	log.Printf("Initializing TUN device...")

	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})

	if err != nil {
		panic(err)
	}

	defer ifce.Close()

	// try to set tun0 up (on linux)
	ipStr := fmt.Sprintf("%d.%d.%d.%d", ip[3], ip[2], ip[1], ip[0])
	err = exec.Command("ifconfig", ifce.Name(), ipStr, "up").Run()
	if err != nil {
		log.Print("Failed to set TUN dev up, Try bring it up manually.")
	}

	// recv conn
	conn, err = tlsConn(server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go recvListen(conn, (*[48]byte)(token), (*[4]byte)(ip), ifce)

	// tlsConn for sending data
	conn, err = tlsConn(server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	SendConnHandshake(conn, (*[48]byte)(token), (*[4]byte)(ip))

	message := make([]byte, 2000)
	for {
		n, err := ifce.Read(message)
		if err != nil {
			panic(err)
		}

		if message[0] != 0x45 {
			log.Print("send: dropping non-acceptable (not ipv4 or ihl != 5) packet.")
			continue
		}

		n, err = conn.Write(message[:n])
		if err != nil {
			panic(err)
		}

		log.Printf("send: wrote %d bytes", n)
		dumpHex([]byte(message[:n]))
	}

	// // HANDSHAKE?
	// // message = message + "\x05\x00\x00\x00" + token + "\x00\x00\x00\x00\x00\x00\x00\x00" + ip

	// // HEARTBEAT?
	// // message = message + "\x03\x00\x00\x00" + token + "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

}
