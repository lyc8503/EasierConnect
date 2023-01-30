package config

import (
	"log"
)

var Ipv4RangeRules *[]Ipv4RangeRule

// Ipv4RangeRule Ipv4 rule with range
type Ipv4RangeRule struct {
	Rule  string
	Ports []int
	CIDR  bool
}

func AppendSingleIpv4RangeRule(rule string, ports []int, cidr bool, debug bool) {
	if Ipv4RangeRules == nil {
		Ipv4RangeRules = &[]Ipv4RangeRule{}
	}

	if debug {
		log.Printf("AppendSingleIpv4RangeRule: %s%v cidr: %v", rule, ports, cidr)
	}

	*Ipv4RangeRules = append(*Ipv4RangeRules, Ipv4RangeRule{Rule: rule, Ports: ports, CIDR: cidr})
}

func GetIpv4Rules() *[]Ipv4RangeRule {
	return Ipv4RangeRules
}

func IsIpv4RuleAvailable() bool {
	return Ipv4RangeRules != nil
}

func GetIpv4RuleLen() int {
	if IsDomainRuleAvailable() {
		return len(*Ipv4RangeRules)
	} else {
		return 0
	}
}
