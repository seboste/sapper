package core

import (
	"reflect"
	"sort"
	"testing"
)

func TestSemanticVersion_String(t *testing.T) {
	type fields struct {
		Prefix string
		Major  uint32
		Minor  uint32
		Patch  uint32
		Suffix string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Simple", fields: fields{Major: 1, Minor: 2, Patch: 3}, want: "1.2.3"},
		{name: "WithSuffix", fields: fields{Major: 1, Minor: 2, Patch: 3, Suffix: "suffix"}, want: "1.2.3suffix"},
		{name: "WithPrefix", fields: fields{Major: 1, Minor: 2, Patch: 3, Prefix: "prefix"}, want: "prefix1.2.3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := SemanticVersion{
				Prefix: tt.fields.Prefix,
				Major:  tt.fields.Major,
				Minor:  tt.fields.Minor,
				Patch:  tt.fields.Patch,
				Suffix: tt.fields.Suffix,
			}
			if got := v.String(); got != tt.want {
				t.Errorf("SemanticVersion.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSemanticVersion(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    SemanticVersion
		wantErr bool
	}{
		{name: "full blown", args: args{s: "prefix1.2.3suffix"}, want: SemanticVersion{Prefix: "prefix", Major: 1, Minor: 2, Patch: 3, Suffix: "suffix"}, wantErr: false},
		{name: "large version", args: args{s: "123.456.789"}, want: SemanticVersion{Major: 123, Minor: 456, Patch: 789}, wantErr: false},
		{name: "invalid version", args: args{s: "123.456"}, want: SemanticVersion{}, wantErr: true},
		{name: "invalid version with non digits", args: args{s: "123.4S6.789"}, want: SemanticVersion{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSemanticVersion(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSemanticVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSemanticVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func SemVer(s string) SemanticVersion {
	v, _ := ParseSemanticVersion(s)
	return v
}

func SemVerPtr(s string) *SemanticVersion {
	v, err := ParseSemanticVersion(s)
	if err != nil {
		return nil
	}
	return &v
}

func Test_ConvertToSemVer(t *testing.T) {
	type args struct {
		versions []string
	}
	tests := []struct {
		name    string
		args    args
		want    []SemanticVersion
		wantErr bool
	}{
		{name: "multiple entries",
			args:    args{versions: []string{"1.2.3", "1.2.4", "0.1.1"}},
			want:    []SemanticVersion{SemVer("1.2.3"), SemVer("1.2.4"), SemVer("0.1.1")},
			wantErr: false,
		},
		{name: "invalid entry",
			args:    args{versions: []string{"1.2.3", "1.invalid.4", "0.1.1"}},
			want:    []SemanticVersion{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToSemVer(tt.args.versions)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToSemVer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToSemVer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByVersion_Less(t *testing.T) {
	testData := []SemanticVersion{
		SemVer("1.1.1"),
		SemVer("0.1.1"),
		SemVer("1.0.1"),
		SemVer("1.1.0"),
		SemVer("0.99.99"),
		SemVer("1.1.99"),
		SemVer("1.1.1abc"),
		SemVer("1.1.1def"),
		SemVer("prefix1.1.1"),
	}
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		v    ByVersion
		args args
		want bool
	}{
		{name: "Identical", v: testData, args: args{i: 0, j: 0}, want: false},
		{name: "Differ by Major", v: testData, args: args{i: 0, j: 1}, want: false},
		{name: "Differ by Major inverse", v: testData, args: args{i: 1, j: 0}, want: true},
		{name: "Differ by Minor", v: testData, args: args{i: 0, j: 2}, want: false},
		{name: "Differ by Minor inverse", v: testData, args: args{i: 2, j: 0}, want: true},
		{name: "Differ by Patch", v: testData, args: args{i: 0, j: 3}, want: false},
		{name: "Differ by Patch inverse", v: testData, args: args{i: 3, j: 0}, want: true},
		{name: "Differ by Major with larger Minor", v: testData, args: args{i: 0, j: 4}, want: false},
		{name: "Differ by Major with larger Minor inverse", v: testData, args: args{i: 4, j: 0}, want: true},
		{name: "Differ by Patch (large)", v: testData, args: args{i: 0, j: 5}, want: true},
		{name: "Differ by Patch (large) inverse", v: testData, args: args{i: 5, j: 0}, want: false},
		{name: "Differ by Suffix", v: testData, args: args{i: 0, j: 6}, want: true},
		{name: "Differ by Suffix inverse", v: testData, args: args{i: 6, j: 0}, want: false},
		{name: "Differ by Suffix 2", v: testData, args: args{i: 6, j: 7}, want: true},
		{name: "Differ by Suffix 2 inverse", v: testData, args: args{i: 7, j: 6}, want: false},
		{name: "Differ by Prefix", v: testData, args: args{i: 0, j: 8}, want: false},
		{name: "Differ by Prefix inverse", v: testData, args: args{i: 8, j: 0}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("ByVersion.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByVersion_Sort(t *testing.T) {
	type args struct {
		v []SemanticVersion
	}
	tests := []struct {
		name string
		args args
		want []SemanticVersion
	}{
		{name: "Some Random Order", args: args{v: []SemanticVersion{

			SemVer("1.0.1"),
			SemVer("1.2.5aa"),
			SemVer("1.0.0"),
			SemVer("v2.0.0ac"),
			SemVer("1.2.5a"),
			SemVer("2.0.0"),
			SemVer("1.2.0"),
			SemVer("1.2.5"),
			SemVer("1.0.2"),
			SemVer("1.2.5b"),
			SemVer("1.2.0"),
			SemVer("2.0.0ab"),
			SemVer("1.0.10"),
		}}, want: []SemanticVersion{
			SemVer("1.0.0"),
			SemVer("1.0.1"),
			SemVer("1.0.2"),
			SemVer("1.0.10"),
			SemVer("1.2.0"),
			SemVer("1.2.0"),
			SemVer("1.2.5"),
			SemVer("1.2.5a"),
			SemVer("1.2.5aa"),
			SemVer("1.2.5b"),
			SemVer("2.0.0"),
			SemVer("2.0.0ab"),
			SemVer("v2.0.0ac"),
		},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := tt.args.v
			sort.Sort(ByVersion(values))
			if !reflect.DeepEqual(values, tt.want) {
				t.Errorf("sort.Sort(ByVersion) = %v, want %v", values, tt.want)
			}
		})
	}
}
