package redtape

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

const (
	maxIterDepth = 10
)

// Role represents a named association to a set of permissionable capability
type Role struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Roles       []*Role `json:"roles"`
}

// NewRole returns a Role configured with the provided options
func NewRole(id string, roles ...*Role) *Role {
	return &Role{
		ID:    id,
		Roles: roles,
	}
}

// AddRole adds a subrole
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

// EffectiveRoles returns a flattened slice of all roles embedded in the Role
func (r *Role) EffectiveRoles() ([]*Role, error) {
	return getEffectiveRoles(r, 0)
}

// RoleManager provides methods to store and retrieve role sets
type RoleManager interface {
	Create(*Role) error
	Update(*Role) error
	Get(string) (*Role, error)
	GetByName(string) (*Role, error)
	Delete(string) error
	All(limit, offset int) ([]*Role, error)

	GetMatching(string) ([]*Role, error)
}

type defaultRoleManager struct {
	roles map[string]*Role
	mu    sync.RWMutex
}

func NewRoleManager() RoleManager {
	return &defaultRoleManager{
		roles: make(map[string]*Role),
	}
}

func (m *defaultRoleManager) Create(r *Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.roles[r.ID]; exists {
		return fmt.Errorf("role %s already registered", r.ID)
	}

	m.roles[r.ID] = r

	return nil
}

func (m *defaultRoleManager) Update(r *Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.roles[r.ID] = r

	return nil
}

func (m *defaultRoleManager) Get(id string) (*Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.roles[id]
	if !ok {
		return nil, fmt.Errorf("role %s does not exist", id)
	}

	return r, nil
}

func (m *defaultRoleManager) GetByName(name string) (*Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	match := new(Role)

	for _, r := range m.roles {
		if name == r.Name {
			match = r
			break
		}
	}

	if match == nil {
		return nil, fmt.Errorf("role %s does not exist", name)
	}

	return match, nil
}

func (m *defaultRoleManager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.roles, id)
	return nil
}

func (m *defaultRoleManager) All(limit, offset int) ([]*Role, error) {
	m.mu.RLock()

	rkeys := make([]string, len(m.roles))
	i := 0
	for k := range m.roles {
		rkeys[i] = k
		i++
	}

	start, end := limitIndices(limit, offset, len(m.roles))
	sort.Strings(rkeys)

	roles := make([]*Role, 0, len(rkeys[start:end]))
	for _, r := range rkeys[start:end] {
		roles = append(roles, m.roles[r])
	}

	m.mu.RUnlock()

	return roles, nil
}

func (m *defaultRoleManager) GetMatching(id string) ([]*Role, error) {
	panic("not implemented")
}
