package utils

import "testing"

func Test_trimSpacesAndQuotes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no spaces or quotes", args{"hello"}, "hello"},
		{"leading and trailing spaces", args{"  hello  "}, "hello"},
		{"single quotes", args{"'hello'"}, "hello"},
		{"double quotes", args{`"hello"`}, "hello"},
		{"mixed spaces and single quotes", args{"  'hello'  "}, "hello"},
		{"mixed spaces and double quotes", args{`  "hello"  `}, "hello"},
		{"no quotes but spaces inside", args{"  hello world  "}, "hello world"},
		{"quotes inside but not around", args{`  he"llo'  `}, `he"llo'`},
		{"empty string", args{""}, ""},
		{"only spaces", args{"     "}, ""},
		{"only single quotes", args{"''"}, ""},
		{"only double quotes", args{`""`}, ""},
		{"single quote at start", args{"'hello"}, "'hello"},
		{"single quote at end", args{"hello'"}, "hello'"},
		{"double quote at start", args{`"hello`}, `"hello`},
		{"double quote at end", args{`hello"`}, `hello"`},
		{"nested quotes", args{`"'hello'"`}, `'hello'`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := TrimSpacesAndQuotes(tt.args.s); got != tt.want {
				t.Errorf("trimSpacesAndQuotes() = %v, want %v", got, tt.want)
			}
		})
	}
}
