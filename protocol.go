package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	tls "github.com/refraction-networking/utls"
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

func recvListen(conn *tls.UConn, token *[48]byte, ip *[4]byte) {
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
		panic("unexpected recv handshake reply.")
	}

	for true {
		n, err = conn.Read(reply)
		log.Printf("recv: read %d bytes", n)
		dumpHex(reply[:n])

		if strings.Contains(string(reply[:n]), "abcdefghijklmnopqrstuvwabcdefghi") {
			panic(">>> PING REPLY RECEIVED   TEST PASSED <<<")
		}

		time.Sleep(time.Second)
	}
}

func send() {

}

func AskIp(conn *tls.UConn, token *[48]byte) []byte {
	// ASK IP PACKET
	message := []byte{0x00, 0x00, 0x00, 0x00}
	message = append(message, token[:]...)
	message = append(message, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}...)

	n, err := conn.Write(message)
	if err != nil {
		panic(err)
	}

	log.Printf("ask ip: wrote %d bytes", n)
	dumpHex(message[:n])

	reply := make([]byte, 0x40)
	n, err = conn.Read(reply)
	log.Printf("ask ip: read %d bytes", n)
	dumpHex(reply[:n])

	if reply[0] != 0x00 {
		panic("unexpected ask ip reply.")
	}

	return reply[4:8]
}

func main() {
	server := "vpn.nju.edu.cn:443"

	token := WebLogin()
	// ask IP
	conn, err := tlsConn(server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ip := AskIp(conn, (*[48]byte)(token))
	log.Printf("IP: %q", ip)
	ip[0], ip[1], ip[2], ip[3] = ip[3], ip[2], ip[1], ip[0] // reverse the ip slice for future use

	// send conn
	conn, err = tlsConn(server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go recvListen(conn, (*[48]byte)(token), (*[4]byte)(ip))

	// tlsConn for sending data
	conn, err = tlsConn(server)
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
		panic(err)
	}
	log.Printf("send handshake: wrote %d bytes", n)
	dumpHex(message[:n])

	reply := make([]byte, 1500)
	n, err = conn.Read(reply)
	if err != nil {
		panic(err)
	}
	log.Printf("send handshake: read %d bytes", n)
	dumpHex(reply[:n])

	if reply[0] != 0x02 {
		panic("unexpected send handshake reply.")
	}

	for true {
		tmp := []byte("\x45\x00\x00\x3c\x0e\x83\x00\x00\x80\x01\x00\x00" + string([]byte{ip[3], ip[2], ip[1], ip[0]}) + "\xac\x1a\x2c\x51")
		checksum := CheckSum(tmp)
		binary.BigEndian.PutUint16(tmp[10:12], checksum)

		message2 := string(tmp) +
			"\x08\x00\x49\x27\x00\x01\x04\x34\x61\x62\x63\x64" +
			"\x65\x66\x67\x68\x69\x6a\x6b\x6c\x6d\x6e\x6f\x70\x71\x72\x73\x74" +
			"\x75\x76\x77\x61\x62\x63\x64\x65\x66\x67\x68\x69"

		n, err = io.WriteString(conn, message2)
		log.Printf("send: wrote %d bytes", n)
		dumpHex([]byte(message2[:n]))

		time.Sleep(time.Second)
	}

	// // HANDSHAKE?
	// // message = message + "\x05\x00\x00\x00" + token + "\x00\x00\x00\x00\x00\x00\x00\x00" + ip

	// // HEARTBEAT?
	// // message = message + "\x03\x00\x00\x00" + token + "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)

	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}

	if length > 0 {
		sum += uint32(data[index])
	}

	sum += (sum >> 16)
	return uint16(^sum)
}
