package adapters

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/seboste/sapper/ports"
)

func TestFileSystemConfiguration_Load(t *testing.T) {

	tempDir, _ := ioutil.TempDir("", "fscLoadTest*")
	defer os.RemoveAll(tempDir) // clean up

	type fields struct {
		Path string
		Yaml string
		Rmts []ports.Remote
	}
	tests := []struct {
		name       string
		fields     fields
		wantConfig FileSystemConfiguration
		wantErr    func(error) bool
	}{
		{name: "correct config", fields: fields{Path: tempDir, Yaml: `Remotes:
    - name: some-remote
      path: some-path
`, Rmts: []ports.Remote{}},
			wantConfig: FileSystemConfiguration{Path: tempDir, Rmts: []ports.Remote{{Name: "some-remote", Path: "some-path"}}},
			wantErr:    func(err error) bool { return err == nil }},
		{name: "invalid yaml syntax", fields: fields{Path: tempDir, Yaml: `Remotes:
    - name: some-remote
         path: some-path
`, Rmts: []ports.Remote{}},
			wantConfig: FileSystemConfiguration{Path: tempDir, Rmts: []ports.Remote{}},
			wantErr:    func(err error) bool { return err != nil && !os.IsNotExist(err) }},
		{name: "file does not exist", fields: fields{Path: tempDir, Yaml: "", Rmts: []ports.Remote{}},
			wantConfig: FileSystemConfiguration{Path: tempDir, Rmts: []ports.Remote{}},
			wantErr:    os.IsNotExist},
		{name: "path does not exist", fields: fields{Path: filepath.Join(tempDir, "somepath"), Yaml: "", Rmts: []ports.Remote{}},
			wantConfig: FileSystemConfiguration{Path: filepath.Join(tempDir, "somepath"), Rmts: []ports.Remote{}},
			wantErr:    os.IsNotExist},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.Yaml != "" {
				ioutil.WriteFile(filepath.Join(tempDir, "config.yaml"), []byte(tt.fields.Yaml), 0644)
			}
			defer os.Remove(filepath.Join(tempDir, "config.yaml"))

			fsc := FileSystemConfiguration{
				Path: tt.fields.Path,
				Rmts: tt.fields.Rmts,
			}
			if err := fsc.Load(); !tt.wantErr(err) {
				t.Errorf("FileSystemConfiguration.Load() error = %v, wantErr %v", err, tt.wantErr(err))
			}

			if !reflect.DeepEqual(fsc, tt.wantConfig) {
				t.Errorf("FileSystemConfiguration.Load() config = %v, want %v", fsc, tt.wantConfig)
			}
		})
	}
}

func TestFileSystemConfiguration_ConfigPath(t *testing.T) {
	type fields struct {
		Path string
		Rmts []ports.Remote
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "returns config.yaml", fields: fields{Path: "/somePath"}, want: "/somePath/config.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsc := FileSystemConfiguration{
				Path: tt.fields.Path,
				Rmts: tt.fields.Rmts,
			}
			if got := fsc.ConfigPath(); got != tt.want {
				t.Errorf("FileSystemConfiguration.ConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileSystemConfiguration_DefaultRemotesDir(t *testing.T) {
	type fields struct {
		Path string
		Rmts []ports.Remote
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "returns remotes", fields: fields{Path: "/somePath"}, want: "/somePath/remotes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsc := FileSystemConfiguration{
				Path: tt.fields.Path,
				Rmts: tt.fields.Rmts,
			}
			if got := fsc.DefaultRemotesDir(); got != tt.want {
				t.Errorf("FileSystemConfiguration.DefaultRemotesDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileSystemConfiguration_Save(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "fscSaveTest*")
	defer os.RemoveAll(tempDir) // clean up

	type fields struct {
		Path string
		Rmts []ports.Remote
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		wantYaml string
	}{
		{name: "config without remotes", fields: fields{Path: filepath.Join(tempDir, "test1")}, wantErr: false, wantYaml: `Remotes: []
`},
		{name: "config with remotes", fields: fields{Path: filepath.Join(tempDir, "test2"), Rmts: []ports.Remote{{Name: "some-remote", Path: "some-path"}}}, wantErr: false, wantYaml: `Remotes:
    - name: some-remote
      path: some-path
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsc := FileSystemConfiguration{
				Path: tt.fields.Path,
				Rmts: tt.fields.Rmts,
			}
			if err := fsc.Save(); (err != nil) != tt.wantErr {
				t.Errorf("FileSystemConfiguration.Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotYaml, _ := ioutil.ReadFile(fsc.ConfigPath())
			if !reflect.DeepEqual(string(gotYaml), tt.wantYaml) {
				t.Errorf("FileSystemConfiguration.Save() yaml = %v, want %v", string(gotYaml), tt.wantYaml)
			}
		})
	}
}

func TestMakeFilesystemConfiguration(t *testing.T) {

	testConfig := FileSystemConfiguration{Rmts: []ports.Remote{{Name: "some-remote", Path: "some-path"}}}

	tests := []struct {
		name          string
		initialConfig *FileSystemConfiguration
		want          FileSystemConfiguration
		wantErr       bool
	}{
		{name: "config exists", initialConfig: &testConfig, want: testConfig, wantErr: false},
		{name: "config does not exist", initialConfig: nil, want: defaultConfiguration, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savePath := defaultConfiguration.Path
			defer func() { defaultConfiguration.Path = savePath }()

			defaultConfiguration.Path, _ = ioutil.TempDir("", "MakeFscTest*")
			os.RemoveAll(defaultConfiguration.Path)

			if tt.initialConfig != nil {
				tt.initialConfig.Path = defaultConfiguration.Path
				tt.initialConfig.Save()
			}
			tt.want.Path = defaultConfiguration.Path

			got, err := MakeFilesystemConfiguration()
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeFilesystemConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeFilesystemConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
