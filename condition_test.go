package redtape

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var jsonCond = []byte(`
[
	{
		"name": "some_condition",
		"type": "bool",
		"options": {
			"value": true
		}
	}
]
`)

func TestNewConditions(t *testing.T) {
	var unmCond []ConditionOptions
	if err := json.Unmarshal(jsonCond, &unmCond); err != nil {
		t.Errorf("failed to unmarshal options: %v", err)
	}

	type args struct {
		opts []ConditionOptions
	}
	tests := []struct {
		name    string
		args    args
		test    string
		val     interface{}
		want    bool
		wantErr bool
	}{
		{
			name: "unmarshal_conditions",
			args: args{
				opts: unmCond,
			},
			test:    "some_condition",
			val:     false,
			want:    false,
			wantErr: false,
		},
		{
			name: "bool_condition",
			args: args{
				opts: []ConditionOptions{
					{Name: "mycond", Type: "bool", Options: map[string]interface{}{
						"value": true,
					}},
				},
			},
			test:    "mycond",
			val:     true,
			want:    true,
			wantErr: false,
		},
		{
			name: "ip_whitelist",
			args: args{
				opts: []ConditionOptions{
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
			},
			test:    "office-ip",
			val:     "192.168.10.111",
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConditions(tt.args.opts, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			spew.Dump(got)
			if tt.want != got[tt.test].Meets(tt.val, nil) {
				t.Errorf("Condition.Meets() = %v, want %v", got, tt.want)
			}
		})
	}
}
