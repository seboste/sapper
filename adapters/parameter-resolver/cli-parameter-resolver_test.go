package parameterResolver

import (
	"reflect"
	"testing"
)

func TestCommandLineInterfaceParameterResolver_Resolve(t *testing.T) {
	type fields struct {
		parameters map[string]string
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
		{name: "existing param retruns value", fields: fields{parameters: map[string]string{"param": "value"}}, args: args{name: "param", defaultValue: ""}, want: "value"},
		{name: "existing param with default returns value", fields: fields{parameters: map[string]string{"param": "value"}}, args: args{name: "param", defaultValue: "default"}, want: "value"},
		{name: "unknown param returns emply", fields: fields{parameters: map[string]string{"param": "value"}}, args: args{name: "unknown", defaultValue: ""}, want: ""},
		{name: "unknown param with default returns emply", fields: fields{parameters: map[string]string{"param": "value"}}, args: args{name: "unknown", defaultValue: "default"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clipr := CommandLineInterfaceParameterResolver{
				parameters: tt.fields.parameters,
			}
			if got := clipr.Resolve(tt.args.name, tt.args.defaultValue); got != tt.want {
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
