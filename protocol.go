package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net"
	"os"
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

		time.Sleep(time.Second)
	}
}

func send() {

}

// func main() {

// 	server := "vpn.nju.edu.cn:443"

// 	random := make([]byte, 16)
// 	rand.Read(random)

// 	token := []byte("d5090d9753544a3e541e4a02b742a27" + "\x00" + "b014487c1c9c622e")
// 	ip := []byte{242, 40, 29, 172}

// 	conn, err := tlsConn(server)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close()

// 	go recvListen(conn, (*[48]byte)(token), (*[4]byte)(ip))

// 	// tlsConn for sending data
// 	conn, err = tlsConn(server)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close()

// 	// SEND STREAM START
// 	message := []byte{0x08, 0x00, 0x00, 0x00}
// 	message = append(message, token[:]...)
// 	message = append(message, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
// 	message = append(message, ip[:]...)

// 	n, err := conn.Write(message)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Printf("send handshake: wrote %d bytes", n)
// 	dumpHex(message[:n])

// 	reply := make([]byte, 1500)
// 	n, err = conn.Read(reply)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Printf("send handshake: read %d bytes", n)
// 	dumpHex(reply[:n])

// 	if reply[0] != 0x02 {
// 		panic("unexpected send handshake reply.")
// 	}

// 	for true {

// 		message2 := "\x45\x00\x00\x3c\x01\x6f\x00\x00\x80\x01\x8e\x40\xac\x1d\x28\x66" +
// 			"\x72\xd4\x63\xba\x08\x00\x49\x2c\x00\x01\x04\x2f\x61\x62\x63\x64" +
// 			"\x65\x66\x67\x68\x69\x6a\x6b\x6c\x6d\x6e\x6f\x70\x71\x72\x73\x74" +
// 			"\x75\x76\x77\x61\x62\x63\x64\x65\x66\x67\x68\x69"
// 		n, err = io.WriteString(conn, message2)
// 		log.Printf("send: wrote %d bytes", n)
// 		dumpHex([]byte(message2[:n]))

// 		time.Sleep(time.Second)
// 	}

// 	// // HANDSHAKE?
// 	// // message = message + "\x05\x00\x00\x00" + token + "\x00\x00\x00\x00\x00\x00\x00\x00" + ip

// 	// // HEARTBEAT?
// 	// // message = message + "\x03\x00\x00\x00" + token + "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

// }
