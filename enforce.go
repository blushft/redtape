package redtape

import (
	"errors"
	"fmt"
)

type Enforcer interface {
	Enforce(*Request) error
}

type enforcer struct {
	manager PolicyManager
	matcher Matcher
	auditor Auditor
}

func NewEnforcer(manager PolicyManager, matcher Matcher, auditor Auditor) (Enforcer, error) {
	return &enforcer{
		manager: manager,
		matcher: matcher,
		auditor: auditor,
	}, nil
}

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

		// check all conditions
		if !e.checkConditions(p, r) {
			continue
		}

		matched = append(matched, p)

		// deny overrides all
		if p.Effect() == PolicyEffectDeny {
			return NewErrRequestDenied(fmt.Errorf("access denied by policy %s", p.ID()))
		}

		allow = true
	}

	if !allow {
		return NewErrRequestDenied(errors.New("access denied because no policy allowed access"))
	}

	return nil
}

func (e *enforcer) checkConditions(p Policy, r *Request) bool {
	for key, cond := range p.Conditions() {
		if pass := cond.Meets(r.Context.getKey(key), r); !pass {
			return false
		}
	}

	return true
}
