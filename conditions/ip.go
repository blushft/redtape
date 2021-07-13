package conditions

import (
	"net"

	"github.com/blushft/redtape"
)

// IPAllowCondition performs CIDR matching for a range of Networks against a provided value and
// allows access on match.
type IPAllowCondition struct {
	Networks []string `json:"networks"`
}

// Name fulfills the Name method of Condition.
func (c *IPAllowCondition) Name() string {
	return "ip_allow"
}

// Meets evaluates true when the network address in val is contained within
// one of the CIDR ranges of IPAllowCondition#Networks.
func (c *IPAllowCondition) Meets(val interface{}, _ *redtape.Request) bool {
	ip, ok := val.(string)
	if !ok {
		return false
	}

	return matchIP(ip, c.Networks)
}

// IPAllowCondition performs CIDR matching for a range of Networks against a provided value and denys.
// access on match.
type IPDenyCondition struct {
	Networks []string `json:"networks"`
}

// Name fulfills the Name method of Condition.
func (c *IPDenyCondition) Name() string {
	return "ip_deny"
}

// Meets evaluates true when the network address in val is contained within
// one of the CIDR ranges of IPDenyCondition#Networks.
func (c *IPDenyCondition) Meets(val interface{}, _ *redtape.Request) bool {
	ip, ok := val.(string)
	if !ok {
		return false
	}

	return !matchIP(ip, c.Networks)
}

func matchIP(val string, nets []string) bool {
	for _, ns := range nets {
		_, cidr, err := net.ParseCIDR(ns)
		if err != nil {
			return false
		}

		tip := net.ParseIP(val)
		if tip == nil {
			return false
		}

		if cidr.Contains(tip) {
			return true
		}
	}

	return false
}
