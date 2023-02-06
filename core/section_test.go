package core

import (
	"reflect"
	"testing"
)

func Test_readSections(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    []section
		wantErr bool
	}{
		{name: "single section", args: args{data: `0
1
2  xyz <<<SAPPER SECTION BEGIN APPEND SOME-NAME>>> abc
3content1
4content2
5content3
6<<<SAPPER SECTION END APPEND SOME-NAME>>>`},
			want: []section{{
				name:      "SOME-NAME",
				verb:      "APPEND",
				lineBegin: 3,
				lineEnd:   6,
				content: `3content1
4content2
5content3`,
			},
			},
			wantErr: false,
		},
		{name: "two sections", args: args{data: `0
1
2  xyz <<<SAPPER SECTION BEGIN APPEND SOME-NAME>>> abc
3content1
4content2
5content3
6<<<SAPPER SECTION END APPEND SOME-NAME>>>
7<<<SAPPER SECTION BEGIN BLA>>>
8<<<SAPPER SECTION END BLA>>>bla
9
`},
			want: []section{{
				name:      "SOME-NAME",
				verb:      "APPEND",
				lineBegin: 3,
				lineEnd:   6,
				content: `3content1
4content2
5content3`,
			}, {
				name:      "BLA",
				verb:      "",
				lineBegin: 8,
				lineEnd:   8,
				content:   "",
			},
			},
			wantErr: false,
		},
		{
			name:    "end tag without begin tag",
			args:    args{data: `<<<SAPPER SECTION END SOME-NAME>>>`},
			want:    []section{},
			wantErr: true,
		},
		{
			name: "intersecting sections",
			args: args{data: `
<<<SAPPER SECTION BEGIN SOME-NAME>>>
<<<SAPPER SECTION BEGIN SOME-OTHER-NAME>>>
<<<SAPPER SECTION END SOME-NAME>>>
<<<SAPPER SECTION END SOME-OTHER-NAME>>>`},
			want:    []section{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readSections(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("readSections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readSections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readTag(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want *tag
	}{
		{
			name: "begin tag",
			args: args{line: "#<<<SAPPER SECTION BEGIN MY-SECTION>>>"},
			want: &tag{name: "MY-SECTION", verb: "", begin: true},
		},
		{
			name: "end tag",
			args: args{line: "#<<<SAPPER SECTION END MY-SECTION>>>"},
			want: &tag{name: "MY-SECTION", verb: "", begin: false},
		},
		{
			name: "end append tag",
			args: args{line: "#<<<SAPPER SECTION END APPEND MY-SECTION>>>"},
			want: &tag{name: "MY-SECTION", verb: "APPEND", begin: false},
		},
		{
			name: "begin replace tag",
			args: args{line: "#<<<SAPPER SECTION BEGIN REPLACE MY-SECTION>>>"},
			want: &tag{name: "MY-SECTION", verb: "REPLACE", begin: true},
		},
		{
			name: "whitespace",
			args: args{line: "#<<<SAPPER 		SECTION  BEGIN	 	MY-SECTION>>>"},
			want: &tag{name: "MY-SECTION", verb: "", begin: true},
		},
		{
			name: "begin tag embedded in stuff",
			args: args{line: "some stuff     <<<SAPPER SECTION BEGIN MY-SECTION>>>>>> some more stuff"},
			want: &tag{name: "MY-SECTION", verb: "", begin: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readTag(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCurrentSection(t *testing.T) {
	type args struct {
		line    string
		section string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "entering section", args: args{line: "<<<SAPPER SECTION BEGIN new>>>", section: ""}, want: "new"},
		{name: "inside section", args: args{line: "bla bla", section: "current"}, want: "current"},
		{name: "leaving section", args: args{line: "<<<SAPPER SECTION END current>>>", section: "current"}, want: ""},
		{name: "nested section", args: args{line: "<<<SAPPER SECTION BEGIN new>>>", section: "current"}, want: "current"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCurrentSection(tt.args.line, tt.args.section); got != tt.want {
				t.Errorf("getCurrentSection() = %v, want %v", got, tt.want)
			}
		})
	}
}
