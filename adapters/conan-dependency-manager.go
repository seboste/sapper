package adapters

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/seboste/sapper/ports"
)

type ConanDependencyManager struct {
}

var dependencyExp = regexp.MustCompile(`\s*(.*)/(.*)\s*`)

func (cdm ConanDependencyManager) Read(s ports.Service) ([]ports.PackageDependency, error) {
	dependencies := []ports.PackageDependency{}

	sectionExp := regexp.MustCompile(`\[(.*)\]`)

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
			dependencyMatch := dependencyExp.FindStringSubmatch(line)
			if len(dependencyMatch) == 3 {
				dependency := ports.PackageDependency{Id: dependencyMatch[1], Version: dependencyMatch[2]}
				dependencies = append(dependencies, dependency)
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
