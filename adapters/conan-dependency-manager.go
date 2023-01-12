package adapters

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/seboste/sapper/ports"
)

type ConanDependencyManager struct {
}

type ConanDependency struct {
	Id        string
	Version   string
	User      string
	Channel   string
	Reference string
}

var dependencyExp = regexp.MustCompile(`([^@\/#\s]+)\/([^@\/#\s]+)(@([^@\/#\s]+)\/([^@\/#\s]+))?(#([0-9a-fA-F]+))?`)

func parseConanDependency(input string) (ConanDependency, error) {
	m := dependencyExp.FindStringSubmatch(input)
	if len(m) != 8 {
		return ConanDependency{}, fmt.Errorf("unable to parse %s. It needs to be in the format 'lib/version@user/channel#reference'. 'lib' and 'version' are required.", input)
	}

	return ConanDependency{Id: m[1], Version: m[2], User: m[4], Channel: m[5], Reference: m[7]}, nil
}

var sectionExp = regexp.MustCompile(`\[(.*)\]`)

func (cdm ConanDependencyManager) Read(s ports.Service) ([]ports.PackageDependency, error) {
	dependencies := []ports.PackageDependency{}

	conanFilePath := filepath.Join(s.Path, "conanfile.txt")
	f, err := os.Open(conanFilePath)
	if err != nil {
		return []ports.PackageDependency{}, err
	}

	scanner := bufio.NewScanner(f)
	currentSection := ""
	for scanner.Scan() {
		line := scanner.Text()

		sectionMatch := sectionExp.FindStringSubmatch(line)
		if len(sectionMatch) == 2 {
			currentSection = sectionMatch[1]
		}

		if currentSection == "requires" {
			dep, err := parseConanDependency(line)
			if err == nil {
				dependencies = append(dependencies, ports.PackageDependency{Id: dep.Id, Version: dep.Version})
			}
		}
	}

	return dependencies, nil
}

func (cdm ConanDependencyManager) Write(s ports.Service, dependency string, version string) error {
	return nil
}

func (cdm ConanDependencyManager) AvailableVersions(dependency string) ([]string, error) {
	var buffer bytes.Buffer
	cmd := exec.Command("conan", "search", "-r=all", dependency)

	cmd.Stdout = &buffer
	//cmd.Stderr = os.Stderr
	err := cmd.Run()

	scanner := bufio.NewScanner(&buffer)
	versions := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		dependencyMatch := dependencyExp.FindStringSubmatch(line)
		if len(dependencyMatch) == 3 && dependencyMatch[1] == dependency {
			versions = append(versions, dependencyMatch[2])
		}
	}

	return versions, err
}

var _ ports.DependencyReader = ConanDependencyManager{}
var _ ports.DependencyWriter = ConanDependencyManager{}
var _ ports.DependencyInfo = ConanDependencyManager{}
