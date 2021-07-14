package manager_test

import (
	"os"
	"testing"

	"github.com/blushft/redtape"
	"github.com/blushft/redtape/manager"
	"github.com/stretchr/testify/assert"
)

func TestFileRoleManager(t *testing.T) {
	f := manager.NewFile()
	rm, err := f.RoleManager()
	if err != nil {
		t.Fatal(err)
	}

	r := redtape.NewRole("test_role")
	if err := rm.Create(r); err != nil {
		t.Fatal(err)
	}

	r.Name = "testing"
	if err := rm.Update(r); err != nil {
		t.Fatal(err)
	}

	rid, err := rm.Get(r.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, r.ID, rid.ID)

	rname, err := rm.GetByName(rid.Name)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, r.Name, rname.Name)

	all, err := rm.All(10, 0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, all, 1)

	if err := os.Remove(f.RolePath()); err != nil {
		t.Fatal(err)
	}
}
