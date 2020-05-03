package redtape

import (
	"errors"
	"fmt"
)

// Enforcer interface provides methods to enforce policies against a request
type Enforcer interface {
	Enforce(*Request) error
}

type enforcer struct {
	manager PolicyManager
	matcher Matcher
	auditor Auditor
}

// NewEnforcer returns a default Enforcer combining a PolicyManager, Matcher, and Auditor
func NewEnforcer(manager PolicyManager, matcher Matcher, auditor Auditor) (Enforcer, error) {
	return &enforcer{
		manager: manager,
		matcher: matcher,
		auditor: auditor,
	}, nil
}

// Enforce fulfills the Enforce method of Enforcer. The default implementation matches the Request against
// the range of stored Policies and evaluating each.
// Polices are matched first by Action, then Role, Resource, Scope and finally Condition. If a match is found, the
// configured Policy Effect is applied.
// TODO: return explicit PolicyEffect and use error to indicate processing failures
func (e *enforcer) Enforce(r *Request) error {
	allow := false
	matched := []Policy{}

	pol, err := e.manager.FindByRequest(r)
	if err != nil {
		return err
	}

	for _, p := range pol {
		// match actions
		am, err := e.matcher.MatchPolicy(p, p.Actions(), r.Action)
		if err != nil {
			return err
		}

		if !am {
			continue
		}

		rm := false
		// match roles
		for _, role := range p.Roles() {
			b, err := e.matcher.MatchRole(role, r.Role)
			if err != nil {
				return err
			}

			if b {
				rm = true
				break
			}
		}

		if !rm {
			continue
		}

		// match resources
		resm, err := e.matcher.MatchPolicy(p, p.Resources(), r.Resource)
		if err != nil {
			return err
		}
		if !resm {
			continue
		}

		// match scopes
		scm, err := e.matcher.MatchPolicy(p, p.Scopes(), r.Scope)
		if err != nil {
			return err
		}
		if !scm {
			continue
		}

		// check all conditions
		if !e.checkConditions(p, r) {
			continue
		}

		matched = append(matched, p)

		// deny overrides all
		if p.Effect() == PolicyEffectDeny {
			return NewErrRequestDeniedExplicit(fmt.Errorf("access denied by policy %s", p.ID()))
		}

		allow = true
	}

	if !allow && DefaultPolicyEffect == PolicyEffectDeny {
		return NewErrRequestDeniedImplicit(errors.New("access denied because no policy allowed access"))
	}

	return nil
}

func (e *enforcer) checkConditions(p Policy, r *Request) bool {
	for key, cond := range p.Conditions() {
		meta := RequestMetadataFromContext(r.Context)
		if pass := cond.Meets(meta[key], r); !pass {
			return false
		}
	}

	return true
}
