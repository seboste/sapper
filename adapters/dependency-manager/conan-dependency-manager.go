package dependencyManager

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

func processLines(r io.Reader, op func(line string, isActive bool), p ports.PackageDependencySectionPredicate) {
	scanner := bufio.NewScanner(r)
	state := ""
	for scanner.Scan() {
		line := scanner.Text()
		var isActive bool
		isActive, state = p(line, state)
		op(line, isActive)
	}
}

func readDependenciesFromConanfile(path string, p ports.PackageDependencySectionPredicate) ([]ports.PackageDependency, error) {
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

func (cdm ConanDependencyManager) ReadFromBrick(b ports.Brick, p ports.PackageDependencySectionPredicate) ([]ports.PackageDependency, error) {
	deps, err := readDependenciesFromConanfile(b.BasePath, p)
	if os.IsNotExist(err) { //it is ok if the conanfile.txt does not exist for a brick => there are just no dependencies
		return []ports.PackageDependency{}, nil
	}
	return deps, err
}

func writeDependenciesToConanfile(path string, dependencies []ports.PackageDependency, p ports.PackageDependencySectionPredicate) error {

	dependencyMap := map[string]string{}
	for _, d := range dependencies {
		dependencyMap[d.Id] = d.Version
	}

	conanfilePath := filepath.Join(path, "conanfile.txt")
	content, err := ioutil.ReadFile(conanfilePath)
	if err != nil {
		return err
	}

	var outputContent string
	processLines(strings.NewReader(string(content)), func(line string, isActive bool) {
		if isActive {
			dep, err := parseConanDependency(line)
			newVersion, ok := dependencyMap[dep.Id]
			if err == nil && ok {
				dep.Version = newVersion
				line = replaceConanDependency(line, dep)
				delete(dependencyMap, dep.Id)
			}
		}

		outputContent = outputContent + fmt.Sprintln(line)
	}, p)
	if len(dependencyMap) > 0 {
		return fmt.Errorf("unable to write all dependencies")
	}
	if err := ioutil.WriteFile(conanfilePath, []byte(outputContent), 0644); err != nil {
		return err
	}

	return nil
}

func (cdm ConanDependencyManager) WriteToService(s ports.Service, d ports.PackageDependency) error {
	return writeDependenciesToConanfile(s.Path, []ports.PackageDependency{d}, isInConanRequiresSection)

}

func (cdm ConanDependencyManager) WriteToBrick(b ports.Brick, d []ports.PackageDependency, p ports.PackageDependencySectionPredicate) error {
	return writeDependenciesToConanfile(b.BasePath, d, p)
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

var _ ports.BrickPackageDependencyReader = ConanDependencyManager{}
var _ ports.ServicePackageDependencyReader = ConanDependencyManager{}
var _ ports.BrickPackageDependencyWriter = ConanDependencyManager{}
var _ ports.ServicePackageDependencyWriter = ConanDependencyManager{}
var _ ports.DependencyInfo = ConanDependencyManager{}
