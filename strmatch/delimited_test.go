package strmatch

import (
	"reflect"
	"testing"
)

func TestExtractDelimited(t *testing.T) {
	type args struct {
		s          string
		delimStart rune
		delimEnd   rune
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "single match",
			args: args{
				s:          "foo.bar.<.*>",
				delimStart: '<',
				delimEnd:   '>',
			},
			want: []string{
				".*",
			},
			wantErr: false,
		},
		{
			name: "two matches",
			args: args{
				s:          "foo.bar.<.*>.<[ABC]>",
				delimStart: '<',
				delimEnd:   '>',
			},
			want: []string{
				".*",
				"[ABC]",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractDelimited(tt.args.s, tt.args.delimStart, tt.args.delimEnd)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractDelimited() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchDelimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompileDelimitedRegex(t *testing.T) {
	type args struct {
		s          string
		delimStart rune
		delimEnd   rune
	}
	tests := []struct {
		name    string
		args    args
		match   string
		want    bool
		wantErr bool
	}{
		{
			name: "test match",
			args: args{
				s:          "foo.bar.<.*>",
				delimStart: '<',
				delimEnd:   '>',
			},
			match:   "foo.bar.test",
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg, err := CompileDelimitedRegex(tt.args.s, tt.args.delimStart, tt.args.delimEnd)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompileDelimitedRegex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got := reg.MatchString(tt.match); got != tt.want {
				t.Errorf("CompileDelimitedRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}
