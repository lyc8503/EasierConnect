package config

import (
	"log"
)

var Ipv4RangeRules *[]Ipv4Range

type Ipv4Range struct {
	Cidr  string
	Ports []int
}

func AppendSingleCIDRRule(cidr string, ports []int, debug bool) {
	if Ipv4RangeRules == nil {
		Ipv4RangeRules = &[]Ipv4Range{}
	}

	if debug {
		log.Printf("AppendSingleCIDRRule: %s%v", cidr, ports)
	}

	*Ipv4RangeRules = append(*Ipv4RangeRules, Ipv4Range{Cidr: cidr, Ports: ports})
}

func GetCIDRRules() *[]Ipv4Range {
	return Ipv4RangeRules
}

func IsCIDRRuleAvailable() bool {
	return Ipv4RangeRules != nil
}

func GetCIDRRuleLen() int {
	if IsDomainRuleAvailable() {
		return len(*Ipv4RangeRules)
	} else {
		return 0
	}
}
