package strmatch

import (
	"bytes"
	"fmt"
	"regexp"
)

func ExtractDelimited(s string, delimStart, delimEnd rune) ([]string, error) {
	idxs, err := delimIndices(s, delimStart, delimEnd)
	if err != nil {
		return nil, err
	}

	matches := make([]string, len(idxs)/2)
	var end int

	for i := 0; i < len(idxs); i += 2 {
		//rval := s[end:idx[i]]
		end = idxs[i+1]
		match := s[idxs[i]+1 : end-1]

		vidx := i / 2

		matches[vidx] = match
	}

	return matches, nil
}

func CompileDelimitedRegex(s string, delimStart, delimEnd rune) (*regexp.Regexp, error) {
	idxs, err := delimIndices(s, delimStart, delimEnd)
	if err != nil {
		return nil, err
	}

	pattern := bytes.NewBufferString("")
	pattern.WriteByte('^')

	var end int

	for i := 0; i < len(idxs); i += 2 {
		raw := s[end:idxs[i]]
		end = idxs[i+1]
		patt := s[idxs[i]+1 : end-1]

		_, err := fmt.Fprintf(pattern, "%s(%s)", regexp.QuoteMeta(raw), patt)
		if err != nil {
			return nil, err
		}

	}

	raw := s[end:]
	pattern.WriteString(regexp.QuoteMeta(raw))
	pattern.WriteByte('$')

	reg, err := regexp.Compile(pattern.String())
	if err != nil {
		return nil, err
	}

	return reg, nil
}

func delimIndices(s string, delimStart, delimEnd rune) ([]int, error) {
	var level, idx int
	idxs := make([]int, 0)

	rs := []rune(s)

	for i := 0; i < len(s); i++ {
		switch rs[i] {
		case delimStart:
			if level++; level == 1 {
				idx = i
			}
		case delimEnd:
			if level--; level == 0 {
				idxs = append(idxs, idx, i+1)
			} else if level < 0 {
				return nil, fmt.Errorf("unbalanced escape sequence %q", s)
			}
		}
	}

	if level != 0 {
		return nil, fmt.Errorf("unbalanced escape sequence %q", s)
	}

	return idxs, nil
}
