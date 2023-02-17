package parameterResolver

import (
	"reflect"
	"testing"

	"github.com/seboste/sapper/ports"
)

type TestResolver struct {
	ReturnValue string
}

func (t TestResolver) Resolve(name string, defaultValue string) string {
	return t.ReturnValue
}

var resolverA TestResolver = TestResolver{ReturnValue: "A"}
var resolverB TestResolver = TestResolver{ReturnValue: "B"}
var resolverEmpty TestResolver = TestResolver{ReturnValue: ""}

func TestMakeCompoundParameterResolver(t *testing.T) {
	type args struct {
		resolver []ports.ParameterResolver
	}
	tests := []struct {
		name string
		args args
		want CompoundParameterResolver
	}{
		{
			name: "single Resolver",
			args: args{resolver: []ports.ParameterResolver{resolverA}},
			want: CompoundParameterResolver{resolver: []ports.ParameterResolver{resolverA}},
		},
		{
			name: "multiple Resolver",
			args: args{resolver: []ports.ParameterResolver{resolverA, resolverB}},
			want: CompoundParameterResolver{resolver: []ports.ParameterResolver{resolverA, resolverB}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeCompoundParameterResolver(tt.args.resolver); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeCompoundParameterResolver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompoundParameterResolver_Resolve(t *testing.T) {
	type fields struct {
		resolver []ports.ParameterResolver
	}
	type args struct {
		name         string
		defaultValue string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "empty",
			fields: fields{},
			args:   args{},
			want:   "",
		},
		{
			name:   "just a",
			fields: fields{resolver: []ports.ParameterResolver{resolverA}},
			args:   args{},
			want:   "A",
		},
		{
			name:   "a before b",
			fields: fields{resolver: []ports.ParameterResolver{resolverA, resolverB}},
			args:   args{},
			want:   "A",
		},
		{
			name:   "a after b",
			fields: fields{resolver: []ports.ParameterResolver{resolverB, resolverA}},
			args:   args{},
			want:   "B",
		},
		{
			name:   "empty before a",
			fields: fields{resolver: []ports.ParameterResolver{resolverEmpty, resolverA}},
			args:   args{},
			want:   "A",
		},
		{
			name:   "empty after a",
			fields: fields{resolver: []ports.ParameterResolver{resolverA, resolverEmpty}},
			args:   args{},
			want:   "A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpr := CompoundParameterResolver{
				resolver: tt.fields.resolver,
			}
			if got := cpr.Resolve(tt.args.name, tt.args.defaultValue); got != tt.want {
				t.Errorf("CompoundParameterResolver.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
