package redtape

import (
	"fmt"
	"sort"
	"sync"
)

type PolicyManager interface {
	Create(Policy) error
	Update(Policy) error
	Get(string) (Policy, error)
	Delete(string) error
	All(limit, offset int) ([]Policy, error)

	FindByRequest(*Request) ([]Policy, error)
	FindByRole(string) ([]Policy, error)
	FindByResource(string) ([]Policy, error)
}

type defaultManager struct {
	policies map[string]Policy
	mu       sync.RWMutex
}

func NewManager() PolicyManager {
	return &defaultManager{
		policies: make(map[string]Policy),
	}
}

func (m *defaultManager) Create(p Policy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.policies[p.ID()]; exists {
		return fmt.Errorf("policy %s already registered", p.ID())
	}

	m.policies[p.ID()] = p

	return nil
}

func (m *defaultManager) Update(p Policy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.policies[p.ID()] = p

	return nil
}

func (m *defaultManager) Get(id string) (Policy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.policies[id]
	if !ok {
		return nil, fmt.Errorf("policy %s does not exist", id)
	}

	return p, nil
}

func (m *defaultManager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.policies, id)
	return nil
}

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

func (m *defaultManager) FindByRequest(r *Request) ([]Policy, error) {
	return m.findAll()
}

func (m *defaultManager) FindByRole(_ string) ([]Policy, error) {
	return m.findAll()
}

func (m *defaultManager) FindByResource(_ string) ([]Policy, error) {
	return m.findAll()
}

func limitIndices(limit, offset, len int) (int, int) {
	if offset > len {
		return len, len
	}

	if limit+offset > len {
		return offset, len
	}

	return offset, offset + limit
}
