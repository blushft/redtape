package redtape

import (
	"context"
	"reflect"
	"testing"
)

func newConditions() Conditions {
	c, err := NewConditions(
		[]ConditionOptions{
			{
				Name: "let-me-in",
				Type: "bool",
				Options: map[string]interface{}{
					"value": true,
				},
			},
		},
		nil,
	)
	if err != nil {
		panic(err)
	}

	return c
}

func Test_policy_MarshalJSON(t *testing.T) {
	type fields struct {
		id         string
		desc       string
		roles      []*Role
		resources  []string
		actions    []string
		conditions Conditions
		effect     PolicyEffect
		ctx        context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "test_marshal",
			fields: fields{
				id:   "test_policy",
				desc: "testing policy",
				roles: []*Role{
					NewRole("test_role"),
				},
				resources: []string{
					"test_res",
				},
				actions: []string{
					"test_action",
				},
				conditions: newConditions(),
				effect:     PolicyEffectAllow,
			},
			want:    []byte(`{"name":"test_policy","description":"testing policy","roles":[{"id":"test_role","name":"","description":"","roles":null}],"resources":["test_res"],"actions":["test_action"],"scopes":null,"conditions":[{"name":"let-me-in","type":"bool","options":{"value":true}}],"effect":"allow"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &policy{
				id:         tt.fields.id,
				desc:       tt.fields.desc,
				roles:      tt.fields.roles,
				resources:  tt.fields.resources,
				actions:    tt.fields.actions,
				conditions: tt.fields.conditions,
				effect:     tt.fields.effect,
				ctx:        tt.fields.ctx,
			}
			got, err := p.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("policy.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("policy.MarshalJSON() = \n%s, want \n%s", got, tt.want)
			}
		})
	}
}
