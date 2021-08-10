package match

import "testing"

func TestSimpleWildcard(t *testing.T) {
	type args struct {
		search string
		val    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test 1",
			args: args{
				search: "test*",
				val:    "test_string",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SimpleWildcard(tt.args.search, tt.args.val); got != tt.want {
				t.Errorf("SimpleWildcard() = %v, want %v", got, tt.want)
			}
		})
	}
}
