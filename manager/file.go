package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/blushft/redtape"
)

type FileOptions struct {
	Name string
	Path string
}

type FileOption func(*FileOptions)

func NewFileOptions(opts ...FileOption) FileOptions {
	o := FileOptions{
		Name: "redtape",
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

type File struct {
	options FileOptions
}

func NewFile(opts ...FileOption) *File {
	return &File{
		options: NewFileOptions(opts...),
	}
}

func (f *File) policyPath() string {
	fn := fmt.Sprintf("%s.policy", f.options.Name)
	return filepath.Join(f.options.Path, fn)
}

func (f *File) loadPolicies() (map[string]redtape.Policy, error) {
	b, err := os.ReadFile(f.policyPath())
	if err != nil {
		return nil, err
	}

	m := make(map[string]redtape.Policy)
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (f *File) savePolicies(m map[string]redtape.Policy) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(f.policyPath(), b, os.ModePerm)
}

func (f *File) RoleManager() (redtape.RoleManager, error) {
	if !fileExists(f.RolePath()) {
		if err := os.WriteFile(f.RolePath(), []byte("{}"), os.ModePerm); err != nil {
			return nil, err
		}
	}

	return &fileRoleMgr{f}, nil
}

func (f *File) RolePath() string {
	fn := fmt.Sprintf("%s.roles", f.options.Name)
	return filepath.Join(f.options.Path, fn)
}

func (f *File) loadRoles() (map[string]*redtape.Role, error) {
	b, err := os.ReadFile(f.RolePath())
	if err != nil {
		return nil, err
	}

	m := make(map[string]*redtape.Role)
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (f *File) saveRoles(roles map[string]*redtape.Role) error {
	b, err := json.Marshal(roles)
	if err != nil {
		return err
	}

	return os.WriteFile(f.RolePath(), b, os.ModePerm)
}

type fileRoleMgr struct {
	mgr *File
}

func (f *fileRoleMgr) Create(role *redtape.Role) error {
	return f.writeRole(role, false)
}

func (f *fileRoleMgr) Update(role *redtape.Role) error {
	return f.writeRole(role, true)
}

func (f *fileRoleMgr) writeRole(role *redtape.Role, overwrite bool) error {
	m, err := f.mgr.loadRoles()
	if err != nil {
		return err
	}

	_, ok := m[role.ID]
	if ok && !overwrite {
		return fmt.Errorf("role %s already registered", role.ID)
	}

	m[role.ID] = role

	return f.mgr.saveRoles(m)
}

func (f *fileRoleMgr) Get(id string) (*redtape.Role, error) {
	m, err := f.mgr.loadRoles()
	if err != nil {
		return nil, err
	}

	r, ok := m[id]
	if !ok {
		return nil, fmt.Errorf("role %s not found", id)
	}

	return r, nil
}

func (f *fileRoleMgr) GetByName(name string) (*redtape.Role, error) {
	m, err := f.mgr.loadRoles()
	if err != nil {
		return nil, err
	}

	for _, r := range m {
		if r.Name == name {
			return r, nil
		}
	}

	return nil, fmt.Errorf("role name %s not found", name)
}

func (f *fileRoleMgr) Delete(id string) error {
	m, err := f.mgr.loadRoles()
	if err != nil {
		return err
	}

	delete(m, id)

	return f.mgr.saveRoles(m)
}

func (f *fileRoleMgr) All(limit, offset int) ([]*redtape.Role, error) {
	m, err := f.mgr.loadRoles()
	if err != nil {
		return nil, err
	}

	rkeys := make([]string, len(m))
	i := 0
	for k := range m {
		rkeys[i] = k
		i++
	}

	start, end := limitIndices(limit, offset, len(m))
	sort.Strings(rkeys)

	roles := make([]*redtape.Role, 0, len(rkeys[start:end]))
	for _, r := range rkeys[start:end] {
		roles = append(roles, m[r])
	}

	return roles, nil
}

func (f *fileRoleMgr) GetMatching(_ string) ([]*redtape.Role, error) {
	panic("not implemented") // TODO: Implement
}

type filePolicyMgr struct {
	mgr *File
}

func (f *filePolicyMgr) Create(p redtape.Policy) error {
	panic("not implemented") // TODO: Implement
}

func (f *filePolicyMgr) Update(_ redtape.Policy) error {
	panic("not implemented") // TODO: Implement
}

func (f *filePolicyMgr) Get(_ string) (redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
}

func (f *filePolicyMgr) Delete(_ string) error {
	panic("not implemented") // TODO: Implement
}

func (f *filePolicyMgr) All(limit int, offset int) ([]redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
}

func (f *filePolicyMgr) FindByRequest(_ *redtape.Request) ([]redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
}

func (f *filePolicyMgr) FindByRole(_ string) ([]redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
}

func (f *filePolicyMgr) FindByResource(_ string) ([]redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
}

func (f *filePolicyMgr) FindByScope(_ string) ([]redtape.Policy, error) {
	panic("not implemented") // TODO: Implement
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

func fileExists(fp string) bool {
	i, err := os.Stat(fp)
	if err != nil || err == os.ErrNotExist {
		return false
	}

	if i.IsDir() {
		return false
	}

	return true
}
