package core

import (
	"testing"
)

func Test_replaceParameters(t *testing.T) {
	type args struct {
		content    string
		parameters map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "single line", args: args{content: "<<<BLA>>>", parameters: map[string]string{"BLA": "XY"}}, want: "XY"},
		{name: "single line with context", args: args{content: "abc<<<BLA>>>def", parameters: map[string]string{"BLA": "XY"}}, want: "abcXYdef"},
		{name: "single line multiple params", args: args{content: "<<<BLA>>><<<BLUB>>>", parameters: map[string]string{"BLA": "XY", "BLUB": "AB"}}, want: "XYAB"},
		{name: "single line multiple occurences", args: args{content: "<<<BLA>>>abc<<<BLA>>>", parameters: map[string]string{"BLA": "XY"}}, want: "XYabcXY"},
		{name: "single line undefined parameter", args: args{content: "abc<<<UNDEFINED>>>def", parameters: map[string]string{"BLA": "XY"}}, want: "abc<<<UNDEFINED>>>def"},
		{name: "multi line", args: args{content: `test
bla
<<<BLA>>>
bla
`, parameters: map[string]string{"BLA": "XY"}}, want: `test
bla
XY
bla
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceParameters(tt.args.content, tt.args.parameters); got != tt.want {
				t.Errorf("replaceParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeSection(t *testing.T) {
	type args struct {
		base     section
		incoming section
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "append", args: args{base: section{content: "a"}, incoming: section{content: "b", verb: "APPEND"}}, want: "a\nb", wantErr: false},
		{name: "prepend", args: args{base: section{content: "a"}, incoming: section{content: "b", verb: "PREPEND"}}, want: "b\na", wantErr: false},
		{name: "replace", args: args{base: section{content: "a"}, incoming: section{content: "b", verb: "REPLACE"}}, want: "b", wantErr: false},
		{name: "append empty b", args: args{base: section{content: "a"}, incoming: section{content: "", verb: "APPEND"}}, want: "a", wantErr: false},
		{name: "append empty a", args: args{base: section{content: ""}, incoming: section{content: "b", verb: "APPEND"}}, want: "b", wantErr: false},
		{name: "prepend empty b", args: args{base: section{content: "a"}, incoming: section{content: "", verb: "PREPEND"}}, want: "a", wantErr: false},
		{name: "prepend empty a", args: args{base: section{content: ""}, incoming: section{content: "b", verb: "PREPEND"}}, want: "b", wantErr: false},
		{name: "error verb a", args: args{base: section{content: "a", verb: "APPEND"}, incoming: section{content: "b", verb: "APPEND"}}, want: "", wantErr: true},
		{name: "error no verb b", args: args{base: section{content: "a"}, incoming: section{content: "b", verb: ""}}, want: "", wantErr: true},
		{name: "error different names", args: args{base: section{content: "a", name: "SECTION-A"}, incoming: section{content: "b", verb: "APPEND", name: "SECTION-B"}}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeSection(tt.args.base, tt.args.incoming)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mergeSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeSections(t *testing.T) {
	type args struct {
		content       string
		inputSections map[string]section
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "no sections", args: args{content: `0
1
2
3
`}, want: `0
1
2
3
`, wantErr: false},
		{name: "single empty section", args: args{content: `0
1<<<SAPPER SECTION BEGIN BLA>>>
2<<<SAPPER SECTION END BLA>>>
3
`}, want: `0
1<<<SAPPER SECTION BEGIN BLA>>>
2<<<SAPPER SECTION END BLA>>>
3
`, wantErr: false},
		{name: "single section append", args: args{inputSections: map[string]section{"BLA": {name: "BLA", content: "b", verb: "APPEND"}},
			content: `0
1<<<SAPPER SECTION BEGIN BLA>>>
a
2<<<SAPPER SECTION END BLA>>>
3
`}, want: `0
1<<<SAPPER SECTION BEGIN BLA>>>
a
b
2<<<SAPPER SECTION END BLA>>>
3
`, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeSections(tt.args.content, tt.args.inputSections)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeSections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mergeSections() = %v, want %v", got, tt.want)
			}
		})
	}
}
