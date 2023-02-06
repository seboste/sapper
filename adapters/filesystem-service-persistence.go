package adapters

import (
	"io/ioutil"
	"path/filepath"

	"github.com/seboste/sapper/ports"
	"gopkg.in/yaml.v3"
)

type FileSystemServicePersistence struct {
	DependencyReader ports.ServicePackageDependencyReader
}

func (fsp FileSystemServicePersistence) Load(path string) (ports.Service, error) {
	s := ports.Service{Path: path}

	sapperFilePath := filepath.Join(s.Path, "sapperfile.yaml")

	yamlData, err := ioutil.ReadFile(sapperFilePath)
	if err != nil {
		return s, err
	}

	if err := yaml.Unmarshal(yamlData, &s); err != nil {
		return s, err
	}

	s.Dependencies, err = fsp.DependencyReader.ReadFromService(s)
	if err != nil {
		return s, err
	}

	return s, nil
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
