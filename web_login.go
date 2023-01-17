package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func WebLogin() string {
	server := "https://" + "vpn.nju.edu.cn:443"
	username := "211250076"
	password := "233.6666"

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

	csrfCode := string(regexp.MustCompile(`<CSRF_RAND_CODE>(.*)</CSRF_RAND_CODE>`).FindSubmatch(buf[:n])[1])
	log.Printf("CSRF Code: %s", csrfCode)

	pubKey := rsa.PublicKey{}
	pubKey.E, _ = strconv.Atoi(rsaExp)
	moduls := big.Int{}
	moduls.SetString(rsaKey, 16)
	pubKey.N = &moduls

	encryptedPassword, err := rsa.EncryptPKCS1v15(rand.Reader, &pubKey, []byte(password+"_"+csrfCode))
	if err != nil {
		panic(err)
	}
	encryptedPasswordHex := hex.EncodeToString(encryptedPassword)
	log.Printf("Encrypted Password: %s", encryptedPasswordHex)

	addr = server + "/por/login_psw.csp?anti_replay=1&encrypt=1"
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

	if !strings.Contains(string(buf[:n]), "Auth is success") {
		panic("Login FAILED: " + string(buf[:n]))
	}

	if strings.Contains(string(buf[:n]), "<NextService>auth/sms</NextService>") {
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
		panic("not implemented: sms not required")
	}

	log.Printf("Web Login process done.")

	addr = server + "/por/conf.csp"
	log.Printf("ECAgent Request: " + addr)
	req, err = http.NewRequest("GET", addr, nil)
	req.Header.Set("Cookie", "TWFID="+twfId)

	resp, err = c.Do(req)
	if err != nil {
		panic(err)
	}

	n, _ = resp.Body.Read(buf)
	defer resp.Body.Close()

	// log.Printf(string(buf[:n]))

}
