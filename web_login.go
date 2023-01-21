package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	utls "github.com/refraction-networking/utls"
)

func WebLogin(server string, username string, password string) string {
	server = "https://" + server + ":443"

	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}

	addr := server + "/por/login_auth.csp?apiversion=1"
	log.Printf("Login Request: %s", addr)

	resp, err := c.Get(addr)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	buf := make([]byte, 40960)
	n, _ := resp.Body.Read(buf)

	twfId := string(regexp.MustCompile(`<TwfID>(.*)</TwfID>`).FindSubmatch(buf[:n])[1])
	log.Printf("Twf Id: %s", twfId)

	rsaKey := string(regexp.MustCompile(`<RSA_ENCRYPT_KEY>(.*)</RSA_ENCRYPT_KEY>`).FindSubmatch(buf[:n])[1])
	log.Printf("RSA Key: %s", rsaKey)

	rsaExp := string(regexp.MustCompile(`<RSA_ENCRYPT_EXP>(.*)</RSA_ENCRYPT_EXP>`).FindSubmatch(buf[:n])[1])
	log.Printf("RSA Exp: %s", rsaExp)

	csrfMatch := regexp.MustCompile(`<CSRF_RAND_CODE>(.*)</CSRF_RAND_CODE>`).FindSubmatch(buf[:n])
 	csrfCode := "WARNING: No Match. Maybe you're connecting to an older server? Continue anyway..."
	if csrfMatch != nil {
		csrfCode = string(csrfMatch[1])
		password += "_" + csrfCode
	}
	log.Printf("CSRF Code: %s", csrfCode)
	log.Printf("Password to encrypt: %s", password)

	pubKey := rsa.PublicKey{}
	pubKey.E, _ = strconv.Atoi(rsaExp)
	moduls := big.Int{}
	moduls.SetString(rsaKey, 16)
	pubKey.N = &moduls

	encryptedPassword, err := rsa.EncryptPKCS1v15(rand.Reader, &pubKey, []byte(password))
	if err != nil {
		panic(err)
	}
	encryptedPasswordHex := hex.EncodeToString(encryptedPassword)
	log.Printf("Encrypted Password: %s", encryptedPasswordHex)

	addr = server + "/por/login_psw.csp?anti_replay=1&encrypt=1&type=cs"
	log.Printf("Login Request: %s", addr)

	form := url.Values{
		"svpn_rand_code":    {""},
		"mitm":              {""},
		"svpn_req_randcode": {csrfCode},
		"svpn_name":         {username},
		"svpn_password":     {encryptedPasswordHex},
	}

	req, err := http.NewRequest("POST", addr, strings.NewReader(form.Encode()))
	req.Header.Set("Cookie", "TWFID="+twfId)

	resp, err = c.Do(req)
	if err != nil {
		panic(err)
	}

	n, _ = resp.Body.Read(buf)
	defer resp.Body.Close()

	if strings.Contains(string(buf[:n]), "<Result>0</Result>") {
		panic("Login FAILED: " + string(buf[:n]))
	}

	if strings.Contains(string(buf[:n]), "<NextAuth>-1</NextAuth>") {
		log.Print("No NextAuth found.")
	} else if strings.Contains(string(buf[:n]), "<NextService>auth/sms</NextService>") {
		log.Print("SMS code required")

		addr = server + "/por/login_sms.csp?apiversion=1"
		log.Printf("SMS Request: " + addr)
		req, err := http.NewRequest("POST", addr, nil)
		req.Header.Set("Cookie", "TWFID="+twfId)

		resp, err = c.Do(req)
		if err != nil {
			panic(err)
		}

		n, _ = resp.Body.Read(buf)
		defer resp.Body.Close()

		if !strings.Contains(string(buf[:n]), "验证码已发送到您的手机") {
			panic("unexpected sms resp: " + string(buf[:n]))
		}

		log.Printf("SMS Code is sent or still valid.")

		fmt.Print(">>>Please enter your sms code<<<:")
		smsCode := ""
		fmt.Scan(&smsCode)

		addr = server + "/por/login_sms1.csp?apiversion=1"
		log.Printf("SMS Request: " + addr)
		form := url.Values{
			"svpn_inputsms": {smsCode},
		}

		req, err = http.NewRequest("POST", addr, strings.NewReader(form.Encode()))
		req.Header.Set("Cookie", "TWFID="+twfId)

		resp, err = c.Do(req)
		if err != nil {
			panic(err)
		}

		n, _ = resp.Body.Read(buf)
		defer resp.Body.Close()

		if !strings.Contains(string(buf[:n]), "Auth sms suc") {
			panic("SMS Code verification FAILED: " + string(buf[:n]))
		}

		twfId = string(regexp.MustCompile(`<TwfID>(.*)</TwfID>`).FindSubmatch(buf[:n])[1])
		log.Print("SMS Code verification SUCCESS")

	} else {
		panic("Not implemented auth: " + string(buf[:n]))
	}

	log.Printf("Web Login process done.")

	return twfId
}

func ECAgentToken(server string, twfId string) string {
	dialConn, err := net.Dial("tcp", server+":443")
	defer dialConn.Close()
	conn := utls.UClient(dialConn, &utls.Config{InsecureSkipVerify: true}, utls.HelloGolang)
	defer conn.Close()

	// WTF???
	// When you establish a HTTPS connection to server and send a valid request with TWFID to it
	// The **TLS ServerHello SessionId** is the first part of token
	log.Printf("ECAgent Request: /por/conf.csp & /por/rclist.csp")
	io.WriteString(conn, "GET /por/conf.csp HTTP/1.1\r\nHost: "+server+"\r\nCookie: TWFID="+twfId+"\r\n\r\nGET /por/rclist.csp HTTP/1.1\r\nHost: "+server+"\r\nCookie: TWFID="+twfId+"\r\n\r\n")

	log.Printf("Server Session ID: %q", conn.HandshakeState.ServerHello.SessionId)

	buf := make([]byte, 40960)
	n, err := conn.Read(buf)
	if n == 0 || err != nil {
		panic("ECAgent Request invalid: error " + err.Error() + "\n" + string(buf[:n]))
	}

	return hex.EncodeToString(conn.HandshakeState.ServerHello.SessionId)[:31] + "\x00"
}
