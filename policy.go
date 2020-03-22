package redtape

import (
	"encoding/json"

	"github.com/fatih/structs"
)

type PolicyEffect string

const (
	PolicyEffectAllow PolicyEffect = "allow"
	PolicyEffectDeny  PolicyEffect = "deny"
)

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

type Scope struct {
	ID    string
	Value interface{}
}

type PolicyContext struct {
	Metadata map[string]interface{} `json:"metadata"`
	Scopes   map[string]Scope       `json:"scopes"`
}

type Policy interface {
	ID() string
	Description() string
	Roles() []*Role
	Resources() []string
	Actions() []string
	Conditions() Conditions
	Effect() PolicyEffect
	Context() PolicyContext
}

type policy struct {
	id         string
	desc       string
	roles      []*Role
	resources  []string
	actions    []string
	conditions Conditions
	effect     PolicyEffect
	ctx        PolicyContext
}

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

func MustNewPolicy(opts ...PolicyOption) Policy {
	p, err := NewPolicy(opts...)

	if err != nil {
		panic("failed to create new policy: " + err.Error())
	}

	return p
}

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

func (p *policy) ID() string {
	return p.id
}

func (p *policy) Description() string {
	return p.desc
}

func (p *policy) Roles() []*Role {
	return p.roles
}

func (p *policy) Resources() []string {
	return p.resources
}

func (p *policy) Actions() []string {
	return p.actions
}

func (p *policy) Context() PolicyContext {
	return p.ctx
}

func (p *policy) Conditions() Conditions {
	return p.conditions
}

func (p *policy) Effect() PolicyEffect {
	return p.effect
}

type PolicyOptions struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Roles       []*Role            `json:"roles"`
	Resources   []string           `json:"resources"`
	Actions     []string           `json:"actions"`
	Conditions  []ConditionOptions `json:"conditions"`
	Effect      string             `json:"effect"`
	Context     PolicyContext      `json:"context"`
}

type PolicyOption func(*PolicyOptions)

func NewPolicyOptions(opts ...PolicyOption) PolicyOptions {
	options := PolicyOptions{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func SetPolicyOptions(opts PolicyOptions) PolicyOption {
	return func(o *PolicyOptions) {
		*o = opts
	}
}

func PolicyName(n string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Name = n
	}
}

func PolicyDescription(d string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Description = d
	}
}

func PolicyDeny() PolicyOption {
	return func(o *PolicyOptions) {
		o.Effect = "deny"
	}
}

func PolicyAllow() PolicyOption {
	return func(o *PolicyOptions) {
		o.Effect = "allow"
	}
}

func SetResources(s ...string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Resources = s
	}
}

func SetActions(s ...string) PolicyOption {
	return func(o *PolicyOptions) {
		o.Actions = s
	}
}

func SetContext(ctx PolicyContext) PolicyOption {
	return func(o *PolicyOptions) {
		o.Context = ctx
	}
}

func WithCondition(co ConditionOptions) PolicyOption {
	return func(o *PolicyOptions) {
		o.Conditions = append(o.Conditions, co)
	}
}

func WithRole(r *Role) PolicyOption {
	return func(o *PolicyOptions) {
		o.Roles = append(o.Roles, r)
	}
}
