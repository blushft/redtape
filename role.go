package redtape

import (
	"fmt"
)

const (
	maxIterDepth = 10
)

// Role represents a named association to a set of permissionable capability.
type Role struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Roles       []*Role `json:"roles"`
}

// NewRole returns a Role configured with the provided options.
func NewRole(id string, roles ...*Role) *Role {
	return &Role{
		ID:    id,
		Roles: roles,
	}
}

// AddRole adds a subrole.
func (r *Role) AddRole(role *Role) error {
	if r.ID == role.ID {
		return fmt.Errorf("sub role id %s cannot match parent", role.ID)
	}
	for _, rs := range r.Roles {
		if rs.ID == role.ID {
			return fmt.Errorf("%s already contains role %s", r.ID, role.ID)
		}
	}

	r.Roles = append(r.Roles, role)

	return nil
}

func getEffectiveRoles(r *Role) []*Role {
	rm := make(map[string]bool)
	rm[r.ID] = true

	var er []*Role
	er = append(er, r)

	for _, rs := range r.Roles {
		if rm[rs.ID] {
			continue
		}

		rm[rs.ID] = true

		sr := getEffectiveRoles(rs)

		er = append(er, sr...)
	}

	return er
}

// EffectiveRoles returns a flattened slice of all roles embedded in the Role.
func (r *Role) EffectiveRoles() []*Role {
	return getEffectiveRoles(r)
}
