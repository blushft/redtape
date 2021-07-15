package redtape

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// ConditionBuilder is a typed function that returns a Condition.
type ConditionBuilder func() Condition

// ConditionRegistry is a map contiaining named ConditionBuilders.
type ConditionRegistry map[string]ConditionBuilder

// NewConditionRegistry returns a ConditionRegistry containing the default Conditions and accepts an array
// of map[string]ConditionBuilder to add custom conditions to the set.
func NewConditionRegistry(conds ...map[string]ConditionBuilder) ConditionRegistry {
	reg := ConditionRegistry{
		new(BoolCondition).Name(): func() Condition {
			return new(BoolCondition)
		},
		new(RoleEqualsCondition).Name(): func() Condition {
			return new(RoleEqualsCondition)
		},
	}

	for _, ce := range conds {
		for k, c := range ce {
			reg[k] = c
		}
	}

	return reg
}

// Condition is the interface allowing different types of conditional expressions.
type Condition interface {
	Name() string
	Meets(interface{}, *Request) bool
}

// Conditions is a map of named Conditions.
type Conditions map[string]Condition

// NewConditions accepts an array of options and an optional ConditionRegistry and returns a Conditions map.
func NewConditions(opts []ConditionOptions, reg ConditionRegistry) (Conditions, error) {
	if reg == nil {
		reg = NewConditionRegistry()
	}

	cond := make(map[string]Condition)

	for _, co := range opts {
		cf, ok := reg[co.Type]
		if !ok {
			return nil, fmt.Errorf("unknown condition type %s, is it registered?", co.Type)
		}

		nc := cf()
		if len(co.Options) > 0 {
			if err := mapstructure.Decode(co.Options, &nc); err != nil {
				return nil, err
			}
		}

		cond[co.Name] = nc
	}

	return cond, nil
}

func (c Conditions) Meets(r *Request) bool {
	meta := RequestMetadataFromContext(r.Context)
	for key, cond := range c {
		if ok := cond.Meets(meta[key], r); !ok {
			return false
		}
	}

	return true
}

// ConditionOptions contains the values used to build a Condition.
type ConditionOptions struct {
	Name    string                 `json:"name"`
	Type    string                 `json:"type"`
	Options map[string]interface{} `json:"options"`
}

// BoolCondition matches a boolean value from context to the preconfigured value.
type BoolCondition struct {
	Value bool `json:"value"`
}

// Name fulfills the Name method of Condition.
func (c *BoolCondition) Name() string {
	return "bool"
}

// Meets evaluates whether parameter val matches the Condition Value.
func (c *BoolCondition) Meets(val interface{}, _ *Request) bool {
	v, ok := val.(bool)

	return ok && v == c.Value
}

// RoleEqualsCondition matches the Request role against the required role passed to the condition.
type SubjectEqualsCondition struct{}

// Name fulfills the Name method of Condition.
func (c *SubjectEqualsCondition) Name() string {
	return "role_equals"
}

// Meets evaluates true when the role val matches Request#Role.
func (c *SubjectEqualsCondition) Meets(val interface{}, r *Request) bool {
	switch v := val.(type) {
	case string:
		return v == r.Subject
	case []string:
		for _, s := range v {
			if s == r.Subject {
				return true
			}
		}
	}

	return false
}
