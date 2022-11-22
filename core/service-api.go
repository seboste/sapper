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

type ServiceApi struct {
	Db                 ports.BrickDB
	ServicePersistence ports.ServicePersistence
	ParameterResolver  ports.ParameterResolver
}

func ResolveParameters(bp []ports.BrickParameters, pr ports.ParameterResolver) (map[string]string, error) {
	parameters := make(map[string]string)
	for _, p := range bp {
		value := pr.Resolve(p.Name)
		if value == "" {
			value = p.Default
		}
		if value == "" {
			return nil, fmt.Errorf("unable to resolve value for parameter %s", p.Name)
		}
		parameters[p.Name] = value
	}
	return parameters, nil
}

func ResolveParameterSlice(bricks []ports.Brick, pr ports.ParameterResolver) (map[string]string, error) {
	combinedParameters := map[string]string{}
	for _, brick := range bricks {
		p, err := ResolveParameters(brick.GetParameters(), pr)
		if err != nil {
			return map[string]string{}, err
		}
		for k, v := range p {
			combinedParameters[k] = v
		}
	}
	return combinedParameters, nil
}

func AddBrick(s *ports.Service, b ports.Brick, parameters map[string]string) error {
	for _, f := range b.GetFiles() {
		inputFilePath := filepath.Join(b.GetBasePath(), f)
		if _, err := os.Stat(inputFilePath); err != nil {
			return err
		}

		outputFilePath := filepath.Join(s.Path, f)
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

	s.BrickIds = append(s.BrickIds, ports.BrickDependency{b.GetId(), b.GetVersion()})

	return nil
}

func GetBricksRecursive(brickId string, db ports.BrickDB) ([]ports.Brick, error) {
	bricks := []ports.Brick{}
	brickIds := make(map[string]bool)

	brick := db.Brick(brickId)
	if brick == nil {
		return bricks, fmt.Errorf("invalid brick %s", brickId)
	}

	bricks = append(bricks, brick)
	brickIds[brick.GetId()] = true
	for _, dependencyId := range brick.GetDependencies() {
		dependencies, err := GetBricksRecursive(dependencyId, db)
		if err != nil {
			return nil, err
		}
		for _, dedependency := range dependencies {

			if brickIds[dedependency.GetId()] == false {
				bricks = append(bricks, dedependency)
				brickIds[dedependency.GetId()] = true
			}
		}
	}

	return bricks, nil
}

func (s ServiceApi) Add(templateName string, parentDir string, parameterResolver ports.ParameterResolver) error {

	bricks, err := GetBricksRecursive(templateName, s.Db)
	if err != nil {
		return err
	}

	parameters, err := ResolveParameterSlice(bricks, parameterResolver)
	if err != nil {
		return err
	}

	template := bricks[0]
	if template == nil {
		return fmt.Errorf("invalid template %s", templateName)
	}
	serviceName := parameters["NAME"]
	if serviceName == "" {
		return fmt.Errorf("invalid service name %s", serviceName)
	}
	outputBasePath := filepath.Join(parentDir, serviceName)
	if err := os.MkdirAll(outputBasePath, os.ModePerm); err != nil {
		return err
	}

	service := ports.Service{Id: serviceName, Path: outputBasePath}

	for _, brick := range bricks {
		if err := AddBrick(&service, brick, parameters); err != nil {
			return err
		}
	}

	if err := s.ServicePersistence.Save(service); err != nil {
		return err
	}
	return nil
}

func (s ServiceApi) Update() {
	fmt.Println("update")
}

func (s ServiceApi) Build(path string) error {
	cmd := exec.Command("make", "build", "-B")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (s ServiceApi) Test() {
	fmt.Println("test")
}

func (s ServiceApi) Deploy() {
	fmt.Println("deploy")
}

var _ ports.ServiceApi = ServiceApi{}
