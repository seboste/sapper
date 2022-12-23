package adapters

import (
	"reflect"
	"testing"
)

func TestCommandLineInterfaceParameterResolver_Resolve(t *testing.T) {
	type fields struct {
		parameters map[string]string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "existing param", fields: fields{parameters: map[string]string{"param": "value"}}, args: args{name: "param"}, want: "value"},
		{name: "unknown param", fields: fields{parameters: map[string]string{"param": "value"}}, args: args{name: "unknown"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clipr := CommandLineInterfaceParameterResolver{
				parameters: tt.fields.parameters,
			}
			if got := clipr.Resolve(tt.args.name); got != tt.want {
				t.Errorf("CommandLineInterfaceParameterResolver.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeCommandLineInterfaceParameterResolver(t *testing.T) {
	type args struct {
		parameters []string
	}
	tests := []struct {
		name    string
		args    args
		want    CommandLineInterfaceParameterResolver
		wantErr bool
	}{
		{name: "single param", args: args{parameters: []string{"param=value"}}, want: CommandLineInterfaceParameterResolver{parameters: map[string]string{"param": "value"}}, wantErr: false},
		{name: "multiple param", args: args{parameters: []string{"param1=value1", "param2=value2"}}, want: CommandLineInterfaceParameterResolver{parameters: map[string]string{"param1": "value1", "param2": "value2"}}, wantErr: false},
		{name: "invalid param syntax", args: args{parameters: []string{"abcdef"}}, want: CommandLineInterfaceParameterResolver{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeCommandLineInterfaceParameterResolver(tt.args.parameters)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeCommandLineInterfaceParameterResolver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeCommandLineInterfaceParameterResolver() = %v, want %v", got, tt.want)
			}
		})
	}
}
