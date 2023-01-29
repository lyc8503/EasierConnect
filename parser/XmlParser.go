package parser

import (
	"crypto/tls"
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

func parseXml(in any, host string, path string, twfid string) string {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}

	addr := "https://" + host + path
	req, err := http.NewRequest("GET", addr, nil)
	req.Header.Set("Cookie", "TWFID="+twfid)
	req.Header.Set("Content-Type", "text/html; charset=utf-8")

	resp, err := c.Do(req)
	if err != nil {
		log.Print(err)
		log.Printf("Cannot request %s \n", path)
		return ""
	}

	buf, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	//    log.Printf("%s \n", string(buf[:]))

	err = xml.Unmarshal(buf[:], &in)
	if err != nil {
		log.Print(err)
		log.Printf("Cannot parse %s \n", path)

		return ""
	} else {
		log.Printf("Parsed %s \n", path)

		return string(buf[:])
	}

}
