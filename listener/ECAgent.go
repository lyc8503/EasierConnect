package listener

import (
	"EasierConnect/core"
	"crypto/rand"
	"crypto/rsa"
	REGEXP "regexp"
	"strconv"

	"encoding/hex"

	"fmt"
	"strings"

	"log"
	"net/http"
)

type Env_ struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	regexp     *REGEXP.Regexp
}

type ECAgentResult_ struct {
	server string
	port   string
	twfid  string
}

var debugDump bool
var Env Env_
var ECAgentResult ECAgentResult_

/** This is a quick implementation of ECAgent listener protocol
 */

func HelloServer(w http.ResponseWriter, req *http.Request) {
	reqMap := make(map[string]string)
	log.Printf("Simple ECAgent #> ClientRequest: %s \n", req.RequestURI)

	constructRespon := func(operate string, result string, message string, debug string) string {
		return "(\"{\\\"type\\\":\\\"" + operate + "\\\",\\\"result\\\":\\\"" + result + "\\\",\\\"message\\\":\\\"" + message + "\\\",\\\"debug\\\":\\\"" + debug + "\\\"}\");"
	}

	form := Env.regexp.FindAllString(req.RequestURI, -1)

	for _, ent := range form {
		entry := strings.Split(ent, "=")
		reqMap[entry[0]] = entry[1]

		fmt.Printf("request > %s\n", ent)
	}

	response := strings.Builder{}

	action := reqMap["op"]
	//token := reqMap["token"]
	//Guid := reqMap["Guid"]
	callback := reqMap["callback"]

	response.WriteString(callback)

	//TODO:: Optimize & reformat
	switch {
	case action == "InitECAgent":
		response.WriteString(constructRespon(action, "1", "", "CSCM_EXIST, init ok"))
		break
	case action == "GetEncryptKey":
		//https://github.com/creationix/jsbn/blob/master/README.md

		modulus := fmt.Sprintf("%02X", Env.publicKey.N)
		response.WriteString(constructRespon(action, modulus, "", "CSCM_EXIST, init ok"))
		break
	case action == "DoConfigure":
		/*
		   op=DoConfigure
		   arg1=
		   token= hex_md5(session + '__md5_salt_for_ecagent_session__')
		   Guid=
		   callback=EA_cbxxxxx
		*/
		arg1 := reqMap["arg1"]

		Configure := strings.Split(arg1, "%20")

		switch Configure[0] {
		case "SET":
			log.Printf("DoConfigure -> Set\n")

			switch Configure[1] {
			case "SERVADDR": // nice typo
				server := Configure[2]
				port := Configure[3]

				ECAgentResult.server = server
				ECAgentResult.port = port

				log.Printf("server: %s port: %s\n", server, port)

				break
			case "TWFID":
				EncryptedTwfidHex := Configure[2]

				EncryptedTwfid, err := hex.DecodeString(EncryptedTwfidHex)
				if err != nil {
					fmt.Println(err)
				}

				DecryptedTwfid, err := rsa.DecryptPKCS1v15(rand.Reader, Env.privateKey, EncryptedTwfid)
				if err != nil {
					fmt.Println(err)
				}

				ECAgentResult.twfid = string(DecryptedTwfid[:])

				log.Printf("Encrypted twfid: %s \n", EncryptedTwfidHex)
				log.Printf("Decrypted twfid: %s \n", ECAgentResult.twfid)

				go startClient(ECAgentResult)

				break
			}
		}

		response.WriteString(constructRespon(action, "1", "", ""))
		break
	case action == "CheckProxySetting":
		response.WriteString(constructRespon(action, "-1", "", ""))
		break
	case action == "TestProxyServer":
		response.WriteString(constructRespon(action, "-1", "", ""))
		break
	case action == "GetConfig":
		//TODO:: finish Config parser
		//op=GetConfig&arg1=1&token=&Guid=&callback=EA_cbxxxxx

		//arg1 == 1 -> /por/conf.csp
		//arg1 == 2 -> rlist.csp
		response.Reset()
		break
	case action == "CheckReLogin":
		response.WriteString(constructRespon(action, "1", "", ""))
		break
	case action == "UpdateControls":
		response.WriteString(constructRespon(action, "1", "", ""))
		break
	case action == "DoQueryService":
		//TODO:: handle diff types
		response.WriteString(constructRespon(action, "26", "", ""))
		break
	case action == "StartService":
		response.WriteString(constructRespon(action, "1", "", ""))
		break
	case action == "doXmlConfigure":
		response.WriteString(constructRespon(action, "1", "", ""))
		break
	case action == "__check_alive__":
		response.Reset()
		response.WriteString("e(\"1\");")
		break
	default:
		log.Printf("Unknown action %s\n", action)
	}

	w.Header().Set("Content-Type", "text/javascript; charset=UTF-8")

	fmt.Printf("response > %s \n", response.String())
	_, err := w.Write([]byte(response.String()))
	if err != nil {
		return
	}
}

func startClient(config ECAgentResult_) {
	port, error := strconv.Atoi(config.port)
	if error != nil {
		log.Fatal("Cannot parse port!")
	}
	core.StartClient(config.server, port, "", "", config.twfid, debugDump)
}

/*
*
generate Pkcs1 keys
*/
func generateKey() {
	var err error
	// generate private key
	Env.privateKey, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal("Failed generating private key")
	}

	Env.privateKey.Precompute()
	// validate private key
	err = Env.privateKey.Validate()
	if err != nil {
		log.Fatal("Failed validating private key")
	}

	Env.publicKey = &Env.privateKey.PublicKey
}

func initRegExp() {
	var err error
	Env.regexp, err = REGEXP.Compile("[a-zA-Z0-9]*=[a-zA-Z0-9_.%]*")

	if err != nil {
		log.Fatal(err)
		return
	}
}

func StartECAgent(debugMode bool) {
	Env = Env_{}
	ECAgentResult = ECAgentResult_{}

	initRegExp()
	generateKey()

	/*
	      54530,
	   	  54541,
	   	  54552,
	   	  54563,
	   	  54574,
	   	  54585,
	   	  54596,
	   	  54607
	*/

	debugDump = debugMode

	//TODO:: Handle port in use error & https://go.dev/src/crypto/tls/generate_cert.go Auto cert generate
	http.HandleFunc("/ECAgent/", HelloServer)
	err := http.ListenAndServeTLS(":54530", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
