package redtape

import (
	"regexp"
	"strings"

	"github.com/blushft/redtape/strmatch"
)

type Matcher interface {
	MatchPolicy(p Policy, def []string, val string) (bool, error)
	MatchRole(r *Role, val string) (bool, error)
}

type simpleMatcher struct{}

func NewMatcher() Matcher {
	return &simpleMatcher{}
}

func (m *simpleMatcher) MatchPolicy(p Policy, def []string, val string) (bool, error) {
	for _, h := range def {
		if strmatch.MatchWildcard(h, val) {
			return true, nil
		}
	}

	return false, nil
}

func (m *simpleMatcher) MatchRole(r *Role, val string) (bool, error) {
	er, err := r.EffectiveRoles()
	if err != nil {
		return false, err
	}

	for _, rr := range er {
		if strmatch.MatchWildcard(val, rr.ID) {
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

func NewRegexMatcher() Matcher {
	return &regexMatcher{
		startDelim: "<",
		stopDelim:  ">",
		pat:        make(map[string]*regexp.Regexp),
	}
}

func (m *regexMatcher) MatchRole(r *Role, val string) (bool, error) {
	ef, err := r.EffectiveRoles()
	if err != nil {
		return false, err
	}

	def := make([]string, 0, len(ef))
	for _, rr := range ef {
		def = append(def, rr.ID)
	}

	return m.match(def, val)
}

func (m *regexMatcher) MatchPolicy(p Policy, def []string, val string) (bool, error) {
	return m.match(def, val)
}

func (m *regexMatcher) match(def []string, val string) (bool, error) {
	for _, h := range def {
		if strings.Count(h, m.startDelim) == 0 {
			if strmatch.MatchWildcard(h, val) {
				return true, nil
			}

			continue
		}

		var reg *regexp.Regexp
		var err error

		reg, ok := m.pat[h]
		if !ok {
			reg, err = strmatch.CompileDelimitedRegex(val, '<', '>')
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
