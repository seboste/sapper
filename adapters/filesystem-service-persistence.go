package adapters

import (
	"io/ioutil"
	"path/filepath"

	"github.com/seboste/sapper/ports"
	"gopkg.in/yaml.v3"
)

type FileSystemServicePersistence struct {
}

func (fsp FileSystemServicePersistence) Load(path string) (ports.Service, error) {
	return ports.Service{}, nil
}

func (fsp FileSystemServicePersistence) Save(s ports.Service) error {

	yamlData, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	sapperFilePath := filepath.Join(s.Path, "sapperfile.yaml")
	err = ioutil.WriteFile(sapperFilePath, yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}
