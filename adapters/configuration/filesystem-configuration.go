package configuration

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/seboste/sapper/ports"
	"gopkg.in/yaml.v3"
)

type FileSystemConfiguration struct {
	Path string         `yaml:"-"`
	Rmts []ports.Remote `yaml:"Remotes"`
}

var defaultRemote ports.Remote = ports.Remote{
	Name: "sapper-bricks",
	Kind: ports.GitRemote,
	Src:  "https://github.com/seboste/sapper-bricks.git",
}

var defaultConfiguration FileSystemConfiguration

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultConfiguration.Path = filepath.Join(homeDir, ".sapper")
	defaultConfiguration.Rmts = []ports.Remote{defaultRemote}
}

func MakeFilesystemConfiguration() (FileSystemConfiguration, error) {
	fsc := FileSystemConfiguration{Path: defaultConfiguration.Path}
	err := fsc.Load()
	if err != nil && os.IsNotExist(err) { //config does not exit => write default config and retry
		err = defaultConfiguration.Save()
		if err != nil {
			return fsc, err
		}
		err = fsc.Load()
	}
	return fsc, nil
}

func (fsc FileSystemConfiguration) ConfigPath() string {
	return filepath.Join(fsc.Path, "config.yaml")
}

func (fsc *FileSystemConfiguration) Load() error {
	yamlFile, err := ioutil.ReadFile(fsc.ConfigPath())
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, fsc)
	if err != nil {
		return err
	}
	return nil
}

func (fsc FileSystemConfiguration) Save() error {
	if err := os.MkdirAll(fsc.Path, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(fsc.DefaultRemotesDir(), os.ModePerm); err != nil {
		return err
	}

	yamlData, err := yaml.Marshal(fsc)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fsc.ConfigPath(), yamlData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (fsc FileSystemConfiguration) DefaultRemotesDir() string {
	return path.Join(fsc.Path, "remotes")
}
func (fsc FileSystemConfiguration) Remotes() []ports.Remote {
	return fsc.Rmts
}

func (fsc *FileSystemConfiguration) UpdateRemotes(remotes []ports.Remote) {
	fsc.Rmts = remotes
}

var _ ports.Configuration = (*FileSystemConfiguration)(nil)
