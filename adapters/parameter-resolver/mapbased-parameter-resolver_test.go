package parameterResolver

import "testing"

func TestMapBasedParameterResolver_Resolve(t *testing.T) {
	type fields struct {
		parameters map[string]string
	}
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "unknown parameter",
			fields: fields{parameters: map[string]string{"a": "b"}},
			args:   args{key: "x"},
			want:   "",
		},
		{
			name:   "known parameter",
			fields: fields{parameters: map[string]string{"a": "b"}},
			args:   args{key: "a"},
			want:   "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := MapBasedParameterResolver{
				parameters: tt.fields.parameters,
			}
			if got := r.Resolve(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("MapBasedParameterResolver.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
