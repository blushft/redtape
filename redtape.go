package redtape

var (
	DefaultMatcher Matcher
)

func init() {
	DefaultMatcher = NewMatcher()
}

func MatchRole(r *Role, val string) (bool, error) {
	return DefaultMatcher.MatchRole(r, val)
}

func MatchPolicy(p Policy, def []string, val string) (bool, error) {
	return DefaultMatcher.MatchPolicy(p, def, val)
}
