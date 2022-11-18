package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/seboste/sapper/ports"
)

type Service struct {
	Db ports.BrickDB
}

func (s Service) Add(name string, version string, templateName string, parentDir string) error {

	template := s.Db.Brick(templateName)
	if template == nil {
		return fmt.Errorf("invalid template %s", templateName)
	}

	parameters := make(map[string]string)
	for _, p := range template.GetParameters() {
		parameters[p.Name] = p.Default
	}

	parameters["NAME"] = name

	outputBasePath := filepath.Join(parentDir, parameters["NAME"])
	if err := os.MkdirAll(outputBasePath, os.ModePerm); err != nil {
		return err
	}

	for _, f := range template.GetFiles() {
		inputFilePath := filepath.Join(template.GetBasePath(), f)
		if _, err := os.Stat(inputFilePath); err != nil {
			return err
		}

		outputFilePath := filepath.Join(outputBasePath, f)
		outputDir, _ := filepath.Split(outputFilePath)
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return err
		}

		content, err := ioutil.ReadFile(inputFilePath)
		if err != nil {
			return err
		}

		for name, value := range parameters {
			pattern := "<<<" + name + ">>>"
			content = []byte(strings.ReplaceAll(string(content), pattern, value))
		}

		err = ioutil.WriteFile(outputFilePath, content, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Service) Update() {
	fmt.Println("update")
}

func (s Service) Build(path string) error {
	cmd := exec.Command("make", "build", "-B")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (s Service) Test() {
	fmt.Println("test")
}

func (s Service) Deploy() {
	fmt.Println("deploy")
}

var _ ports.ServiceApi = Service{}
