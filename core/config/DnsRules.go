package config

import "github.com/cornelk/hashmap"

var dnsRules *hashmap.Map[string, string]

func AppendSingleDnsRule(domain string, ip string) {
	if dnsRules == nil {
		dnsRules = hashmap.New[string, string]()
	}

	dnsRules.Set(domain, ip)
}

func GetSingleDnsRule(domain string) (string, bool) {
	return dnsRules.Get(domain)
}
