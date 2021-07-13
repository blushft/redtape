package redtape

import (
	"fmt"
	"sort"
	"sync"
)

// PolicyManager contains methods to allow query, update, and removal of policies.
type PolicyManager interface {
	Create(Policy) error
	Update(Policy) error
	Get(string) (Policy, error)
	Delete(string) error
	All(limit, offset int) ([]Policy, error)

	FindByRequest(*Request) ([]Policy, error)
	FindByRole(string) ([]Policy, error)
	FindByResource(string) ([]Policy, error)
	FindByScope(string) ([]Policy, error)
}

type defaultManager struct {
	policies map[string]Policy
	mu       sync.RWMutex
}

// NewManager returns a default memory backed policy manager.
func NewManager() PolicyManager {
	return &defaultManager{
		policies: make(map[string]Policy),
	}
}

// Create adds a policy to the manager.
func (m *defaultManager) Create(p Policy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.policies[p.ID()]; exists {
		return fmt.Errorf("policy %s already registered", p.ID())
	}

	m.policies[p.ID()] = p

	return nil
}

// Update replaces a named policy with the provided policy.
func (m *defaultManager) Update(p Policy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.policies[p.ID()] = p

	return nil
}

// Get retrieves a policy by id or error if one does not exist.
func (m *defaultManager) Get(id string) (Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.policies[id]
	if !ok {
		return nil, fmt.Errorf("policy %s does not exist", id)
	}

	return p, nil
}

// Delete removes a policy by id.
func (m *defaultManager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.policies, id)
	return nil
}

// All returns a slice containing all policies.
func (m *defaultManager) All(limit int, offset int) ([]Policy, error) {
	m.mu.RLock()

	pkeys := make([]string, len(m.policies))
	i := 0
	for k := range m.policies {
		pkeys[i] = k
		i++
	}

	start, end := limitIndices(limit, offset, len(m.policies))
	sort.Strings(pkeys)

	pols := make([]Policy, 0, len(pkeys[start:end]))
	for _, p := range pkeys[start:end] {
		pols = append(pols, m.policies[p])
	}

	m.mu.RUnlock()

	return pols, nil
}

func (m *defaultManager) findAll() ([]Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ps := make([]Policy, 0, len(m.policies))
	for _, p := range m.policies {
		ps = append(ps, p)
	}

	return ps, nil
}

// FindByRequest returns all policies matching a Request.
func (m *defaultManager) FindByRequest(r *Request) ([]Policy, error) {
	return m.findAll()
}

// FindByRole returns all policies matching a Role.
func (m *defaultManager) FindByRole(_ string) ([]Policy, error) {
	return m.findAll()
}

// FindByResource returns all policies matching a Resource.
func (m *defaultManager) FindByResource(_ string) ([]Policy, error) {
	return m.findAll()
}

// FindByResource returns all policies matching a Resource.
func (m *defaultManager) FindByScope(_ string) ([]Policy, error) {
	return m.findAll()
}

// RoleManager provides methods to store and retrieve role sets.
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

	for _, r := range m.roles {
		if name == r.Name {
			return r, nil
		}
	}

	return nil, fmt.Errorf("role %s does not exist", name)
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

func limitIndices(limit, offset, length int) (int, int) {
	if offset > length {
		return length, length
	}

	if limit+offset > length {
		return offset, length
	}

	return offset, offset + limit
}
