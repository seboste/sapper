package core

import (
	"bufio"
	"errors"
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
		p, err := ResolveParameters(brick.Parameters, pr)
		if err != nil {
			return nil, err
		}
		for k, v := range p {
			combinedParameters[k] = v
		}
	}
	return combinedParameters, nil
}

func AddSingleBrick(s *ports.Service, b ports.Brick, parameters map[string]string) error {
	for _, f := range b.Files {
		inputFilePath := filepath.Join(b.BasePath, f)
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
		contentStr := string(content)

		contentStr = replaceParameters(contentStr, parameters)

		if _, err := os.Stat(outputFilePath); errors.Is(err, os.ErrNotExist) { //file does not exit => just write out the content
			err = ioutil.WriteFile(outputFilePath, []byte(contentStr), 0644)
			if err != nil {
				return err
			}
		} else { //file exists => merge Sections
			inputSectionSlice, err := readSections(contentStr)
			if err != nil {
				return err
			}
			inputSections := toMap(inputSectionSlice)

			outputContent, err := ioutil.ReadFile(outputFilePath)
			if err != nil {
				return err
			}

			outputContentStr := string(outputContent)

			mergedOutputContentStr, err := mergeSections(outputContentStr, inputSections)

			err = ioutil.WriteFile(outputFilePath, []byte(mergedOutputContentStr), 0644)
			if err != nil {
				return err
			}
		}
	}

	s.BrickIds = append(s.BrickIds, ports.BrickDependency{b.Id, b.Version})

	return nil
}

func mergeSection(base section, incoming section) (string, error) {

	if base.name != incoming.name {
		return "", fmt.Errorf("Unable to merge section %s with section %s. Names must match.", base.name, incoming.name)
	}

	if base.verb != "" {
		return "", fmt.Errorf("Unable to merge section %s. Base operation must not be defined.", base.name)
	}

	if incoming.verb == "REPLACE" {
		return incoming.content, nil
	} else if incoming.verb == "PREPEND" {
		var content string
		content = content + incoming.content
		if base.content != "" && incoming.content != "" {
			content = content + fmt.Sprintln("")
		}
		content = content + base.content
		return content, nil
	} else if incoming.verb == "APPEND" {
		var content string
		content = content + base.content
		if base.content != "" && incoming.content != "" {
			content = content + fmt.Sprintln("")
		}
		content = content + incoming.content
		return content, nil
	} else {
		return "", fmt.Errorf("Unable to merge section %s. Invalid incoming operation %s.", base.name, incoming.verb)
	}
}

func mergeSections(content string, inputSections map[string]section) (string, error) {
	outputSections, err := readSections(content)
	if err != nil {
		return content, err
	}

	scanner := bufio.NewScanner(strings.NewReader(content))

	outputContent := ""
	lineNumber := 0
	for _, s := range outputSections {
		//advance to next section
		for lineNumber < s.lineBegin && scanner.Scan() {
			outputContent = outputContent + fmt.Sprintln(scanner.Text())
			lineNumber = lineNumber + 1
		}

		if incomingSection, ok := inputSections[s.name]; ok {
			mergedSectionContent, err := mergeSection(s, incomingSection)
			if err != nil {
				return "", err
			}
			if mergedSectionContent != "" {
				outputContent = outputContent + fmt.Sprintln(mergedSectionContent)
			}
		} else { // no incoming section => just use base content
			outputContent = outputContent + s.content
		}
		for lineNumber < s.lineEnd && scanner.Scan() { //skip section content
			lineNumber = lineNumber + 1
		}

		//take care of end tag
		scanner.Scan()
		outputContent = outputContent + fmt.Sprintln(scanner.Text())
		lineNumber = lineNumber + 1
	}

	//copy incoming content after last section
	for scanner.Scan() {
		outputContent = outputContent + fmt.Sprintln(scanner.Text())
		lineNumber = lineNumber + 1
	}

	return outputContent, nil
}

func replaceParameters(content string, parameters map[string]string) string {
	for name, value := range parameters {
		pattern := "<<<" + name + ">>>"
		content = strings.ReplaceAll(content, pattern, value)
	}
	return content
}

func GetBricksRecursive(brickId string, db ports.BrickDB) ([]ports.Brick, error) {
	bricks := []ports.Brick{}
	brickIds := make(map[string]bool)

	brick, err := db.Brick(brickId)
	if err != nil {
		return bricks, fmt.Errorf("invalid brick %s", brickId)
	}

	bricks = append(bricks, brick)
	brickIds[brick.Id] = true
	for _, dependencyId := range brick.Dependencies {
		dependencies, err := GetBricksRecursive(dependencyId, db)
		if err != nil {
			return nil, err
		}
		for _, dedependency := range dependencies {

			if brickIds[dedependency.Id] == false {
				bricks = append(bricks, dedependency)
				brickIds[dedependency.Id] = true
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

	if len(bricks) == 0 {
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
		if err := AddSingleBrick(&service, brick, parameters); err != nil {
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
