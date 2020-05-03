package redtape

var (
	// DefaultMatcher is a simple matcher
	DefaultMatcher = NewMatcher()
	// DefaultPolicyEffect is the policy effect to apply when no other matches can be found
	DefaultPolicyEffect = PolicyEffectDeny
)

// MatchRole is a utility function that uses the DefaultMatcher to evaluate whether role val matches the
// effective roles of r
func MatchRole(r *Role, val string) (bool, error) {
	return DefaultMatcher.MatchRole(r, val)
}

// MatchPolicy is a utility function that uses DefaultMatcher to evaluate whether p can be matched by val
func MatchPolicy(p Policy, def []string, val string) (bool, error) {
	return DefaultMatcher.MatchPolicy(p, def, val)
}
