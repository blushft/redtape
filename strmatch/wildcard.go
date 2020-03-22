package strmatch

func MatchWildcard(search, val string) bool {
	return matchWildcard(search, val, false)
}

func MatchSimpleWildcard(search, val string) bool {
	return matchWildcard(search, val, true)
}

func matchWildcard(search, val string, simple bool) bool {
	if val == search {
		return true
	}

	if search == "*" {
		return true
	}

	rval := []rune(val)
	rsearch := []rune(search)

	return runeSearch(rval, rsearch, simple)
}

func runeSearch(val, search []rune, simple bool) bool {
	for len(search) > 0 {
		switch search[0] {
		default:
			if len(val) == 0 || val[0] != search[0] {
				return false
			}
		case '?':
			if len(val) == 0 && !simple {
				return false
			}
		case '*':
			return runeSearch(val, search[1:], simple) || (len(val) > 0 && runeSearch(val[1:], search, simple))
		}

		val = val[1:]
		search = search[1:]
	}

	return len(val) == 0 && len(search) == 0
}
