package redtape

import (
	"errors"
	"fmt"
)

const (
	maxIterDepth = 10
)

type Role struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Roles       []*Role `json:"roles"`
}

func NewRole(id, name, desc string, roles ...*Role) *Role {
	return &Role{
		ID:          id,
		Name:        name,
		Description: desc,
		Roles:       roles,
	}
}

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

func getEffectiveRoles(r *Role, iter int) ([]*Role, error) {
	if iter > maxIterDepth {
		return nil, errors.New("maximum recursion reached")
	}

	var er []*Role

	er = append(er, r)
	for _, rs := range r.Roles {
		iter++
		sr, err := getEffectiveRoles(rs, iter)
		if err != nil {
			break
		}

		er = append(er, sr...)
	}

	return er, nil
}

func (r *Role) EffectiveRoles() ([]*Role, error) {
	return getEffectiveRoles(r, 0)
}

type RoleManager struct {
	Roles []Role
}
