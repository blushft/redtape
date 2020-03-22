package redtape

import (
	"net"

	"github.com/mitchellh/mapstructure"
)

type ConditionBuilder func() Condition
type ConditionRegistry map[string]ConditionBuilder

func NewConditionRegistry(conds ...map[string]ConditionBuilder) ConditionRegistry {
	reg := ConditionRegistry{
		new(BoolCondition).Name(): func() Condition {
			return new(BoolCondition)
		},
		new(RoleEqualsCondition).Name(): func() Condition {
			return new(RoleEqualsCondition)
		},
		new(IPWhitelistCondition).Name(): func() Condition {
			return new(IPWhitelistCondition)
		},
	}

	for _, ce := range conds {
		for k, c := range ce {
			reg[k] = c
		}
	}

	return reg
}

type Condition interface {
	Name() string
	Meets(interface{}, *Request) bool
}

type Conditions map[string]Condition

func NewConditions(opts []ConditionOptions, reg ConditionRegistry) (Conditions, error) {
	if reg == nil {
		reg = NewConditionRegistry()
	}

	cond := make(map[string]Condition)

	for _, co := range opts {
		if cf, ok := reg[co.Type]; ok {
			nc := cf()
			if len(co.Options) > 0 {
				if err := mapstructure.Decode(co.Options, &nc); err != nil {
					return nil, err
				}
			}

			cond[co.Name] = nc
		}
	}

	return cond, nil
}

type ConditionOptions struct {
	Name    string                 `json:"name"`
	Type    string                 `json:"type"`
	Options map[string]interface{} `json:"options"`
}

// BoolCondition matches a boolean value from context to the preconfigured value
type BoolCondition struct {
	Value bool `json:"value"`
}

func (c *BoolCondition) Name() string {
	return "bool"
}

func (c *BoolCondition) Meets(val interface{}, _ *Request) bool {
	v, ok := val.(bool)

	return ok && v == c.Value
}

// SubjectEqualsCondition matches when the named value is equal to the request subject
type RoleEqualsCondition struct{}

func (c *RoleEqualsCondition) Name() string {
	return "subject_equals"
}

func (c *RoleEqualsCondition) Meets(val interface{}, r *Request) bool {
	s, ok := val.(string)

	return ok && s == r.Role
}

type IPWhitelistCondition struct {
	Networks []string `json:"networks" structs:"networks"`
}

func (c *IPWhitelistCondition) Name() string {
	return "ip_whitelist"
}

func (c *IPWhitelistCondition) Meets(val interface{}, _ *Request) bool {
	ip, ok := val.(string)
	if !ok {
		return false
	}

	for _, ns := range c.Networks {
		_, cidr, err := net.ParseCIDR(ns)
		if err != nil {
			return false
		}

		tip := net.ParseIP(ip)
		if tip == nil {
			return false
		}

		if cidr.Contains(tip) {
			return true
		}
	}

	return false
}
