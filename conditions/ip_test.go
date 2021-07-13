package conditions

import (
	"testing"

	"github.com/blushft/redtape"
	"github.com/davecgh/go-spew/spew"
)

func TestIPConditions(t *testing.T) {
	reg := redtape.NewConditionRegistry(map[string]redtape.ConditionBuilder{
		new(IPAllowCondition).Name(): func() redtape.Condition {
			return new(IPAllowCondition)
		},
		new(IPDenyCondition).Name(): func() redtape.Condition {
			return new(IPDenyCondition)
		},
	})

	type args struct {
		opts []redtape.ConditionOptions
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
			name: "ip_whitelist",
			args: args{
				opts: []redtape.ConditionOptions{
					{
						Name: "office-ip",
						Type: "ip_allow",
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
			got, err := redtape.NewConditions(tt.args.opts, reg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConditions() error = %v, wantErr %v", err, tt.wantErr)
				t.FailNow()
			}

			spew.Dump(got)
			if tt.want != got[tt.test].Meets(tt.val, nil) {
				t.Errorf("Condition.Meets() = %v, want %v", got, tt.want)
			}
		})
	}
}
