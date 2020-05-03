package redtape

import (
	"context"
	"encoding/json"

	"github.com/fatih/structs"
)

// PolicyEffect type is returned by Enforcer to describe the outcome of a policy evaluation
type PolicyEffect string

const (
	// PolicyEffectAllow indicates explicit permission of the request
	PolicyEffectAllow PolicyEffect = "allow"
	// PolicyEffectDeny indicates explicti denial of the request
	PolicyEffectDeny PolicyEffect = "deny"
)

// NewPolicyEffect returns a PolicyEffect for a given string
func NewPolicyEffect(s string) PolicyEffect {
	switch s {
	case "allow":
		return PolicyEffectAllow
	case "deny":
		return PolicyEffectDeny
	default:
		return PolicyEffectDeny
	}
}

// Policy provides methods to return data about a configured policy
type Policy interface {
	ID() string
	Description() string
	Roles() []*Role
	Resources() []string
	Actions() []string
	Scopes() []string
	Conditions() Conditions
	Effect() PolicyEffect
	Context() context.Context
}

type policy struct {
	id         string
	desc       string
	roles      []*Role
	resources  []string
	actions    []string
	scopes     []string
	conditions Conditions
	effect     PolicyEffect
	ctx        context.Context
}

// NewPolicy returns a default policy implementation from a set of provided options
func NewPolicy(opts ...PolicyOption) (Policy, error) {
	o := NewPolicyOptions(opts...)

	p := &policy{
		id:        o.Name,
		desc:      o.Description,
		roles:     o.Roles,
		resources: o.Resources,
		actions:   o.Actions,
		effect:    NewPolicyEffect(o.Effect),
		ctx:       o.Context,
	}

	conds, err := NewConditions(o.Conditions, nil)
	if err != nil {
		return nil, err
	}

	p.conditions = conds

	return p, nil
}

// MustNewPolicy returns a default policy implimenation or panics on error
func MustNewPolicy(opts ...PolicyOption) Policy {
	p, err := NewPolicy(opts...)

	if err != nil {
		panic("failed to create new policy: " + err.Error())
	}

	return p
}

// MarshalJSON returns a JSON byte slice representation of the default policy implimentation
func (p *policy) MarshalJSON() ([]byte, error) {
	opts := PolicyOptions{
		Name:        p.id,
		Description: p.desc,
		Roles:       p.roles,
		Resources:   p.resources,
		Actions:     p.actions,
	}

	var copts []ConditionOptions
	for k, c := range p.conditions {
		cov := structs.Map(c)
		co := ConditionOptions{
			Name:    k,
			Type:    c.Name(),
			Options: cov,
		}
		copts = append(copts, co)
	}

	opts.Conditions = copts

	return json.Marshal(opts)
}

// ID returns the policy ID
func (p *policy) ID() string {
	return p.id
}

// Description returns the policy Description
func (p *policy) Description() string {
	return p.desc
}

// Roles returns the roles the policy applies to
func (p *policy) Roles() []*Role {
	return p.roles
}

// Resources returns the resources the policy applies to
func (p *policy) Resources() []string {
	return p.resources
}

// Actions returns the actions the policy applies to
func (p *policy) Actions() []string {
	return p.actions
}

// Scopes returns the scopes the policy applies to
func (p *policy) Scopes() []string {
	return p.scopes
}

func (p *policy) Context() context.Context {
	return p.ctx
}

// Conditions returns the Conditions used to apply the policy
func (p *policy) Conditions() Conditions {
	return p.conditions
}

// Effect returns the configured PolicyEffect
func (p *policy) Effect() PolicyEffect {
	return p.effect
}

// PolicyOptions struct allows different Policy implementations to be configured with marshalable data
type PolicyOptions struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Roles       []*Role            `json:"roles"`
	Resources   []string           `json:"resources"`
	Actions     []string           `json:"actions"`
	Scopes      []string           `json:"scopes"`
	Conditions  []ConditionOptions `json:"conditions"`
	Effect      string             `json:"effect"`
	Context     context.Context    `json:"context"`
}

// PolicyOption is a typed function allowing updates to PolicyOptions through functional options
type PolicyOption func(*PolicyOptions)

// NewPolicyOptions returns PolicyOptions configured with the provided functional options
func NewPolicyOptions(opts ...PolicyOption) PolicyOptions {
	options := PolicyOptions{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// SetPolicyOptions is a PolicyOption setting all PolicyOptions to the provided values
func SetPolicyOptions(opts PolicyOptions) PolicyOption {
	return func(o *PolicyOptions) {
		*o = opts
	}
}

// PolicyName sets the policy Name Option
func PolicyName(n string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Name = n
	}
}

// PolicyDescription sets the policy description Option
func PolicyDescription(d string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Description = d
	}
}

// PolicyDeny sets the PolicyEffect to deny
func PolicyDeny() PolicyOption {
	return func(o *PolicyOptions) {
		o.Effect = "deny"
	}
}

// PolicyAllow sets the PolicyEffect to allow
func PolicyAllow() PolicyOption {
	return func(o *PolicyOptions) {
		o.Effect = "allow"
	}
}

// SetResources replaces the option Resources with the provided values
func SetResources(s ...string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Resources = s
	}
}

// SetActions replaces the option Actions with the provided values
func SetActions(s ...string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Actions = s
	}
}

// SetContext sets the Context option
func SetContext(ctx context.Context) PolicyOption {
	return func(o *PolicyOptions) {
		o.Context = ctx
	}
}

// WithCondition adds a Condition to the Conditions option
func WithCondition(co ConditionOptions) PolicyOption {
	return func(o *PolicyOptions) {
		o.Conditions = append(o.Conditions, co)
	}
}

// WithRole adds a Role to the Roles option
func WithRole(r *Role) PolicyOption {
	return func(o *PolicyOptions) {
		o.Roles = append(o.Roles, r)
	}
}
