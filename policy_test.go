package redtape

import (
	"reflect"
	"testing"
)

func newConditions() Conditions {
	c, err := NewConditions(
		[]ConditionOptions{
			{
				Name: "office-ip",
				Type: "ip_whitelist",
				Options: map[string]interface{}{
					"networks": []string{
						"192.168.1.0/24",
					},
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
		ctx        PolicyContext
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
					NewRole("test_role", "Test", "Testing"),
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
			want:    nil,
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
				t.Errorf("policy.MarshalJSON() = %s, want %v", got, tt.want)
			}
		})
	}
}
