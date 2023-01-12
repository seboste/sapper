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
		{name: "user, channel", conanfile: `[requires]
my_lib/1.2.3@user/channel
	`, args: args{s: testService}, want: []ports.PackageDependency{{Id: "my_lib", Version: "1.2.3"}}, wantErr: false},
		{name: "terminating @", conanfile: `[requires]
my_lib/1.2.3@
	`, args: args{s: testService}, want: []ports.PackageDependency{{Id: "my_lib", Version: "1.2.3"}}, wantErr: false},
		{name: "reference", conanfile: `[requires]
my_lib/1.2.3@user/channel#6af9cc7cb931c5ad94
	`, args: args{s: testService}, want: []ports.PackageDependency{{Id: "my_lib", Version: "1.2.3"}}, wantErr: false},
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
