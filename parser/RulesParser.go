package parser

import (
	"EasierConnect/core/config"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var domainRegExp *regexp.Regexp

func StringArrToIntArr(strArr []string) [4]int {
	var intArr [4]int

	for index, str := range strArr {
		intArr[index], _ = strconv.Atoi(str)
	}

	return intArr
}

// from https://github.com/dromara/hutool/blob/fc091b01a23271e02f3174c45942105048155c90/hutool-core/src/main/java/cn/hutool/core/net/Ipv4Util.java#L114
func getIPsInRange(from string, to string) *[]string {
	ipf := StringArrToIntArr(strings.Split(from, "."))
	ipt := StringArrToIntArr(strings.Split(to, "."))

	var ips []string

	equ := func(condition bool, yes int, no int) int {
		if condition {
			return yes
		} else {
			return no
		}
	}

	var a, b, c int

	//TODO::handle with a better method
	for a = ipf[0]; a <= ipt[0]; a++ {
		for b = equ(a == ipf[0], ipf[1], 0); b <= equ(a == ipt[0], ipt[1], 255); b++ {
			for c = equ(b == ipf[1], ipf[2], 0); c <= equ(b == ipt[1], ipt[2], 255); c++ {
				for d := equ(c == ipf[2], ipf[3], 0); d <= equ(c == ipt[2], ipt[3], 255); d++ {
					ips = append(ips, strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(c)+"."+strconv.Itoa(d))
				}
			}
		}
	}

	return &ips
}

func processSingleIpRule(rule string, port string, debug bool, waitChan chan int) {
	appendRule := func(domain *string) {
		minValue := port
		maxValue := port

		if strings.Contains(port, "~") {
			minValue = strings.Split(port, "~")[0]
			maxValue = strings.Split(port, "~")[1]
		}

		minValueInt, err := strconv.Atoi(minValue)
		if err != nil {
			log.Printf("Cannot parse port value from string")
			return
		}

		maxValueInt, err := strconv.Atoi(maxValue)
		if err != nil {
			log.Printf("Cannot parse port value from string")
			return
		}

		//	if debug {
		log.Printf("Appending Domain rule for: %s%v", *domain, []int{minValueInt, maxValueInt})
		//	}
		config.AppendSingleDomainRule(*domain, []int{minValueInt, maxValueInt}, debug)
	}

	if strings.Contains(rule, "~") { // ip range 1.1.1.7~1.1.7.9
		from := strings.Split(rule, "~")[0]
		to := strings.Split(rule, "~")[1]

		if debug {
			log.Printf("Handling rule for: %s-%s", from, to)
		}
		for _, domain := range *getIPsInRange(from, to) {
			appendRule(&domain)
		}
	} else { // http://domain.example.com/path/to&something=good#extra
		appendRule(&rule)

		if domainRegExp == nil {
			domainRegExp, _ = regexp.Compile("(?:\\w+\\.)+\\w+")
		}

		pureDomain := domainRegExp.FindString(rule)

		appendRule(&pureDomain) //TODO::FIXME:: remove this when using Http(s) proxy (i think it works on socks5)
	}

	waitChan <- 1
}

func ParseResourceLists(host string, twfID string, debug bool) {
	ResourceList := config.Resource{}
	parseXml(&ResourceList, host, config.PathRlist, twfID)

	RcsLen := len(ResourceList.Rcs.Rc)

	cpuNumber := runtime.NumCPU()
	waitChan := make(chan int, cpuNumber)

	for RcsIndex, ent := range ResourceList.Rcs.Rc {
		//	if debug {
		log.Printf("[%s] %s %s", ent.Name, ent.Host, ent.Port)
		//	}

		if ent.Host == "" || ent.Port == "" {
			break
		}

		domains := strings.Split(ent.Host, ";")
		ports := strings.Split(ent.Port, ";")

		if len(domains) >= 1 && len(ports) >= 1 {
			for index, domain := range domains {
				portRange := ports[index]

				if cpuNumber > 0 {
					cpuNumber--
				} else {
					<-waitChan
				}
				processSingleIpRule(domain, portRange, debug, waitChan)
			}
		}

		log.Printf("Progress: %v/100.00 (ResourceList.Rcs)", (float32(RcsIndex)/float32(RcsLen))*100)
	}

	log.Printf("Loaded ResourceList.Rcs")

	for _, ent := range strings.Split(ResourceList.Dns.Data, ";") {
		dnsEntry := strings.Split(ent, ":")

		if len(dnsEntry) >= 3 {
			RcID := dnsEntry[0]
			domain := dnsEntry[1]
			ip := dnsEntry[2]

			//	if debug {
			log.Printf("[%s] %s %s", RcID, domain, ip)
			//	}

			if domain != "" && ip != "" {
				config.AppendSingleDnsRule(domain, ip, debug)
			}
		}
	}

	log.Printf("Parsed %v Domain rules", config.GetDomainRuleLen())
	log.Printf("Parsed %v Dns rules", config.GetDnsRuleLen())
}

func ParseConfLists(host string, twfID string, debug bool) {
	conf := config.Conf{}
	_ = parseXml(&conf, host, config.PathConf, twfID)
}
