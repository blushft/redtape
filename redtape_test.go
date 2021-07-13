package redtape

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RedtapeSuite struct {
	suite.Suite
}

func TestRedtapeSuite(t *testing.T) {
	suite.Run(t, new(RedtapeSuite))
}

func (s *RedtapeSuite) TestARoles() {
	subRole := NewRole("sub_role")
	table := []struct {
		role *Role
	}{
		{
			role: NewRole("test_role"),
		},
	}

	for _, tt := range table {
		err := tt.role.AddRole(subRole)
		s.NoError(err)

		eff, err := tt.role.EffectiveRoles()
		s.NoError(err)
		s.Greater(len(eff), 1)

		err = tt.role.AddRole(tt.role)
		s.Error(err, "should not be able to add subrole that matches parent")
		err = tt.role.AddRole(subRole)
		s.Error(err, "should not be able to add duplicate subrole")

		b, err := MatchRole(tt.role, "test*")
		s.NoError(err)
		s.True(b)
	}
}

func (s *RedtapeSuite) TestBPolicies() {
	table := []struct {
		opts PolicyOptions
	}{
		{
			opts: NewPolicyOptions(
				PolicyName("test_policy"),
				PolicyDescription("just a test"),
				SetActions("create", "delete", "update", "read"),
				SetResources("database"),
				PolicyAllow(),
				WithCondition(ConditionOptions{
					Name: "test_cond",
					Type: "bool",
					Options: map[string]interface{}{
						"value": true,
					},
				}),
				WithRole(NewRole("allow_test")),
			),
		},
	}

	man := NewManager()

	for _, tt := range table {
		p, err := NewPolicy(SetPolicyOptions(tt.opts))
		s.Require().NoError(err)

		err = man.Create(p)
		s.Require().NoError(err)
	}
}

func (s *RedtapeSuite) TestCEnforce() {
	m := NewMatcher()
	pm := NewManager()

	allow := NewRole("test.A")
	deny := NewRole("test.B")

	popts := []PolicyOptions{
		{
			Name:        "test_policy_allow",
			Description: "testing",
			Roles: []*Role{
				allow,
			},
			Resources: []string{
				"test_resource",
			},
			Actions: []string{
				"test",
			},
			Effect: "allow",
			Conditions: []ConditionOptions{
				{
					Name: "match_me",
					Type: "bool",
					Options: map[string]interface{}{
						"value": true,
					},
				},
			},
		},
		{
			Name:        "test_policy",
			Description: "testing",
			Roles: []*Role{
				deny,
			},
			Resources: []string{
				"test_resource",
			},
			Actions: []string{
				"test",
			},
			Effect: "deny",
			Conditions: []ConditionOptions{
				{
					Name: "match_me",
					Type: "bool",
					Options: map[string]interface{}{
						"value": true,
					},
				},
			},
		},
	}

	for _, po := range popts {
		err := pm.Create(MustNewPolicy(SetPolicyOptions(po)))
		s.Require().NoError(err)
	}

	e, err := NewEnforcer(pm, m, nil)
	s.Require().NoError(err)

	req := &Request{
		Resource: "test_resource",
		Action:   "test",
		Role:     "test.A",
		Context: NewRequestContext(context.TODO(), map[string]interface{}{
			"match_me": true,
		}),
	}

	err = e.Enforce(req)
	s.Require().NoError(err, "should be allowed")

	req.Role = "test.B"

	err = e.Enforce(req)
	s.Require().Error(err, "should be denied")
}
