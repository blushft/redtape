package redtape

import (
	"regexp"
	"strings"

	"github.com/blushft/redtape/match"
)

// Matcher provides methods to facilitate matching policies to different request elements.
type Matcher interface {
	MatchPolicy(p Policy, def []string, val string) (bool, error)
	MatchRole(r *Role, val string) (bool, error)
}

type simpleMatcher struct{}

// NewMatcher returns the default Matcher implementation.
func NewMatcher() Matcher {
	return &simpleMatcher{}
}

// MatchPolicy evaluates true when the provided val wildcard matches at least one element in def.
// If def is nil, a match is assumed against any value.
func (m *simpleMatcher) MatchPolicy(p Policy, def []string, val string) (bool, error) {
	if def == nil {
		return true, nil
	}

	for _, h := range def {
		if match.Wildcard(h, val) {
			return true, nil
		}
	}

	return false, nil
}

// MatchRole evaluates true when the provided val wildcard matches at least one role in Role#EffectiveRoles.
func (m *simpleMatcher) MatchRole(r *Role, val string) (bool, error) {
	er := r.EffectiveRoles()

	for _, rr := range er {
		if match.Wildcard(val, rr.ID) {
			return true, nil
		}
	}

	return false, nil
}

type regexMatcher struct {
	startDelim string
	stopDelim  string
	pat        map[string]*regexp.Regexp
}

// NewRegexMatcher returns a Matcher using delimited regex for matching.
func NewRegexMatcher() Matcher {
	return &regexMatcher{
		startDelim: "<",
		stopDelim:  ">",
		pat:        make(map[string]*regexp.Regexp),
	}
}

// MatchRole evaluates true when the provided val regex matches at least one role in Role#EffectiveRoles.
func (m *regexMatcher) MatchRole(r *Role, val string) (bool, error) {
	ef := r.EffectiveRoles()

	def := make([]string, 0, len(ef))
	for _, rr := range ef {
		def = append(def, rr.ID)
	}

	return m.match(def, val)
}

// MatchPolicy evaluates true when the provided val regex matches at least one element in def.
func (m *regexMatcher) MatchPolicy(p Policy, def []string, val string) (bool, error) {
	return m.match(def, val)
}

func (m *regexMatcher) match(def []string, val string) (bool, error) {
	for _, h := range def {
		if strings.Count(h, m.startDelim) == 0 {
			if match.Wildcard(h, val) {
				return true, nil
			}

			continue
		}

		var reg *regexp.Regexp
		var err error

		reg, ok := m.pat[h]
		if !ok {
			reg, err = match.CompileDelimitedRegex(val, '<', '>')
			if err != nil {
				return false, err
			}

			m.pat[h] = reg
		}

		if reg.MatchString(val) {
			return true, nil
		}
	}

	return false, nil
}
