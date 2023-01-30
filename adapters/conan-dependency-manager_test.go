package adapters

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/seboste/sapper/ports"
)

func TestConanDependencyManager_Read(t *testing.T) {
	serviceDir, err := ioutil.TempDir("", "service")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(serviceDir) // clean up
	testService := ports.Service{Id: "test", Path: serviceDir}

	type args struct {
		s ports.Service
	}
	tests := []struct {
		name      string
		conanfile string
		args      args
		want      []ports.PackageDependency
		wantErr   bool
	}{
		{name: "no conanfile", conanfile: "", args: args{s: testService}, want: []ports.PackageDependency{}, wantErr: true},
		{name: "no dependencies", conanfile: `[some section]
a
b
c
`, args: args{s: testService}, want: []ports.PackageDependency{}, wantErr: false},
		{name: "single dependency", conanfile: `[some section]
a
[requires]
my_lib/1.2.3
b
c
`, args: args{s: testService}, want: []ports.PackageDependency{{Id: "my_lib", Version: "1.2.3"}}, wantErr: false},
		{name: "dependency outside section", conanfile: `[some section]
a
bla/2.3.4
b
[requires]
my_lib/1.2.3
b
c
	`, args: args{s: testService}, want: []ports.PackageDependency{{Id: "my_lib", Version: "1.2.3"}}, wantErr: false},
		{name: "multiple dependencies", conanfile: `
[requires]
my_lib/1.2.3
		dep1/bla	
	dep2/1.2.3b
invalid1 /0.0.1
invalid2/ 0.0.1
dep3/2.3.4
invalid3
	`, args: args{s: testService}, want: []ports.PackageDependency{{Id: "my_lib", Version: "1.2.3"}, {Id: "dep1", Version: "bla"}, {Id: "dep2", Version: "1.2.3b"}, {Id: "dep3", Version: "2.3.4"}}, wantErr: false},
		{name: "disabled dependencies", conanfile: `
[requires]
#dep1/v1
# dep2/v1
 #dep3/v1
dep4/v1
#dep4/v2
	`, args: args{s: testService}, want: []ports.PackageDependency{{Id: "dep4", Version: "v1"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//1. prepare test
			cdm := ConanDependencyManager{}

			conanfilePath := filepath.Join(serviceDir, "conanfile.txt")
			os.Remove(conanfilePath)
			if tt.conanfile != "" {
				ioutil.WriteFile(conanfilePath, []byte(tt.conanfile), 0666)
			}

			//2. execute test
			got, err := cdm.Read(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConanDependencyManager.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConanDependencyManager.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseConanDependency(t *testing.T) {
	type args struct {
		input string
	}

	tests := []struct {
		name    string
		args    args
		want    ConanDependency
		wantErr bool
	}{
		{name: "lib, version only", args: args{input: "lib/version"}, want: ConanDependency{Id: "lib", Version: "version"}, wantErr: false},
		{name: "terminating @", args: args{input: "lib/version@"}, want: ConanDependency{Id: "lib", Version: "version"}, wantErr: false},
		{name: "with user and channel", args: args{input: "lib/version@user/channel"}, want: ConanDependency{Id: "lib", Version: "version", User: "user", Channel: "channel"}, wantErr: false},
		{name: "with reference", args: args{input: "lib/version#abcdef0123456789"}, want: ConanDependency{Id: "lib", Version: "version", Reference: "abcdef0123456789"}, wantErr: false},
		{name: "full blown", args: args{input: "lib/version@user/channel#abcdef0123456789"}, want: ConanDependency{Id: "lib", Version: "version", User: "user", Channel: "channel", Reference: "abcdef0123456789"}, wantErr: false},
		{name: "invalid syntax", args: args{input: "lib_without_version"}, want: ConanDependency{}, wantErr: true},
		{name: "invalid reference", args: args{input: "lib/version#invalid_reference"}, want: ConanDependency{Id: "lib", Version: "version"}, wantErr: false}, //TODO: error instead of empty reference?
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConanDependency(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConanDependency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConanDependency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConanDependencyManager_Write(t *testing.T) {
	serviceDir, err := ioutil.TempDir("", "service")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(serviceDir) // clean up
	testService := ports.Service{Id: "test", Path: serviceDir}

	type args struct {
		s          ports.Service
		dependency string
		version    string
	}
	tests := []struct {
		name          string
		conanfile     string
		cdm           ConanDependencyManager
		args          args
		wantConanfile string
		wantErr       bool
	}{
		{name: "single dependency", cdm: ConanDependencyManager{}, conanfile: `
[requires]
mylib/1.2.3
`, args: args{s: testService, dependency: "mylib", version: "1.2.4"}, wantErr: false, wantConanfile: `
[requires]
mylib/1.2.4
`,
		},
		{name: "full blown dependency", cdm: ConanDependencyManager{}, conanfile: `
[requires]
lib1/0.1.2
abc mylib/1.2.3@user/channel#01234ebc	xyz
lib2/0.1.2
`, args: args{s: testService, dependency: "mylib", version: "1.2.4"}, wantErr: false, wantConanfile: `
[requires]
lib1/0.1.2
abc mylib/1.2.4@user/channel#01234ebc	xyz
lib2/0.1.2
`,
		},
		{name: "lib duplicate error", cdm: ConanDependencyManager{}, conanfile: `
[requires]
mylib/1.2.3
mylib/1.2.2
`, args: args{s: testService, dependency: "mylib", version: "1.2.4"}, wantErr: true, wantConanfile: `
[requires]
mylib/1.2.3
mylib/1.2.2
`,
		},
		{name: "disabled lib error", cdm: ConanDependencyManager{}, conanfile: `
		[requires]
		#mylib/1.2.3
		`, args: args{s: testService, dependency: "mylib", version: "1.2.4"}, wantErr: true, wantConanfile: `
		[requires]
		#mylib/1.2.3
		`,
		},
		{name: "disabled lib preserved", cdm: ConanDependencyManager{}, conanfile: `
[requires]
#mylib/1.2.3
lib/v1
#this comment needs to be preserved as well
`, args: args{s: testService, dependency: "lib", version: "v2"}, wantErr: false, wantConanfile: `
[requires]
#mylib/1.2.3
lib/v2
#this comment needs to be preserved as well
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//1. prepare test
			cdm := ConanDependencyManager{}

			conanfilePath := filepath.Join(serviceDir, "conanfile.txt")
			os.Remove(conanfilePath)
			if tt.conanfile != "" {
				ioutil.WriteFile(conanfilePath, []byte(tt.conanfile), 0666)
			}

			//2. execute test
			if err := cdm.Write(tt.args.s, tt.args.dependency, tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("ConanDependencyManager.Write() error = %v, wantErr %v", err, tt.wantErr)
			}

			content, err := ioutil.ReadFile(conanfilePath)
			if err != nil {
				t.Errorf("ConanDependencyManager.Write() error. Unable to read from conanfile")
			}
			if string(content) != tt.wantConanfile {
				t.Errorf("ConanDependencyManager.Write() content = %v, wantContent %v", string(content), tt.wantConanfile)
			}

		})
	}
}

func TestConanDependency_String(t *testing.T) {
	type fields struct {
		Id        string
		Version   string
		User      string
		Channel   string
		Reference string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Basic", fields: fields{Id: "lib", Version: "version"}, want: "lib/version"},
		{name: "User and Channel", fields: fields{Id: "lib", Version: "version", User: "user", Channel: "channel"}, want: "lib/version@user/channel"},
		{name: "Reference", fields: fields{Id: "lib", Version: "version", Reference: "123abc"}, want: "lib/version#123abc"},
		{name: "Full Blown", fields: fields{Id: "lib", Version: "version", User: "user", Channel: "channel", Reference: "123abc"}, want: "lib/version@user/channel#123abc"},
		{name: "Missing Channel", fields: fields{Id: "lib", Version: "version", User: "user"}, want: "lib/version@user/_"},
		{name: "Missing User", fields: fields{Id: "lib", Version: "version", Channel: "channel"}, want: "lib/version@_/channel"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := ConanDependency{
				Id:        tt.fields.Id,
				Version:   tt.fields.Version,
				User:      tt.fields.User,
				Channel:   tt.fields.Channel,
				Reference: tt.fields.Reference,
			}
			if got := dep.String(); got != tt.want {
				t.Errorf("ConanDependency.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
