package adapters

import (
	"os"
	"path/filepath"

	"github.com/seboste/sapper/ports"
)

type FileSystemConfiguration struct {
	path    string
	remotes []ports.Remote
}

func MakeFilesystemConfiguration(basePath string) (FileSystemConfiguration, error) {
	path := filepath.Join(basePath, ".sapper")
	fsc := FileSystemConfiguration{path: path}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fsc, err
	}

	return fsc, nil
}

func (fsc FileSystemConfiguration) ConfigurationDir() string {
	return fsc.path
}
func (fsc FileSystemConfiguration) Remotes() []ports.Remote {
	return fsc.remotes
}

func (fsc FileSystemConfiguration) UpdateRemotes(remotes []ports.Remote) error {
	return nil
}

var _ ports.Configuration = FileSystemConfiguration{}
