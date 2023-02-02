package adapters

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

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

func (dep ConanDependency) String() string {
	var str string
	str = str + fmt.Sprintf("%s/%s", dep.Id, dep.Version)
	if dep.User != "" || dep.Channel != "" {
		if dep.User == "" {
			dep.User = "_"
		}
		if dep.Channel == "" {
			dep.Channel = "_"
		}
		str = str + fmt.Sprintf("@%s/%s", dep.User, dep.Channel)
	}
	if dep.Reference != "" {
		str = str + fmt.Sprintf("#%s", dep.Reference)
	}
	return str
}

var dependencyExp = regexp.MustCompile(`([^@\/#\s]+)\/([^@\/#\s]+)(@([^@\/#\s]+)\/([^@\/#\s]+))?(#([0-9a-fA-F]+))?`)

func parseConanDependency(input string) (ConanDependency, error) {
	m := dependencyExp.FindStringSubmatch(input)
	if len(m) != 8 {
		return ConanDependency{}, fmt.Errorf("unable to parse %s. It needs to be in the format 'lib/version@user/channel#reference'. 'lib' and 'version' are required.", input)
	}

	return ConanDependency{Id: m[1], Version: m[2], User: m[4], Channel: m[5], Reference: m[7]}, nil
}

func replaceConanDependency(line string, dep ConanDependency) string {
	return dependencyExp.ReplaceAllString(line, fmt.Sprint(dep))
}

var sectionExp = regexp.MustCompile(`\[(.*)\]`)
var commentExp = regexp.MustCompile(`[\s]*#.*`)

func isInConanRequiresSection(line string, state string) (bool, string) {
	section := state                                                        //return this section until a different section is active
	if loc := commentExp.FindStringIndex(line); loc != nil && loc[0] == 0 { // '#'s only indicate a comment if they are at the beginnig of a line
		section = "comment" //use special section 'comment' to indicate that this line is in a comment. Do not change the state here
	} else {
		sectionMatch := sectionExp.FindStringSubmatch(line)
		if len(sectionMatch) == 2 {
			section = sectionMatch[1]
			state = section //state change
		}
	}
	return section == "requires", state
}

func processLines(r io.Reader, op func(line string, isActive bool), p ports.PacakgeDependencySectionPredicate) {
	scanner := bufio.NewScanner(r)
	state := ""
	for scanner.Scan() {
		line := scanner.Text()
		var isActive bool
		isActive, state = p(line, state)
		op(line, isActive)
	}
}

func readDependenciesFromConanfile(path string, p ports.PacakgeDependencySectionPredicate) ([]ports.PackageDependency, error) {
	dependencies := []ports.PackageDependency{}

	conanFilePath := filepath.Join(path, "conanfile.txt")
	f, err := os.Open(conanFilePath)
	if err != nil {
		return []ports.PackageDependency{}, err
	}

	processLines(f, func(line string, isActive bool) {
		if isActive {
			dep, err := parseConanDependency(line)
			if err == nil {
				dependencies = append(dependencies, ports.PackageDependency{Id: dep.Id, Version: dep.Version})
			}
		}
	}, p)

	return dependencies, nil
}

func (cdm ConanDependencyManager) ReadFromService(s ports.Service) ([]ports.PackageDependency, error) {
	return readDependenciesFromConanfile(s.Path, isInConanRequiresSection)
}

func (cdm ConanDependencyManager) ReadFromBrick(b ports.Brick, p ports.PacakgeDependencySectionPredicate) ([]ports.PackageDependency, error) {
	return readDependenciesFromConanfile(b.BasePath, p)
}

func (cdm ConanDependencyManager) WriteToService(s ports.Service, dependency string, version string) error {

	conanfilePath := filepath.Join(s.Path, "conanfile.txt")
	content, err := ioutil.ReadFile(conanfilePath)
	if err != nil {
		return err
	}

	var outputContent string
	replaceCount := 0
	processLines(strings.NewReader(string(content)), func(line string, isActive bool) {
		if isActive {
			dep, err := parseConanDependency(line)
			if err == nil && dep.Id == dependency {
				dep.Version = version
				line = replaceConanDependency(line, dep)
				replaceCount = replaceCount + 1
			}
		}

		outputContent = outputContent + fmt.Sprintln(line)
	}, isInConanRequiresSection)
	if replaceCount != 1 {
		return fmt.Errorf("unable to set version %s of package %s", version, dependency)
	}
	if err := ioutil.WriteFile(conanfilePath, []byte(outputContent), 0644); err != nil {
		return err
	}

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
		if len(dependencyMatch) == 8 && dependencyMatch[1] == dependency {
			versions = append(versions, dependencyMatch[2])
		}
	}

	return versions, err
}

var _ ports.PackageDependencyReader = ConanDependencyManager{}
var _ ports.PackageDependencyWriter = ConanDependencyManager{}
var _ ports.DependencyInfo = ConanDependencyManager{}
