package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/seboste/sapper/ports"
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
		{name: "merge", args: args{base: section{content: "a"}, incoming: section{content: "b", verb: "MERGE"}}, want: "a\nb", wantErr: false},
		{name: "merge same", args: args{base: section{content: "a"}, incoming: section{content: "a", verb: "MERGE"}}, want: "a", wantErr: false},
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
		{name: "multiple sections", args: args{inputSections: map[string]section{"BLA": {name: "BLA", content: "b", verb: "REPLACE"}},
			content: `
//<<<SAPPER SECTION BEGIN BLUB>>>
//<<<SAPPER SECTION END BLUB>>>

//<<<SAPPER SECTION BEGIN WURST>>>
abc
//<<<SAPPER SECTION END WURST>>>

//<<<SAPPER SECTION BEGIN BLA>>>
a
//<<<SAPPER SECTION END BLA>>>

xyz
`}, want: `
//<<<SAPPER SECTION BEGIN BLUB>>>
//<<<SAPPER SECTION END BLUB>>>

//<<<SAPPER SECTION BEGIN WURST>>>
abc
//<<<SAPPER SECTION END WURST>>>

//<<<SAPPER SECTION BEGIN BLA>>>
b
//<<<SAPPER SECTION END BLA>>>

xyz
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

type testResolver struct{}

func (tr testResolver) Resolve(name string, defaultValue string) string {
	if name == "a" {
		return "1"
	}
	if name == "b" {
		return "2"
	}
	return defaultValue
}

func TestResolveParameters(t *testing.T) {
	test_resolver := testResolver{}
	type args struct {
		bp []ports.BrickParameters
		pr ports.ParameterResolver
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "custom and default",
			args:    args{bp: []ports.BrickParameters{{Name: "a", Default: "d1"}, {Name: "c", Default: "d2"}}, pr: test_resolver},
			want:    map[string]string{"a": "1", "c": "d2"},
			wantErr: false,
		},
		{
			name:    "no default available",
			args:    args{bp: []ports.BrickParameters{{Name: "c"}}, pr: test_resolver},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveParameters(tt.args.bp, tt.args.pr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveParameterSlice(t *testing.T) {
	test_resolver := testResolver{}
	type args struct {
		bricks []ports.Brick
		pr     ports.ParameterResolver
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "simple slice",
			args:    args{bricks: []ports.Brick{{Id: "1234", Parameters: []ports.BrickParameters{{Name: "a"}}}}, pr: test_resolver},
			want:    map[string]string{"a": "1"},
			wantErr: false,
		},
		{
			name:    "simple slice",
			args:    args{bricks: []ports.Brick{{Id: "1234", Parameters: []ports.BrickParameters{{Name: "c"}}}}, pr: test_resolver},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveParameterSlice(tt.args.bricks, tt.args.pr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveParameterSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveParameterSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddSingleBrick(t *testing.T) {

	brick1TempDir, _ := ioutil.TempDir("", "brick1")
	brick2TempDir, _ := ioutil.TempDir("", "brick2")
	serviceTempDir, _ := ioutil.TempDir("", "service")
	defer os.RemoveAll(brick1TempDir)  // clean up
	defer os.RemoveAll(brick2TempDir)  // clean up
	defer os.RemoveAll(serviceTempDir) // clean up

	ioutil.WriteFile(filepath.Join(brick1TempDir, "test.txt"), []byte(`this is some file
with some parameter 'bla' which has the value '<<<bla>>>'
`), 0666)

	ioutil.WriteFile(filepath.Join(serviceTempDir, "some_file_with_section.txt"), []byte(`this is some file
which has a section
<<<SAPPER SECTION BEGIN my_section>>>
with some content
<<<SAPPER SECTION END my_section>>>
and that's it.
`), 0666)

	ioutil.WriteFile(filepath.Join(brick2TempDir, "some_file_with_section.txt"), []byte(`some irrelevant content before the section
<<<SAPPER SECTION BEGIN APPEND my_section>>>
and even more content
<<<SAPPER SECTION END APPEND my_section>>>
`), 0666)

	type args struct {
		s          *ports.Service
		b          ports.Brick
		parameters map[string]string
	}
	tests := []struct {
		name        string
		args        args
		wantService ports.Service
		wantFiles   map[string]string //filename -> content
		wantErr     bool
	}{
		{
			name: "dependency",
			args: args{
				s:          &ports.Service{Id: "my_service", Dependencies: []ports.PackageDependency{}},
				b:          ports.Brick{Id: "b1", Version: "1.0.0"},
				parameters: map[string]string{},
			},
			wantService: ports.Service{
				Id:           "my_service",
				BrickIds:     []ports.BrickDependency{{Id: "b1", Version: "1.0.0"}},
				Dependencies: []ports.PackageDependency{},
				Parameters:   map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "copy file",
			args: args{
				s:          &ports.Service{Id: "my_service", Path: serviceTempDir, Dependencies: []ports.PackageDependency{}},
				b:          ports.Brick{Id: "b1", Version: "1.0.0", BasePath: brick1TempDir, Files: []string{"test.txt"}},
				parameters: map[string]string{"bla": "the_bla_value"},
			},
			wantService: ports.Service{
				Id:           "my_service",
				Path:         serviceTempDir,
				BrickIds:     []ports.BrickDependency{{Id: "b1", Version: "1.0.0"}},
				Dependencies: []ports.PackageDependency{},
				Parameters:   map[string]string{"bla": "the_bla_value"},
			},
			wantFiles: map[string]string{
				"test.txt": `this is some file
with some parameter 'bla' which has the value 'the_bla_value'
`},
			wantErr: false,
		},
		{
			name: "merge sections",
			args: args{
				s:          &ports.Service{Id: "my_service", Path: serviceTempDir, Dependencies: []ports.PackageDependency{}},
				b:          ports.Brick{Id: "b1", Version: "1.0.0", BasePath: brick2TempDir, Files: []string{"some_file_with_section.txt"}},
				parameters: map[string]string{},
			},
			wantService: ports.Service{
				Id:           "my_service",
				Path:         serviceTempDir,
				BrickIds:     []ports.BrickDependency{{Id: "b1", Version: "1.0.0"}},
				Dependencies: []ports.PackageDependency{},
				Parameters:   map[string]string{},
			},
			wantFiles: map[string]string{"some_file_with_section.txt": `this is some file
which has a section
<<<SAPPER SECTION BEGIN my_section>>>
with some content
and even more content
<<<SAPPER SECTION END my_section>>>
and that's it.
`},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddSingleBrick(tt.args.s, tt.args.b, tt.args.parameters); (err != nil) != tt.wantErr {
				t.Errorf("AddSingleBrick() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(*tt.args.s, tt.wantService) {
				t.Errorf("AddSingleBrick() service = %v, wantService %v", *tt.args.s, tt.wantService)
			}

			for filename, wantContent := range tt.wantFiles {
				content, err := ioutil.ReadFile(filepath.Join(serviceTempDir, filename))
				if err != nil {
					t.Errorf("Unable to read from %s", filename)
				}
				contentStr := string(content)
				if contentStr != wantContent {
					t.Errorf("AddSingleBrick() file %s = %s, wantFile %s", filename, contentStr, wantContent)
				}
			}
		})
	}
}

type TestBrickDB struct {
	initCalled   *bool
	updateCalled *bool
}

func (db *TestBrickDB) Init(Path string) error {
	if db.initCalled != nil {
		*db.initCalled = true
	}
	return nil
}

func (db TestBrickDB) Bricks(kind ports.BrickKind) []ports.Brick {
	return []ports.Brick{}
}

func (db TestBrickDB) Brick(id string) (ports.Brick, error) {
	if id == "brick1" {
		return ports.Brick{Id: "brick1", Dependencies: []string{"brick2", "brick3"}}, nil
	}
	if id == "brick2" {
		return ports.Brick{Id: "brick2", Dependencies: []string{"brick4"}}, nil
	}
	if id == "brick3" {
		return ports.Brick{Id: "brick3", Dependencies: []string{"brick4"}}, nil
	}
	if id == "brick4" {
		return ports.Brick{Id: "brick4", Dependencies: []string{}}, nil
	}
	if id == "brick5" {
		return ports.Brick{Id: "brick5", Dependencies: []string{"brick6"}}, nil
	}
	if id == "brick6" {
		return ports.Brick{Id: "brick6", Dependencies: []string{"brick5"}}, nil
	}
	return ports.Brick{}, fmt.Errorf("brick with id %s does not exist", id)
}

func (db *TestBrickDB) Update() error {
	if db.updateCalled != nil {
		*db.updateCalled = true
	}
	return nil
}

func (db TestBrickDB) IsModified() (bool, string) {
	return false, ""
}

var _ ports.BrickDB = (*TestBrickDB)(nil)

func TestGetBricksRecursive(t *testing.T) {
	type args struct {
		brickId string
		db      ports.BrickDB
	}
	tests := []struct {
		name    string
		args    args
		want    []ports.Brick
		wantErr bool
	}{
		{
			name:    "no_dependencies",
			args:    args{brickId: "brick4", db: &TestBrickDB{}},
			want:    []ports.Brick{{Id: "brick4", Dependencies: []string{}}},
			wantErr: false,
		},
		{
			name:    "single_dependency",
			args:    args{brickId: "brick2", db: &TestBrickDB{}},
			want:    []ports.Brick{{Id: "brick4", Dependencies: []string{}}, {Id: "brick2", Dependencies: []string{"brick4"}}},
			wantErr: false,
		},
		{
			name:    "diamond",
			args:    args{brickId: "brick1", db: &TestBrickDB{}},
			want:    []ports.Brick{{Id: "brick4", Dependencies: []string{}}, {Id: "brick2", Dependencies: []string{"brick4"}}, {Id: "brick3", Dependencies: []string{"brick4"}}, {Id: "brick1", Dependencies: []string{"brick2", "brick3"}}},
			wantErr: false,
		},
		{
			name:    "cycle",
			args:    args{brickId: "brick5", db: &TestBrickDB{}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBricksRecursive(tt.args.brickId, tt.args.db, map[string]bool{})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBricksRecursive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBricksRecursive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findLatestWorkingVersion(t *testing.T) {

	type args struct {
		versions  []SemanticVersion
		isWorking map[SemanticVersion]bool
	}
	tests := []struct {
		name string
		args args
		want SemanticVersion
	}{
		{name: "single not working version", args: args{versions: []SemanticVersion{SemVer("1.2.2"), SemVer("1.2.3")}}, want: SemVer("1.2.2")},
		{name: "single working version", args: args{versions: []SemanticVersion{SemVer("1.2.2"), SemVer("1.2.3")}, isWorking: map[SemanticVersion]bool{SemVer("1.2.3"): true}}, want: SemVer("1.2.3")},
		{name: "two working versions", args: args{versions: []SemanticVersion{SemVer("1.2.2"), SemVer("1.2.3"), SemVer("1.2.4")},
			isWorking: map[SemanticVersion]bool{SemVer("1.2.3"): true, SemVer("1.2.4"): true}},
			want: SemVer("1.2.4")},
		{name: "two versions: higher version working", args: args{versions: []SemanticVersion{SemVer("1.2.2"), SemVer("1.2.3"), SemVer("1.2.4")},
			isWorking: map[SemanticVersion]bool{SemVer("1.2.3"): false, SemVer("1.2.4"): true}},
			want: SemVer("1.2.4")},
		{name: "two versions: lower version working", args: args{versions: []SemanticVersion{SemVer("1.2.2"), SemVer("1.2.3"), SemVer("1.2.4")},
			isWorking: map[SemanticVersion]bool{SemVer("1.2.3"): true, SemVer("1.2.4"): false}},
			want: SemVer("1.2.3")},
		{name: "three versions: center version working", args: args{versions: []SemanticVersion{SemVer("1.2.2"), SemVer("1.2.3"), SemVer("1.2.4"), SemVer("1.2.5")},
			isWorking: map[SemanticVersion]bool{SemVer("1.2.3"): true, SemVer("1.2.4"): true, SemVer("1.2.5"): false}},
			want: SemVer("1.2.4")},
		{name: "five versions: 2nd version working", args: args{versions: []SemanticVersion{SemVer("1.2.2"), SemVer("1.2.3"), SemVer("1.2.4"), SemVer("1.2.5"), SemVer("1.2.6"), SemVer("1.2.7")},
			isWorking: map[SemanticVersion]bool{SemVer("1.2.3"): true, SemVer("1.2.4"): true, SemVer("1.2.5"): false, SemVer("1.2.6"): false, SemVer("1.2.7"): false}},
			want: SemVer("1.2.4")},
		{name: "five versions: 4th version working", args: args{versions: []SemanticVersion{SemVer("1.2.2"), SemVer("1.2.3"), SemVer("1.2.4"), SemVer("1.2.5"), SemVer("1.2.6"), SemVer("1.2.7")},
			isWorking: map[SemanticVersion]bool{SemVer("1.2.3"): true, SemVer("1.2.4"): true, SemVer("1.2.5"): true, SemVer("1.2.6"): true, SemVer("1.2.7"): false}},
			want: SemVer("1.2.6")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findLatestWorkingVersion(tt.args.versions, func(v SemanticVersion) bool {
				val, ok := tt.args.isWorking[v]
				return ok == true && val == true
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findLatestWorkingVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterSemvers(t *testing.T) {
	type args struct {
		in        []SemanticVersion
		predicate func(SemanticVersion) bool
	}
	tests := []struct {
		name string
		args args
		want []SemanticVersion
	}{
		{
			name: "two entries, keep major 1",
			args: args{
				in:        []SemanticVersion{SemVer("1.0.0"), SemVer("2.0.0")},
				predicate: func(v SemanticVersion) bool { return v.Major == 1 },
			},
			want: []SemanticVersion{SemVer("1.0.0")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterSemvers(tt.args.in, tt.args.predicate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterSemvers() = %v, want %v", got, tt.want)
			}
		})
	}
}

type versionBuildinfo struct {
	Version         string
	BuildSuccessful bool
}

type upgradeDependencyMock struct {
	AvailableVersionMap map[string]([]versionBuildinfo) //dependency->versions
	VersionMap          map[string]string               //dependency->version
}

func (udm upgradeDependencyMock) WriteToService(s ports.Service, d ports.PackageDependency) error {
	udm.VersionMap[d.Id] = d.Version
	return nil
}

func (udm upgradeDependencyMock) AvailableVersions(dependency string) ([]string, error) {
	availableVersions := []string{}
	for _, v := range udm.AvailableVersionMap[dependency] {
		availableVersions = append(availableVersions, v.Version)
	}
	return availableVersions, nil
}

func (udm upgradeDependencyMock) Load(path string) (ports.Service, error) {
	s := ports.Service{
		Id:   "my_service",
		Path: path,
	}
	for k, v := range udm.VersionMap {
		s.Dependencies = append(s.Dependencies, ports.PackageDependency{Id: k, Version: v})
	}
	return s, nil
}

func (udm upgradeDependencyMock) Save(service ports.Service) error {
	udm.VersionMap = make(map[string]string)
	for _, v := range service.Dependencies {
		udm.VersionMap[v.Id] = v.Version

	}
	return nil
}

func (udm upgradeDependencyMock) Build(service ports.Service, output io.Writer) error {
	for _, d := range service.Dependencies {
		dependencyFound := false
		for _, v := range udm.AvailableVersionMap[d.Id] {
			if v.Version == d.Version {
				if !v.BuildSuccessful {
					return fmt.Errorf("build with %s in version %s failed", d.Id, d.Version)
				} else {
					dependencyFound = true
				}
			}
		}
		if !dependencyFound {
			return fmt.Errorf("dependency %s not found", d.Id)
		}
	}
	return nil
}

func (udm upgradeDependencyMock) Test(service ports.Service, output io.Writer) error {
	return nil
}

func (udm upgradeDependencyMock) Run(service ports.Service, output io.Writer) error {
	return nil
}

func (udm upgradeDependencyMock) Deploy(service ports.Service, output io.Writer) error {
	return nil
}

func TestServiceApi_upgradeDependency(t *testing.T) {
	type fields struct {
		udm upgradeDependencyMock
	}
	type args struct {
		service          ports.Service
		d                ports.PackageDependency
		keepMajorVersion bool
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           VersionUpgradeSpec
		wantVersionMap map[string]string
		wantErr        bool
	}{
		{name: "non semantic version dependency up to date", fields: fields{
			udm: upgradeDependencyMock{
				AvailableVersionMap: map[string][]versionBuildinfo{"lib": {{Version: "version A", BuildSuccessful: true}}},
				VersionMap:          map[string]string{"lib": "version A"},
			},
		}, args: args{
			service:          ports.Service{Id: "my_service", Path: "path", Dependencies: []ports.PackageDependency{{Id: "lib", Version: "version A"}}},
			d:                ports.PackageDependency{Id: "lib", Version: "version A"},
			keepMajorVersion: false,
		}, want: VersionUpgradeSpec{
			previous:        "version A",
			target:          "version A",
			latestAvailable: "version A",
			latestWorking:   "version A",
		},
			wantVersionMap: map[string]string{"lib": "version A"},
			wantErr:        false,
		},
		{name: "non semantic version dependency upgrade", fields: fields{
			udm: upgradeDependencyMock{
				AvailableVersionMap: map[string][]versionBuildinfo{"lib": {{Version: "version A", BuildSuccessful: true}, {Version: "version B", BuildSuccessful: true}, {Version: "version C", BuildSuccessful: true}}},
				VersionMap:          map[string]string{"lib": "version A"},
			},
		}, args: args{
			service:          ports.Service{Id: "my_service", Path: "path", Dependencies: []ports.PackageDependency{{Id: "lib", Version: "version A"}}},
			d:                ports.PackageDependency{Id: "lib", Version: "version A"},
			keepMajorVersion: false,
		}, want: VersionUpgradeSpec{
			previous:        "version A",
			target:          "version C",
			latestAvailable: "version C",
			latestWorking:   "version C",
		},
			wantVersionMap: map[string]string{"lib": "version C"},
			wantErr:        false,
		},
		{name: "non semantic version dependency upgrade not possible", fields: fields{
			udm: upgradeDependencyMock{
				AvailableVersionMap: map[string][]versionBuildinfo{"lib": {{Version: "version A", BuildSuccessful: true}, {Version: "version B", BuildSuccessful: true}, {Version: "version C", BuildSuccessful: false}}},
				VersionMap:          map[string]string{"lib": "version A"},
			},
		}, args: args{
			service:          ports.Service{Id: "my_service", Path: "path", Dependencies: []ports.PackageDependency{{Id: "lib", Version: "version A"}}},
			d:                ports.PackageDependency{Id: "lib", Version: "version A"},
			keepMajorVersion: false,
		}, want: VersionUpgradeSpec{
			previous:        "version A",
			target:          "version C",
			latestAvailable: "version C",
			latestWorking:   "version A",
		},
			wantVersionMap: map[string]string{"lib": "version A"},
			wantErr:        false,
		},
		{name: "semantic version dependency up to date", fields: fields{
			udm: upgradeDependencyMock{
				AvailableVersionMap: map[string][]versionBuildinfo{"lib": {{Version: "1.0.0", BuildSuccessful: true}}},
				VersionMap:          map[string]string{"lib": "1.0.0"},
			},
		}, args: args{
			service:          ports.Service{Id: "my_service", Path: "path", Dependencies: []ports.PackageDependency{{Id: "lib", Version: "1.0.0"}}},
			d:                ports.PackageDependency{Id: "lib", Version: "1.0.0"},
			keepMajorVersion: false,
		}, want: VersionUpgradeSpec{
			previous:        "1.0.0",
			target:          "1.0.0",
			latestAvailable: "1.0.0",
			latestWorking:   "1.0.0",
		},
			wantVersionMap: map[string]string{"lib": "1.0.0"},
			wantErr:        false,
		},
		{name: "semantic version dependency upgrade", fields: fields{
			udm: upgradeDependencyMock{
				AvailableVersionMap: map[string][]versionBuildinfo{"lib": {{Version: "1.0.0", BuildSuccessful: true}, {Version: "1.1.0", BuildSuccessful: true}, {Version: "2.0.0", BuildSuccessful: true}}},
				VersionMap:          map[string]string{"lib": "1.0.0"},
			},
		}, args: args{
			service:          ports.Service{Id: "my_service", Path: "path", Dependencies: []ports.PackageDependency{{Id: "lib", Version: "1.0.0"}}},
			d:                ports.PackageDependency{Id: "lib", Version: "1.0.0"},
			keepMajorVersion: false,
		}, want: VersionUpgradeSpec{
			previous:        "1.0.0",
			target:          "2.0.0",
			latestAvailable: "2.0.0",
			latestWorking:   "2.0.0",
		},
			wantVersionMap: map[string]string{"lib": "2.0.0"},
			wantErr:        false,
		},
		{name: "semantic version dependency upgrade keep major", fields: fields{
			udm: upgradeDependencyMock{
				AvailableVersionMap: map[string][]versionBuildinfo{"lib": {{Version: "1.0.0", BuildSuccessful: true}, {Version: "1.1.0", BuildSuccessful: true}, {Version: "2.0.0", BuildSuccessful: true}}},
				VersionMap:          map[string]string{"lib": "1.0.0"},
			},
		}, args: args{
			service:          ports.Service{Id: "my_service", Path: "path", Dependencies: []ports.PackageDependency{{Id: "lib", Version: "1.0.0"}}},
			d:                ports.PackageDependency{Id: "lib", Version: "1.0.0"},
			keepMajorVersion: true,
		}, want: VersionUpgradeSpec{
			previous:        "1.0.0",
			target:          "1.1.0",
			latestAvailable: "2.0.0",
			latestWorking:   "1.1.0",
		},
			wantVersionMap: map[string]string{"lib": "1.1.0"},
			wantErr:        false,
		},
		{name: "semantic version dependency upgrade not fully possible", fields: fields{
			udm: upgradeDependencyMock{
				AvailableVersionMap: map[string][]versionBuildinfo{"lib": {{Version: "1.0.0", BuildSuccessful: true}, {Version: "1.1.0", BuildSuccessful: true}, {Version: "2.0.0", BuildSuccessful: false}}},
				VersionMap:          map[string]string{"lib": "1.0.0"},
			},
		}, args: args{
			service:          ports.Service{Id: "my_service", Path: "path", Dependencies: []ports.PackageDependency{{Id: "lib", Version: "1.0.0"}}},
			d:                ports.PackageDependency{Id: "lib", Version: "1.0.0"},
			keepMajorVersion: false,
		}, want: VersionUpgradeSpec{
			previous:        "1.0.0",
			target:          "2.0.0",
			latestAvailable: "2.0.0",
			latestWorking:   "1.1.0",
		},
			wantVersionMap: map[string]string{"lib": "1.1.0"},
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ServiceApi{
				ServicePersistence: tt.fields.udm,
				ServiceBuilder:     tt.fields.udm,
				DependencyInfo:     tt.fields.udm,
				DependencyWriter:   tt.fields.udm,
				Stdout:             ioutil.Discard,
				Stderr:             ioutil.Discard,
			}
			s.upgradeDependency(tt.args.service, tt.args.d, tt.args.keepMajorVersion)
			got, err := s.upgradeDependency(tt.args.service, tt.args.d, tt.args.keepMajorVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceApi.upgradeDependency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.fields.udm.VersionMap, tt.wantVersionMap) {
				t.Errorf("ServiceApi.upgradeDependency() version map= %v, want %v", tt.fields.udm.VersionMap, tt.wantVersionMap)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceApi.upgradeDependency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toText(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "single line", args: args{lines: []string{"abc"}}, want: "abc"},
		{name: "multiple lines", args: args{lines: []string{"abc", "xyz"}}, want: "abc\nxyz"},
		{name: "empty", args: args{lines: []string{}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toText(tt.args.lines); got != tt.want {
				t.Errorf("toText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lines(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "single line", args: args{s: "abc"}, want: []string{"abc"}},
		{name: "single line with newline", args: args{s: "abc\n"}, want: []string{"abc"}},
		{name: "multiple lines", args: args{s: "abc\nxyz\n"}, want: []string{"abc", "xyz"}},
		{name: "empty", args: args{s: ""}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lines(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeLines(t *testing.T) {
	type args struct {
		base     []string
		incoming []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "single incoming line already there", args: args{base: []string{"abc", "def"}, incoming: []string{"def"}}, want: []string{"abc", "def"}},
		{name: "single incoming line not there", args: args{base: []string{"abc", "def"}, incoming: []string{"xyz"}}, want: []string{"abc", "def", "xyz"}},
		{name: "multiple incoming lines", args: args{base: []string{"abc", "def"}, incoming: []string{"def", "xyz"}}, want: []string{"abc", "def", "xyz"}},
		{name: "empty incoming lines", args: args{base: []string{"abc", "def"}, incoming: []string{}}, want: []string{"abc", "def"}},
		{name: "empty base lines", args: args{base: []string{}, incoming: []string{"def"}}, want: []string{"def"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeLines(tt.args.base, tt.args.incoming); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeLines() = %v, want %v", got, tt.want)
			}
		})
	}
}
