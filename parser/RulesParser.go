package parser

import (
	"EasierConnect/core/config"
	"log"
	"strings"
)

func ParseResourceLists(host string, twfID string, debug bool) {
	ResourceList := config.Resource{}
	parseXml(&ResourceList, host, config.PathRlist, twfID)

	//for _, ent := range ResourceList.Rcs.Rc {
	//    if debug {
	//        log.Printf("[%s] %s %s", ent.Name, ent.Host, ent.Port)
	//    }
	//}

	for _, ent := range strings.Split(ResourceList.Dns.Data, ";") {
		dnsEntry := strings.Split(ent, ":")
		RcID := dnsEntry[0]
		domain := dnsEntry[1]
		ip := dnsEntry[2]

		if debug {
			log.Printf("[%s] %s %s", RcID, domain, ip)
		}

		config.AppendSingleDnsRule(domain, ip)
	}
}

func ParseConfLists(host string, twfID string, debug bool) {
	conf := config.Conf{}
	_ = parseXml(&conf, host, config.PathConf, twfID)
}
